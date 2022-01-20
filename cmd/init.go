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
			errorMessage.Println("Error: Project Name is required to run this command")
			return
		}

		projectName := args[0]
		websocketConfig, environmentConfig := Questions(projectName)
		combinedConfig := structs.Map(websocketConfig)

		for key, value := range environmentConfig {
			combinedConfig[key] = value
		}

		writeJsonFileError := WriteJsonFile(WEBSOCKET_CONFIG_FILE_PATH, combinedConfig)

		if writeJsonFileError != nil {
			errorMessage.Println(writeJsonFileError)
		}

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
	// Work to make this dynamic will be done later
	choosenIaC := "terraform"
	sourceFilePath := websocketGeneratorSrcLocation + "infrastructure/" + choosenIaC
	destinationFilePath := currentDirectory + websocketConfig.InfrastructureFilePath

	copyFolderToRepo("/modules", sourceFilePath, destinationFilePath)
	copyFileToRepo("/.gitignore", sourceFilePath, destinationFilePath)
	copyFileToRepo("/main.tf", sourceFilePath, destinationFilePath)
	// Checks to see if either vars.tf of variables.tf already exist in user's repo before copong and adding the file
	copyFileToRepo("/vars.tf", sourceFilePath, destinationFilePath)
}

func copyFolderToRepo(filePath string, sourceFilePath string, destinationFilePath string)  {
	fileError := CopyAndMoveFolder(sourceFilePath + filePath, destinationFilePath + filePath)
	if fileError != nil {
		errorMessage.Println(fmt.Sprintf("Error: The infrastructure WebSocket %s file failed to be copied and move to it's destination", filePath), fileError)
	}
}

func copyFileToRepo(filePath string, sourceFilePath string, destinationFilePath string)  {
	fileExists, fileExistError := checkIfFileExists(destinationFilePath + filePath)

	if fileExistError != nil {
		errorMessage.Println(fileExistError)
		return
	}

	if !fileExists && filePath == "/vars.tf" {
		copyFileToRepo("/variables.tf", sourceFilePath, destinationFilePath)
		return
	}

	if !fileExists {
		fileCopyError := CopyAndMoveFile(sourceFilePath + filePath, destinationFilePath + filePath)
		if fileCopyError != nil {
			errorMessage.Println(fmt.Sprintf("Error: The infrastructure WebSocket %s file failed to be copied and move to it's destination", filePath), fileCopyError)
		}
		return
	}

	fmt.Println("")
	informationHeading.Println(fmt.Sprintf(FILE_ALREADY_EXISTS_MESSAGE, filePath))
	informationDetails.Println(getMessage(filePath))
	fmt.Println("")
}

func constructServerlessDirectory(websocketConfig WebsocketConfig, currentDirectory string, websocketGeneratorSrcLocation string)  {
	sourceFilePath := websocketGeneratorSrcLocation + "api/websockets/" + strings.ToLower(websocketConfig.Language)
	destinationFilePath := currentDirectory + websocketConfig.WebsocketFilePath
	copyFolderToRepo("/", sourceFilePath, destinationFilePath)

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

	writeFileError := WriteFile(currentDirectory + filePath + "/config.js", data)

	if writeFileError != nil {
		errorMessage.Println(writeFileError)
	}
}

func getWebsocketGeneratorSrcLocation() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	return gopath + "/src/" + WEBSOCKET_GENERATOR_PACKAGE
}

func getMessage(reference string) string {
	messages := map[string]string{
		"/main.tf": WEBSOCKET_MODULE_MESSAGE,
		"/variables.tf": VARIABLES_MESSAGE,
		"/vars.tf": VARIABLES_MESSAGE,
		"/.gitignore": GIT_IGNORE_MESSAGE,
	}
	return messages[reference]
}
