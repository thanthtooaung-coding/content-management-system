package types

import (
	"github.com/google/uuid"
	"time"
)

type Role struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Users []User `gorm:"foreignKey:RoleID" json:"users,omitempty"`
}

type User struct {
	ID               uuid.UUID `gorm:"primaryKey;autoIncrement" json:"id"`
	Username         string    `gorm:"type:varchar(255);not null;unique" json:"username"`
	Password         string    `gorm:"type:varchar(255);not null" json:"-"`
	Email            string    `gorm:"type:varchar(255);not null;unique" json:"email"`
	Name             string    `gorm:"type:varchar(255)" json:"name"`
	RoleID           uint64    `gorm:"not null" json:"role_id"`
	RegistrationDate time.Time `json:"registration_date"`
	Address          string    `gorm:"type:text" json:"address"`
	PhoneNumber      string    `gorm:"type:varchar(255)" json:"phone_number"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	Role Role `gorm:"foreignKey:RoleID" json:"role"` // relation to Role
}
