package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes all code for creating a WebSocket API",
	Long: `This command will run serverless remove, terraform workspace select, terraform destroy and terraform workspace delete.
It will delete all the backend infrastructure and the CloudFormation stack, lambdas and logs in AWS.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		configFile := ReadFile(WEBSOCKET_CONFIG_FILE_PATH)

		if len(configFile) == 0 {
			errorMessage.Println(CONFIG_FILE_NOT_FOUND_MESSAGE)
			return
		}

		argumentsValid, argumentInvalidMessage, argumentsForMessage := checkForValidArguments("delete", args, configFile)

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

		RemoveServerless(configFile, currentDirectory, environment)
		SelectWorkSpaceTerraform(configFile, currentDirectory, environment)
		DestroyTerraform(configFile, currentDirectory, projectName, environment)
		DeleteWorkSpaceTerraform(configFile, currentDirectory, projectName, environment)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
