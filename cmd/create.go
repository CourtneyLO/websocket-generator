package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Deploys all code for creating a WebSocket API",
	Long: `This command will run terraform init, terraform workspace create, terraform apply and serverless deployment.
It will create the backend infrastructure and the CloudFormation stack, lambdas and logs in AWS.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		configFile, readFileError := ReadFile(WEBSOCKET_CONFIG_FILE_PATH)

		if len(configFile) == 0 {
			errorMessage.Println(CONFIG_FILE_NOT_FOUND_MESSAGE)
			return
		}

		if readFileError != nil {
			errorMessage.Println(readFileError)
			return
		}

		argumentsValid, argumentInvalidMessage, argumentsForMessage := checkForValidArguments( "create", args, configFile)

		if !argumentsValid {
			errorMessage.Printf(argumentInvalidMessage, strings.Join(argumentsForMessage, " "))
			return
		}

		currentDirectory, error := os.Getwd()
		if error != nil {
			errorMessage.Println(CURRENT_DIRECTORY_ERROR_MESSAGE, error)
			return
		}

		projectName := strings.ToLower(args[0])
		environment := strings.ToLower(args[1])

		InitTerraform(configFile, currentDirectory, projectName, environment)
		SelectWorkSpaceTerraform(configFile, currentDirectory, environment)
		ApplyTerraform(configFile, currentDirectory, projectName, environment)
		InstallNodePackages(configFile, currentDirectory, environment)
		DeployServerless(configFile, currentDirectory, environment)
		PrintWebsocketUrl(configFile, environment)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
