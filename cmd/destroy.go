package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroys all AWS resources created by the apply/create",
	Long: `This command runs terraform workspace select, terraform destroy and terraform workspace delete`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println(`
				Project and Environment name are required,
				i.e. 'websocket-generator terraform destroy helloworld development'
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
		SelectWorkSpaceTerraform(webSocketConfig.InfrastructureFilePath, environment)
		DestroyTerraform(webSocketConfig.InfrastructureFilePath, environment)
		DeleteWorkSpaceTerraform(webSocketConfig.InfrastructureFilePath, environment)
	},
}

func init() {
	terraformCmd.AddCommand(destroyCmd)
}
