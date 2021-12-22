package cmd

import (
	"fmt"
	"os"
)


func GetQuestionDetails(reference string) QuestionDetails {
	currentDirectory, error := os.Getwd()
	if error != nil {
		fmt.Println(error)
	}

	QuestionsConfig := map[string]QuestionDetails{
		"InfrastructureFilePath":  {
			QuestionLabel: "Infrastructure Code Location Path. Default value is <current-directory>/infrastructure",
			DefaultResponse: currentDirectory + "/infrastructure",
			ResponseLabel: "Your Infrastructure Code Location Path:",
		},
		"WebsocketFilePath":  {
			QuestionLabel: "WebSocket Code Location Path. Default value is <current-directory>/websocket",
			DefaultResponse: currentDirectory + "/websocket",
			ResponseLabel: "Your WebSocket Code Location Path:",
		},
		"AWSRegion":  {
			QuestionLabel: "AWS Region. Default value is eu-west-2",
			DefaultResponse: "eu-west-2",
			ResponseLabel: "Your AWS Region:",
		},
		"AuthorizationKey":  {
			QuestionLabel: "Authorization Key Name. Default value is Authorization",
			DefaultResponse: "Authorization",
			ResponseLabel: "Your Authorization Query Parameter Key is:",
		},
		"Environment":  {
			QuestionLabel: "Environment Name. Default value is Development",
			DefaultResponse: "Development",
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
