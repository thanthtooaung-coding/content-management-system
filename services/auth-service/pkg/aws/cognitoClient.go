package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

func NewCognitoClient(cfg aws.Config) *cognitoidentityprovider.Client {
	client := cognitoidentityprovider.NewFromConfig(cfg)

	return client
}
