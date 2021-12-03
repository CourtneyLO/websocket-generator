package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Shows the changes to all terraform code to be used to create AWS resources",
	Long: `This command will run terraform init, terraform workspace new and terraform plan`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println(`
				Project and Environment name are required,
				i.e. 'websocket-generator terraform apply helloworld development'
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
		CreateWorkSpaceTerraform(webSocketConfig.InfrastructureFilePath, environment)
		PlanTerraform(webSocketConfig.InfrastructureFilePath, environment)
	},
}

func init() {
	terraformCmd.AddCommand(planCmd)
}
