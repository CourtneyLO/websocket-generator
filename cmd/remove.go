package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Handles Remove Serverless command",
	Long: `This command will run serverless remove`,
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

		argumentsValid, argumentInvalidMessage, argumentsForMessage := checkForValidArguments("severless remove", args, configFile)

		if !argumentsValid {
			errorMessage.Printf(argumentInvalidMessage, strings.Join(argumentsForMessage, " "))
			return
		}

		currentDirectory, error := os.Getwd()
		if error != nil {
			errorMessage.Println(CURRENT_DIRECTORY_ERROR_MESSAGE, error)
			return
		}

		environment := strings.ToLower(args[1])

		RemoveServerless(configFile, currentDirectory, environment)
	},
}

func init() {
	serverlessCmd.AddCommand(removeCmd)
}
