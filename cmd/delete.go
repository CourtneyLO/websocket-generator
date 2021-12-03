package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes all code for creating a WebSocket API",
	Long: `This command will run serverless remove, terraform workspace select, terraform destroy and terraform workspace delete.
It will delete all the backend infrastructure and the CloudFormation stack, lambdas and logs in AWS.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println(`
				Project and Environment name are required,
				i.e. 'websocket-generator delete helloworld development'
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
		RemoveSeverless(webSocketConfig.WebsocketFilePath, environment)
		SelectWorkSpaceTerraform(webSocketConfig.InfrastructureFilePath, environment)
		DestroyTerraform(webSocketConfig.InfrastructureFilePath, environment)
		DeleteWorkSpaceTerraform(webSocketConfig.InfrastructureFilePath, environment)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
