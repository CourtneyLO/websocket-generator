package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroys all AWS resources created by the apply/create",
	Long: `This command runs terraform workspace select, terraform destroy and terraform workspace delete`,
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

		argumentsValid, argumentInvalidMessage, argumentsForMessage := checkForValidArguments("terraform destroy", args, configFile)

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

		SelectWorkSpaceTerraform(configFile, currentDirectory, environment)
		DestroyTerraform(configFile, currentDirectory, projectName, environment)
		DeleteWorkSpaceTerraform(configFile, currentDirectory, projectName, environment)
	},
}

func init() {
	terraformCmd.AddCommand(destroyCmd)
}
