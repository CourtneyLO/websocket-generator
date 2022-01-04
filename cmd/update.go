package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Deploys all code changes for WebSocket API",
	Long: `This command will run terraform init, terraform workspace select, terraform apply and serverless deployment.
It will update the backend infrastructure and the CloudFormation stack and lambdas in AWS.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		configFile := ReadFile(WEBSOCKET_CONFIG_FILE_PATH)

		if len(configFile) == 0 {
			errorMessage.Println(CONFIG_FILE_NOT_FOUND_MESSAGE)
			return
		}

		argumentsValid, argumentInvalidMessage, argumentsForMessage := checkForValidArguments("update", args, configFile)

		if !argumentsValid {
			errorMessage.Printf(argumentInvalidMessage, strings.Join(argumentsForMessage, " "))
			return
		}

		currentDirectory, error := os.Getwd()
		if error != nil {
			errorMessage.Println("ERROR: The current directory path was not retrieved: %v", error)
			return
		}

		projectName := strings.ToLower(args[0])
		environment := strings.ToLower(args[1])

		InitTerraform(configFile, currentDirectory, projectName, environment)
		SelectWorkSpaceTerraform(configFile, currentDirectory, environment)
		ApplyTerraform(configFile, currentDirectory, projectName, environment)
		DeployServerless(configFile, currentDirectory, environment)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
