package cmd

func environmentQuestions(projectName string, environmentConfig map[string]EnvironmentData, environments []string) []string {
	choosenEnvironment := prompt("Environment")
	awsAccountID := requiredPrompt("AWSAccountID")

	environmentData := EnvironmentData {
		Environment: choosenEnvironment,
		AWSAccountID: awsAccountID,
	}

	environmentConfig[choosenEnvironment] = environmentData
	environments = append(environments, choosenEnvironment)

	addAnotherEnvironment := yesNo("Would you like to add another environment (recommended if your AWS account ID differs per environment)")

	if addAnotherEnvironment {
		return environmentQuestions(projectName, environmentConfig, environments)
	}

	return environments
}

func Questions(projectName string) (WebSocketConfig, map[string]EnvironmentData) {
	choosenLanguage := language()
	infrastructureDestinationFilePath := prompt("InfrastructureFilePath")

	websocketDestinationFilePath := prompt("WebsocketFilePath")
	awsRegion := prompt("AWSRegion")
	authorizationKey := prompt("AuthorizationKey")

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
