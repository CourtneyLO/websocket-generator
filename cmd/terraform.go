package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"bytes"

	"github.com/spf13/cobra"
)

var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Handles Terraform commands only",
	Long: `If you wish to only apply or destroy infrastructure code without deploying/removing serveless infrastructure use this command with either apply or delete.
Note: if you choose to destroy the infrastructure but keep the serveless functionlity you may experience some issues.
It is recommended to destroy the infrastructure first.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please add apply, destroy or --help to your command")
	},
}

func init() {
	rootCmd.AddCommand(terraformCmd)
}

func terraformExecCommand(action string, configFile map[string]interface{}, currentDirectory string, projectName string, environment string) {
	environmentVariable := fmt.Sprintf("TF_VAR_ENVIRONMENT=%v", environment)
	awsAccountIdVariable := fmt.Sprintf("TF_VAR_AWS_ACCOUNT_ID=%v", configFile[environment].(map[string]interface{})["awsAccountId"])
	AwsRegionVariable := fmt.Sprintf("TF_VAR_AWS_REGION=%v", configFile["awsRegion"])
	projectNameVariable := fmt.Sprintf("TF_VAR_PROJECT_NAME=%v", projectName)
	infrastructureFilePath := configFile["infrastructureFilePath"]

	if infrastructureFilePath == nil {
		errorMessage.Println("ERROR: Infrastructure file path could not be found")
		return
	}

	directory := currentDirectory + fmt.Sprintf("%v", infrastructureFilePath)

	command := exec.Command("terraform", action)
	command.Env = os.Environ()
	command.Env = append(command.Env, environmentVariable)
	command.Env = append(command.Env, awsAccountIdVariable)
	command.Env = append(command.Env, AwsRegionVariable )
	command.Env = append(command.Env, projectNameVariable)
	command.Dir = directory
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Run()
}

func workspaceCommands(action string, configFile map[string]interface{}, currentDirectory string, environment string) {
	command := exec.Command("terraform", "workspace", action, environment)

	var errb bytes.Buffer
	command.Dir = fmt.Sprintf("%s%v", currentDirectory, configFile["infrastructureFilePath"])
	command.Stderr = &errb
	command.Run()

	if errb.Len() > 0 {
		workspaceCommands("new", configFile, currentDirectory, environment)
	}
}

func ApplyTerraform(configFile map[string]interface{}, currentDirectory string, projectName string, environment string) {
	terraformExecCommand("apply", configFile, currentDirectory, projectName, environment)
}

func DestroyTerraform(configFile map[string]interface{}, currentDirectory string, projectName string, environment string) {
	terraformExecCommand("destroy", configFile, currentDirectory, projectName, environment)
}

func InitTerraform(configFile map[string]interface{}, currentDirectory string, projectName string, environment string) {
	terraformExecCommand("init", configFile, currentDirectory, projectName, environment)
}

func SelectWorkSpaceTerraform(configFile map[string]interface{}, currentDirectory string, environment string) {
	workspaceCommands("select", configFile, currentDirectory, environment)
}

func DeleteWorkSpaceTerraform(configFile map[string]interface{}, currentDirectory string, projectName string, environment string) {
	// Move to default directory before trying to delete current terraform environment workspace
	SelectWorkSpaceTerraform(configFile, currentDirectory, "default")
	workspaceCommands("delete", configFile, currentDirectory, environment)
}

func PlanTerraform(configFile map[string]interface{}, currentDirectory string, projectName string, environment string) {
	terraformExecCommand("plan", configFile, currentDirectory, projectName, environment)
}
