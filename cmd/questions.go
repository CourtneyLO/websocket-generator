package cmd

func environmentQuestions(projectName string, environmentConfig map[string]EnvironmentData, environments []string) []string {
	choosenEnvironment := prompt("Environment")
	awsAccountId := requiredPrompt("AWSAccountID")

	environmentData := EnvironmentData {
		Environment: choosenEnvironment,
		AwsAccountId: awsAccountId,
	}

	environmentConfig[choosenEnvironment] = environmentData
	environments = append(environments, choosenEnvironment)

	addAnotherEnvironment := yesNo("Would you like to add another environment (recommended if your AWS account ID differs per environment)")

	if addAnotherEnvironment {
		return environmentQuestions(projectName, environmentConfig, environments)
	}

	return environments
}

func Questions(projectName string) (WebsocketConfig, map[string]EnvironmentData) {
	choosenLanguage := language()
	infrastructureDestinationFilePath := prompt("InfrastructureFilePath")

	websocketDestinationFilePath := prompt("WebsocketFilePath")
	AwsRegion := prompt("AwsRegion")
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

	websocketConfig := WebsocketConfig {
		Environments: uniqueEnvironments,
		ProjectName: projectName,
		Language: choosenLanguage,
		InfrastructureFilePath: infrastructureDestinationFilePath,
		WebsocketFilePath: websocketDestinationFilePath,
		AwsRegion: AwsRegion,
		AuthorizationKey: authorizationKey,
	}

	return websocketConfig, environmentConfig
}

func GetQuestionDetails(reference string) QuestionDetails {
	QuestionsConfig := map[string]QuestionDetails{
		"InfrastructureFilePath":  {
			QuestionLabel: "Infrastructure Code Location Path. Default value is <current-directory>/infrastructure",
			DefaultResponse: "/infrastructure",
			ResponseLabel: "Your Infrastructure Code Location Path:",
		},
		"WebsocketFilePath":  {
			QuestionLabel: "WebSocket Code Location Path. Default value is <current-directory>/websocket",
			DefaultResponse: "/websocket",
			ResponseLabel: "Your WebSocket Code Location Path:",
		},
		"AwsRegion":  {
			QuestionLabel: "AWS Region. Default value is eu-west-2",
			DefaultResponse: "eu-west-2",
			ResponseLabel: "Your AWS Region:",
		},
		"AuthorizationKey":  {
			QuestionLabel: "Authorization Key Name. Default value is authorization",
			DefaultResponse: "authorization",
			ResponseLabel: "Your Authorization Query Parameter Key is:",
		},
		"Environment":  {
			QuestionLabel: "Environment Name. Default value is development",
			DefaultResponse: "development",
			ResponseLabel: "Environment Choosen:",
		},
		"AWSAccountID":  {
			QuestionLabel: "AWS Account ID",
			DefaultResponse: "",
			ResponseLabel: "Your AWS Account ID:",
		},
	}
	return QuestionsConfig[reference]
}
