package cmd

import (
	"fmt"
	"strings"
	"os"

	"github.com/fatih/structs"
	"github.com/spf13/cobra"
)

const DEFAULT_ENVIRONMENT = "development"
const DEFAULT_AWS_REGION = "eu-west-2"
const DEFAULT_AUTHORIZATION_KEY_NAME = "Authorization"

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates the config needed in order to create or delete the WebSocket API",
	Run: func(cmd *cobra.Command, args []string) {
		// This is hard coded because this value will not changed in the near future.
		// Work  to make this dynamic will be done later
		choosenIaC := "terraform"

		if len(args) < 1 {
			fmt.Println("ERROR: Project Name is required to run this command")
			return
		}

		projectName := args[0]


		webSocketConfig, environmentConfig := Questions(projectName)

		combinedConfig := structs.Map(webSocketConfig)

		for key, value := range environmentConfig {
			combinedConfig[key] = value
		}

		configFilePath := getFilePath(args)

		if configFilePath == "" {
			return
		}

		WriteFile(configFilePath, combinedConfig)

		sourceFileInfrastructure := "infrastructure/" + strings.ToLower(choosenIaC)

		infrastructureFilePath := webSocketConfig.InfrastructureFilePath

		modulesFolderExists := checkIfFileExists(infrastructureFilePath + "/modules")

		if modulesFolderExists {
			sourceFileInfrastructureError := CopyAndMoveFolder(sourceFileInfrastructure + "/modules", infrastructureFilePath + "/modules")
			if sourceFileInfrastructureError != nil {
				fmt.Println("The WebSocket modules folder failed to be copied and move to it's destination", sourceFileInfrastructureError)
			}
		} else {
			sourceFileInfrastructureError := CopyAndMoveFolder(sourceFileInfrastructure + "/modules", infrastructureFilePath + "/modules")
			if sourceFileInfrastructureError != nil {
				fmt.Println("The WebSocket modules folder failed to be copied and move to it's destination", sourceFileInfrastructureError)
			}
		}

		mainFileExists := checkIfFileExists(infrastructureFilePath + "/main.tf")

		if mainFileExists {
			fmt.Println("")
			fmt.Println("It appears you already have a main.tf file. Please add following code to your existing file:")
			fmt.Println(`
module "websocket" {
  source      = "./modules/websocket"
  ENVIRONMENT = var.ENVIRONMENT
  PROJECT_NAME = var.PROJECT_NAME
  AWS_REGION = var.AWS_REGION
  AWS_ACCOUNT_ID = var.AWS_ACCOUNT_ID
  WEBSOCKET_AUTHORIZATION_SECRET_VALUE = var.WEBSOCKET_AUTHORIZATION_SECRET_VALUE
}`)
		} else {
			sourceFileInfrastructureError := CopyAndMoveFile(sourceFileInfrastructure + "/main.tf", infrastructureFilePath + "/main.tf")
			if sourceFileInfrastructureError != nil {
				fmt.Println("The WebSocket main.tf file failed to be copied and move to it's destination", sourceFileInfrastructureError)
			}
		}

		variableFileExists := checkIfFileExists(infrastructureFilePath + "/variables.tf")
		varsFileExists := checkIfFileExists(infrastructureFilePath + "/vars.tf")
		variableTypeFileExists := variableFileExists || varsFileExists

		if variableTypeFileExists {
			fmt.Println("")
			fmt.Println("It appears you already have a variables file. Please add following code to your existing file:")
			fmt.Println(`
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
	type = string
	sensitive = true
}`)
		} else {
			sourceFileInfrastructureError := CopyAndMoveFile(sourceFileInfrastructure + "/variables.tf", infrastructureFilePath + "/variables.tf")
			if sourceFileInfrastructureError != nil {
				fmt.Println("The WebSocket variables.tf file failed to be copied and move to it's destination", sourceFileInfrastructureError)
			}
		}

		// Terraform config
		for _, environment := range webSocketConfig.Environments {
			terraformConfig := TerraformConfig {
				ENVIRONMENT: environment,
				AWS_REGION: webSocketConfig.AWSRegion,
				AWS_ACCOUNT_ID: environmentConfig[environment].AWSAccountID,
				PROJECT_NAME: webSocketConfig.ProjectName,
			}

			mappedTerraformConfig := structs.Map(terraformConfig)
			infrastructureConfigFileDirctory := webSocketConfig.InfrastructureFilePath + "/config/"

			error := os.MkdirAll(infrastructureConfigFileDirctory , os.ModePerm)

			if error != nil {
				fmt.Println("MkdirAll command for infrastructure config folder failed")
				return
			}

			terraformConfigFilePath := infrastructureConfigFileDirctory + environment + ".json"
			WriteFile(terraformConfigFilePath, mappedTerraformConfig)

			CreateWorkSpaceTerraform(webSocketConfig.InfrastructureFilePath, environment)
		}

		sourceFileWebsockets := "websockets/" + strings.ToLower(webSocketConfig.Language)
		sourceFileWebsocketsError := CopyAndMoveFolder(sourceFileWebsockets, webSocketConfig.WebsocketFilePath)

		if sourceFileWebsocketsError != nil {
			fmt.Println("The WebSocket folder failed to be copied and move to it's destination", sourceFileWebsocketsError)
			return
		}

		delete(combinedConfig, "Language")
		delete(combinedConfig, "WebsocketFilePath")
		delete(combinedConfig, "InfrastructureFilePath")

		WriteFile(webSocketConfig.WebsocketFilePath + "/config.json", combinedConfig)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
