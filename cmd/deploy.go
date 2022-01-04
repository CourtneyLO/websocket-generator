package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Handles Serverless Deployment",
	Long: `This command will run serverless deploy`,
	Run: func(cmd *cobra.Command, args []string) {
		configFile := ReadFile(WEBSOCKET_CONFIG_FILE_PATH)

		if len(configFile) == 0 {
			errorMessage.Println(CONFIG_FILE_NOT_FOUND_MESSAGE)
			return
		}

		argumentsValid, argumentInvalidMessage, argumentsForMessage := checkForValidArguments("serverless deploy", args, configFile)

		if !argumentsValid {
			errorMessage.Printf(argumentInvalidMessage, strings.Join(argumentsForMessage, " "))
			return
		}

		currentDirectory, error := os.Getwd()
		if error != nil {
			errorMessage.Println("ERROR: The current directory path was not retrieved: %v", error)
			return
		}

		environment := strings.ToLower(args[1])

		DeployServerless(configFile, currentDirectory, environment)
	},
}

func init() {
	serverlessCmd.AddCommand(deployCmd)
}
