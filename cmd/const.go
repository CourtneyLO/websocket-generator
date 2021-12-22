package cmd

const WEBSOCKET_CONFIG_FILE_PATH = "websocket-generator-config-files/"

const WEBSOCKET_MODULE = `
module "websocket" {
	source                               = "./modules/websocket"
	ENVIRONMENT                          = var.ENVIRONMENT
	PROJECT_NAME                         = var.PROJECT_NAME
	AWS_REGION                           = var.AWS_REGION
	AWS_ACCOUNT_ID                       = var.AWS_ACCOUNT_ID
	WEBSOCKET_AUTHORIZATION_SECRET_VALUE = var.WEBSOCKET_AUTHORIZATION_SECRET_VALUE
}`

const VARIABLES = `
variable "ENVIRONMENT" {
	type = string
}

variable "PROJECT_NAME" {
	type = string
}

variable "AWS_REGION" {
	type = string
}

variable "AWS_ACCOUNT_ID" {
	type = string
}

variable "WEBSOCKET_AUTHORIZATION_SECRET_VALUE" {
	type      = string
	sensitive = true
}`
