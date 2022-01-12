package cmd

import "fmt"

const WEBSOCKET_MODULE_MESSAGE = `
module "websocket" {
	source                               = "./modules/websocket"
	ENVIRONMENT                          = var.ENVIRONMENT
	PROJECT_NAME                         = var.PROJECT_NAME
	AWS_REGION                           = var.AWS_REGION
	AWS_ACCOUNT_ID                       = var.AWS_ACCOUNT_ID
	WEBSOCKET_AUTHORIZATION_SECRET_VALUE = var.WEBSOCKET_AUTHORIZATION_SECRET_VALUE
}`

const VARIABLES_MESSAGE = `
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

const UNACCEPTED_ENVIRONMENT_MESSAGE = `
Environment name does not match any of the environments in the WebSocket generator config file.
Please try again using one of the environments you set when running websocket-generator %s.
`

const PROJECT_ENVIRONMENT_MISSING_MESSAGE = `
Project and Environment name are required,
i.e. 'websocket-generator %s <project-name> <environment>'.
`

const UNACCEPTED_PROJECT_NAME_MESSAGE = `
Project name does not match the project name in the WebSocket generator config file.
Did you mean to run websocket-generator %s?
`

const WEBSOCKET_GENERATOR_FILE_MESSAGE =
`// This is to stop Serverless complaining that the confirguration file is not in the same directory as serverless file.
const websocketGeneratorConfig = require('%s');
module.exports = websocketGeneratorConfig;
`

var CONFIG_FILE_NOT_FOUND_MESSAGE = fmt.Sprintf("The file %s could not be found.", WEBSOCKET_CONFIG_FILE_PATH)

var CURRENT_DIRECTORY_ERROR_MESSAGE = "Error currentDirectory: The current directory path was not retrieved"
