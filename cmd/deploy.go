package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Handles Serverless Deployment",
	Long: `This command will run severless deploy`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println(`
				Project and Environment name are required,
				i.e. 'websocket-generator serveless deploy helloworld development'
			`)
			return
		}

		configFilePath := getFilePath(args)
		configFileExists := checkIfFileExists(configFilePath)

		if !configFileExists {
			return
		}

		environment := args[1]
		webSocketConfig := ReadFile(configFilePath)
		DeploySeverless(webSocketConfig.WebsocketFilePath, environment)
	},
}

func init() {
	serverlessCmd.AddCommand(deployCmd)
}
