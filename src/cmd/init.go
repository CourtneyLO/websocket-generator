package cmd

import (
	"fmt"
	"strings"
	"go/build"
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
		if len(args) < 1 {
			fmt.Println("ERROR: Project Name is required to run this command")
			return
		}

		projectName := args[0]
		websocketConfig, environmentConfig := Questions(projectName)
		combinedConfig := structs.Map(websocketConfig)

		for key, value := range environmentConfig {
			combinedConfig[key] = value
		}

		WriteJsonFile(WEBSOCKET_CONFIG_FILE_PATH, combinedConfig)

		websocketGeneratorSrcLocation := getWebsocketGeneratorSrcLocation()

		currentDirectory, error := os.Getwd()
		if error != nil {
			fmt.Println(error)
		}

		constructInfrastructureDirectory(websocketConfig, currentDirectory, websocketGeneratorSrcLocation)
		constructServerlessDirectory(websocketConfig, currentDirectory, websocketGeneratorSrcLocation)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func constructInfrastructureDirectory(websocketConfig WebsocketConfig, currentDirectory string, websocketGeneratorSrcLocation string)  {


	// This is hard coded because this value will not changed in the near future.
	// Work  to make this dynamic will be done later
	choosenIaC := "terraform"

	sourceFileInfrastructure := websocketGeneratorSrcLocation + "infrastructure/" + choosenIaC

	destinationInfrastructureFilePath := currentDirectory + websocketConfig.InfrastructureFilePath
	modulesFolderExists := checkIfFileExists(destinationInfrastructureFilePath + "/modules")

	if modulesFolderExists {
		sourceFileInfrastructureError := CopyAndMoveFolder(sourceFileInfrastructure + "/modules", destinationInfrastructureFilePath + "/modules")
		if sourceFileInfrastructureError != nil {
			fmt.Println("The WebSocket modules folder failed to be copied and move to it's destination", sourceFileInfrastructureError)
		}
	} else {
		sourceFileInfrastructureError := CopyAndMoveFolder(sourceFileInfrastructure + "/modules", destinationInfrastructureFilePath + "/modules")
		if sourceFileInfrastructureError != nil {
			fmt.Println("The WebSocket modules folder failed to be copied and move to it's destination", sourceFileInfrastructureError)
		}
	}

	mainFileExists := checkIfFileExists(destinationInfrastructureFilePath + "/main.tf")

	if mainFileExists {
		fmt.Println("")
		informationHeading.Println("It appears you already have a main.tf file. Please add following code to your existing file:")
		informationDetails.Println(WEBSOCKET_MODULE_MESSAGE)
	} else {
		sourceFileInfrastructureError := CopyAndMoveFile(sourceFileInfrastructure + "/main.tf", destinationInfrastructureFilePath + "/main.tf")
		if sourceFileInfrastructureError != nil {
			fmt.Println("The WebSocket main.tf file failed to be copied and move to it's destination", sourceFileInfrastructureError)
		}
	}

	variableFileExists := checkIfFileExists(destinationInfrastructureFilePath + "/variables.tf")
	varsFileExists := checkIfFileExists(destinationInfrastructureFilePath + "/vars.tf")
	variableTypeFileExists := variableFileExists || varsFileExists

	if variableTypeFileExists {
		fmt.Println("")
		informationHeading.Println("It appears you already have a variables file. Please add following code to your existing file:")
		informationDetails.Println(VARIABLES_MESSAGE)
		fmt.Println("")
	} else {
		sourceFileInfrastructureError := CopyAndMoveFile(sourceFileInfrastructure + "/variables.tf", destinationInfrastructureFilePath + "/variables.tf")
		if sourceFileInfrastructureError != nil {
			fmt.Println("The WebSocket variables.tf file failed to be copied and move to it's destination", sourceFileInfrastructureError)
		}
	}
}

func constructServerlessDirectory(websocketConfig WebsocketConfig, currentDirectory string, websocketGeneratorSrcLocation string)  {
	sourceFileWebSockets := websocketGeneratorSrcLocation + "websockets/" + strings.ToLower(websocketConfig.Language)
	destinationWebsocketFilePath := currentDirectory + websocketConfig.WebsocketFilePath
	sourceFileWebSocketsError := CopyAndMoveFolder(sourceFileWebSockets, destinationWebsocketFilePath)

	if sourceFileWebSocketsError != nil {
		fmt.Println("The WebSocket folder failed to be copied and move to it's destination", sourceFileWebSocketsError)
		return
	}

	// Contructs config file to prevent Serverless complaining that the confirguration file is not in the same directory as serverless file.
	constructServerlessConfig(currentDirectory, websocketConfig.WebsocketFilePath)
}

func constructServerlessConfig(currentDirectory string, filePath string)  {
	filePathSplitIntoCharacters := strings.Split(filePath, "")

	constructedFilePath := ""
	for _, charater := range filePathSplitIntoCharacters {
		if charater == "/" {
			constructedFilePath = constructedFilePath + "../"
		}
	}

	constructedFilePath = constructedFilePath + WEBSOCKET_CONFIG_FILE_PATH
	data := []byte(fmt.Sprintf(WEBSOCKET_GENERATOR_FILE_MESSAGE, constructedFilePath))

	WriteFile(currentDirectory + filePath + "/config.js", data)
}

func getWebsocketGeneratorSrcLocation() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	return gopath + "/src/" + WEBSOCKET_GENERATOR_PACKAGE
}
