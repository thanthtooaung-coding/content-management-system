package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Role struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Username         string    `gorm:"type:varchar(255);not null;unique" json:"username"`
	Password         string    `gorm:"type:varchar(255);not null" json:"-"`
	Email            string    `gorm:"type:varchar(255);not null;unique" json:"email"`
	Name             string    `gorm:"type:varchar(255)" json:"name"`
	RoleID           uuid.UUID `gorm:"type:uuid;not null" json:"role_id"`
	RegistrationDate time.Time `json:"registration_date"`
	Address          string    `gorm:"type:text" json:"address"`
	PhoneNumber      string    `gorm:"type:varchar(255)" json:"phone_number"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	Role Role `gorm:"foreignKey:RoleID" json:"role"`
}

func (r *Role) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
func databaseConnection() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" || sslmode == "" {
		return nil, fmt.Errorf("missing one or more required environment variables")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	conn, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open GORM connection: %w", err)
	}

	sqlDB, err := conn.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return conn, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := databaseConnection()
	if err != nil {
		panic(err.Error())
	}

	// Drop existing tables first
	err = db.Migrator().DropTable(&User{}, &Role{})
	if err != nil {
		log.Printf("Warning: Could not drop tables: %v", err)
	}

	err = db.AutoMigrate(&Role{}, &User{})
	if err != nil {
		panic(err.Error())
	}

	// Seed data
	roles := []Role{
		{Name: "Administrator"},
		{Name: "Customer"},
	}

	for _, role := range roles {
		err = db.Create(&role).Error
		if err != nil {
			panic(err.Error())
		}
	}

	log.Println("Successfully created")
}
