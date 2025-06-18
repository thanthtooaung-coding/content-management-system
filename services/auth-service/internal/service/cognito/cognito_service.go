package cognito

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	cognitoTypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/content-management-system/auth-service/internal/model/types"
	"github.com/sirupsen/logrus"
	"os"
)

type CognitoService struct {
	CognitoClientInterface
	userPoolClient *cognitoidentityprovider.Client
	userPoolID     string
	clientID       string
	identityPoolID string
	log            *logrus.Logger
}

func NewCognitoService(log *logrus.Logger, cfg aws.Config) *CognitoService {
	userPoolID := os.Getenv("USER_POOL_ID")
	clientID := os.Getenv("CLIENT_ID")
	log.Infof("NewCognitoService: USER_POOL_ID=%s, CLIENT_ID=%s", userPoolID, clientID)

	client := cognitoidentityprovider.NewFromConfig(cfg)
	return &CognitoService{
		userPoolClient: client,
		clientID:       clientID,
		userPoolID:     userPoolID,
		log:            log,
	}
}

func (cg *CognitoService) Login(email string, password string) (*types.AuthResult, error) {
	signInInput := cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: cognitoTypes.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(cg.clientID),
		AuthParameters: map[string]string{
			"USERNAME": email,
			"PASSWORD": password,
		},
	}
	authResp, err := cg.userPoolClient.InitiateAuth(context.Background(), &signInInput)
	if err != nil {
		cg.log.Errorf("failed to login user: %s", err.Error())
		return nil, err
	}
	if authResp.ChallengeName != "" {
		return &types.AuthResult{
			ChallengeName:       string(authResp.ChallengeName),
			Session:             aws.ToString(authResp.Session),
			ChallengeParameters: authResp.ChallengeParameters,
		}, nil
	}

	return &types.AuthResult{
		AccessToken:  aws.ToString(authResp.AuthenticationResult.AccessToken),
		IdToken:      aws.ToString(authResp.AuthenticationResult.IdToken),
		RefreshToken: aws.ToString(authResp.AuthenticationResult.RefreshToken),
		TokenType:    aws.ToString(authResp.AuthenticationResult.TokenType),
		ExpiresIn:    authResp.AuthenticationResult.ExpiresIn,
	}, nil
}

func (cg *CognitoService) Register(email, password string, userAttributes map[string]string) error {
	var attributes []cognitoTypes.AttributeType

	attributes = append(attributes, cognitoTypes.AttributeType{
		Name:  aws.String("email"),
		Value: aws.String(email),
	})
	for key, value := range userAttributes {
		attributes = append(attributes, cognitoTypes.AttributeType{
			Name:  aws.String(key),
			Value: aws.String(value),
		})
	}
	registerInput := cognitoidentityprovider.SignUpInput{
		ClientId:       aws.String(cg.clientID),
		Username:       aws.String(email),
		Password:       aws.String(password),
		UserAttributes: attributes,
	}

	_, err := cg.userPoolClient.SignUp(context.Background(), &registerInput)
	if err != nil {
		cg.log.Errorf("failed to register user: %s", err.Error())
		return err
	}

	// Catch the user sub and create user with the cognito ID in the neon database

	return nil

}

func (cg *CognitoService) ConfirmSignUp(email, code string) error {
	signUpInput := cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(cg.clientID),
		ConfirmationCode: aws.String(code),
		Username:         aws.String(email),
	}
	up, err := cg.userPoolClient.ConfirmSignUp(context.Background(), &signUpInput)
	if err != nil {
		cg.log.Errorf("failed to confirm user: %s", err.Error())
		return err
	}

	//
	cg.log.Infof("user %s confirmed successfully", &up.ResultMetadata)

	return nil
}

func (cg *CognitoService) RefreshToken(refreshToken string) (*types.AuthResult, error) {
	refreshInput := cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: cognitoTypes.AuthFlowTypeRefreshTokenAuth,
		ClientId: aws.String(cg.clientID),
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": refreshToken,
		},
	}
	authResp, err := cg.userPoolClient.InitiateAuth(context.Background(), &refreshInput)
	if err != nil {
		cg.log.Errorf("failed to refresh token: %s", err.Error())
		return nil, err
	}
	return &types.AuthResult{
		AccessToken:  aws.ToString(authResp.AuthenticationResult.AccessToken),
		IdToken:      aws.ToString(authResp.AuthenticationResult.IdToken),
		RefreshToken: aws.ToString(authResp.AuthenticationResult.RefreshToken),
		TokenType:    aws.ToString(authResp.AuthenticationResult.TokenType),
		ExpiresIn:    authResp.AuthenticationResult.ExpiresIn,
	}, nil
}

func (cg *CognitoService) GetUser(accessToken string) (*cognitoTypes.UserType, error) {
	getUserInput := cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(accessToken),
	}

	userResp, err := cg.userPoolClient.GetUser(context.Background(), &getUserInput)
	if err != nil {
		cg.log.Errorf("failed to get user: %s", err.Error())
		return nil, err
	}
	return &cognitoTypes.UserType{
		Username:   userResp.Username,
		Attributes: userResp.UserAttributes,
	}, nil
}

func (cg *CognitoService) AddUserToGroup(username, groupName string) error {
	adminAddInput := cognitoidentityprovider.AdminAddUserToGroupInput{
		GroupName:  aws.String(groupName),
		UserPoolId: aws.String(cg.userPoolID),
		Username:   aws.String(username),
	}
	group, err := cg.userPoolClient.AdminAddUserToGroup(context.Background(), &adminAddInput)
	if err != nil {
		return err
	}
	cg.log.Info("successfully added user to group", group)
	return nil
}
