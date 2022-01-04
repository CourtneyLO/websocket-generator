package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies all terraform code to create AWS resources",
	Long: `This command will run terraform init, terraform workspace new and terraform apply`,
	Run: func(cmd *cobra.Command, args []string) {
		configFile := ReadFile(WEBSOCKET_CONFIG_FILE_PATH)

		if len(configFile) == 0 {
			errorMessage.Println(CONFIG_FILE_NOT_FOUND_MESSAGE)
			return
		}

		argumentsValid, argumentInvalidMessage, argumentsForMessage := checkForValidArguments("apply", args, configFile)

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
		CreateWorkSpaceTerraform(configFile, currentDirectory, environment)
		ApplyTerraform(configFile, currentDirectory, projectName, environment)
	},
}

func init() {
	terraformCmd.AddCommand(applyCmd)
}
