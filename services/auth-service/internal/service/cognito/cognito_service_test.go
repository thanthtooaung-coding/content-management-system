package cognito

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func setupCognitoService(t *testing.T) *CognitoService {
	var service *CognitoService
	app := fx.New(
		fx.Provide(func() *logrus.Logger {
			return logrus.New()
		}),
		fx.Provide(func() (aws.Config, error) {
			return config.LoadDefaultConfig(context.TODO())
		}),
		fx.Provide(NewCognitoService),
		fx.Populate(&service),
	)
	require.NoError(t, app.Start(context.Background()))
	t.Cleanup(func() {
		err := app.Stop(context.Background())
		require.NoError(t, err)
	})
	return service
}
func TestCognitoServiceLogin(t *testing.T) {
	var service = setupCognitoService(t)
	username := "swanhetaungp@gmail.com"
	password := "TesTUSER123!"
	login, err := service.Login(username, password)
	require.NoError(t, err)
	require.NotEmpty(t, login)
}

func TestCogntoServiceSignUp(t *testing.T) {
	service := setupCognitoService(t)
	userName := "swanhetaungp@gmail.com"
	password := "TesTUSER123!"
	var attrMap = map[string]string{
		"email": userName,
	}

	err := service.Register(userName, password, attrMap)
	require.NoError(t, err)
}
