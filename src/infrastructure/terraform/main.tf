terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
}

provider "aws" {
  region = "eu-west-2"
}

module "websocket" {
  source      = "./modules/websocket"
  ENVIRONMENT = var.ENVIRONMENT
  PROJECT_NAME = var.PROJECT_NAME
  AWS_REGION = var.AWS_REGION
  AWS_ACCOUNT_ID = var.AWS_ACCOUNT_ID
  WEBSOCKET_AUTHORIZATION_SECRET_VALUE = var.WEBSOCKET_AUTHORIZATION_SECRET_VALUE
}
