
# Author Swan Htet Aung Phyo

# Provider Configuration
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.81.0"
    }
    ad = {
      source  = "hashicorp/ad"
      version = "0.5.0"
    }
  }
}
provider "aws" {
  region = "us-east-1"
}


resource "aws_cognito_user_pool" "cms_user_pool" {
  name = "cmsApplicationPool"

  auto_verified_attributes = ["email"]
  username_attributes      = ["email"]
  schema {
    name                = "email"
    attribute_data_type = "String"
    required            = true
    mutable             = false
  }

  admin_create_user_config {
    allow_admin_create_user_only = false
    invite_message_template {
      email_message = "Hello {username}, your temporary password is {####}. Please change it after you log in."
      email_subject = "Initial Account Setup"
      sms_message   = "Hi {username}, your password is {####}."
    }
  }


  mfa_configuration = "OFF"
  account_recovery_setting {
    recovery_mechanism {
      name     = "verified_email"
      priority = 1
    }
  }

  password_policy {
    minimum_length    = 8
    require_lowercase = true
    require_numbers   = true
    require_symbols   = true
    require_uppercase = true
  }

  lifecycle {
    prevent_destroy = false
  }
}

# Client configuration to communicate with user pool
resource "aws_cognito_user_pool_client" "cms_app_client" {
  name         = "cms_client"
  user_pool_id = aws_cognito_user_pool.cms_user_pool.id

  generate_secret               = false
  explicit_auth_flows           = ["ALLOW_USER_PASSWORD_AUTH", "ALLOW_REFRESH_TOKEN_AUTH"]
  prevent_user_existence_errors = "ENABLED"
  supported_identity_providers  = ["COGNITO"]
  refresh_token_validity        = 30
  access_token_validity         = 60
  id_token_validity             = 60
  token_validity_units {
    access_token  = "minutes"
    id_token      = "minutes"
    refresh_token = "days"
  }
}


# User Group in User pool configuration
resource "aws_cognito_user_group" "admin_group" {
  name         = "AdminGroup"
  user_pool_id = aws_cognito_user_pool.cms_user_pool.id
  description  = "CMS admin who manage platform setting"
}

resource "aws_cognito_user_group" "owner_group" {
  name         = "SubSystemOwnerGroup"
  user_pool_id = aws_cognito_user_pool.cms_user_pool.id
  description  = "Page Owners who manage their own content and team"
}
resource "aws_cognito_user_group" "staff_group" {
  name         = "StaffGroup"
  user_pool_id = aws_cognito_user_pool.cms_user_pool.id
  description  = "Staff members with limited access"
}

output "user_pool_id" {
  value = aws_cognito_user_pool.cms_user_pool.id
}

output "user_pool_client_id" {
  value = aws_cognito_user_pool_client.cms_app_client.id
}