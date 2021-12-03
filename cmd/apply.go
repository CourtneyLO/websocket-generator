package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies all terraform code to create AWS resources",
	Long: `This command will run terraform init, terraform workspace new and terraform apply`,
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
		ApplyTerraform(webSocketConfig.InfrastructureFilePath, environment)
	},
}

func init() {
	terraformCmd.AddCommand(applyCmd)
}
