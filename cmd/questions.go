package cmd

import (
	"fmt"
	"path/filepath"
	"os"
)

func environmentQuestions(projectName string, environmentConfig map[string]EnvironmentData, environments []string) []string {
	environmentLable := "Environment Name. Default value is " + DEFAULT_ENVIRONMENT
	choosenEnvironment := prompt(environmentLable, DEFAULT_ENVIRONMENT)
	fmt.Printf("The environment you choose is %q\n", choosenEnvironment)

	awsAccountID := requiredPrompt("AWS Account ID")
	fmt.Printf("Your AWS Account ID %q\n", awsAccountID)

	environmentData := EnvironmentData {
		Environment: choosenEnvironment,
		AWSAccountID: awsAccountID,
	}

	environmentConfig[choosenEnvironment] = environmentData
	environments = append(environments, choosenEnvironment)

	addAnotherEnvironment, addAnotherEnvironmentInput := yesNo("Would you like to add another environment (recommended if your AWS account ID differs per environment)")
	fmt.Println("You chose " + addAnotherEnvironmentInput + " to adding another environment")

	if addAnotherEnvironment {
		return environmentQuestions(projectName, environmentConfig, environments)
	}

	return environments
}

func Questions(projectName string) (WebSocketConfig, map[string]EnvironmentData) {
	currentDirectory, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	choosenLanguage := language()
	fmt.Printf("Language choosen %q\n", choosenLanguage)

	infrastructureFilePathLabel := "Infrastructure Code Location Path. Default value is <current-directory>/infrastructure"
	infrastructureFilePathEnding := prompt(infrastructureFilePathLabel,  "infrastructure")

	// MUST DELETE
	currentDirectory = "/Users/osborncourtney/Development/DV/own-work/test"
	// MUST DELETE

	infrastructureDestinationFilePath := filepath.Join(currentDirectory, infrastructureFilePathEnding)
	fmt.Printf("Your Infrastructure Code Location Path %q\n", infrastructureDestinationFilePath)

	websocketFilePathLabel := "WebSocket Code Location Path. Default value is <current-directory>/websocket"
	websocketFilePathEnding := prompt(websocketFilePathLabel, "websocket")
	websocketDestinationFilePath := filepath.Join(currentDirectory, websocketFilePathEnding)
	fmt.Printf("Your WebSocket Code Location Path %q\n", websocketDestinationFilePath)

	awsRegionLabel := "AWS Region. Default value is " + DEFAULT_AWS_REGION
	awsRegion := prompt(awsRegionLabel, DEFAULT_AWS_REGION)
	fmt.Printf("Your AWS Region %q\n", awsRegion)

	authorizationKeyLabel := "Authorization Key Name. Default value is " + DEFAULT_AUTHORIZATION_KEY_NAME
	authorizationKey := prompt(authorizationKeyLabel, DEFAULT_AUTHORIZATION_KEY_NAME)
	fmt.Printf("Your Serverless Authorization key is: %q\n", authorizationKey)

	var environments []string
	environmentConfig := make(map[string]EnvironmentData)
	choosenEnvironments := environmentQuestions(projectName, environmentConfig, environments)

	existingEnvironments := make(map[string]bool)
	var uniqueEnvironments []string
	for _, environment := range choosenEnvironments {
    _, environmentAlreadyExists := existingEnvironments[environment]

		if !environmentAlreadyExists {
      existingEnvironments[environment] = true
      uniqueEnvironments = append(uniqueEnvironments, environment)
    }
  }

	webSocketConfig := WebSocketConfig {
		Environments: uniqueEnvironments,
		ProjectName: projectName,
		Language: choosenLanguage,
		InfrastructureFilePath: infrastructureDestinationFilePath,
		WebsocketFilePath: websocketDestinationFilePath,
		AWSRegion: awsRegion,
		AuthorizationKey: authorizationKey,
	}

	return webSocketConfig, environmentConfig
}
