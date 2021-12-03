package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Deploys all code for creating a WebSocket API",
	Long: `This command will run terraform init, terraform workspace select, terraform apply and serverless deployment.
It will create the backend infrastructure and the CloudFormation stack, lambdas and logs in AWS.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println(`
				Project and Environment name are required,
				i.e. 'websocket-generator create helloworld development'
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
		InitTerraform(webSocketConfig.InfrastructureFilePath, environment)
		SelectWorkSpaceTerraform(webSocketConfig.InfrastructureFilePath, environment)
		ApplyTerraform(webSocketConfig.InfrastructureFilePath, environment)
		DeploySeverless(webSocketConfig.WebsocketFilePath, environment)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
