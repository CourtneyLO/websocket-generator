package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var serverlessCmd = &cobra.Command{
	Use:   "serverless",
	Short: "Handles Serverless commands only",
	Long: `If you wish to only deploy/redeploy/remove serverless code without affecting the infrastructure use this command with either deploy or remove.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please add deploy, remove or --help to your command")
	},
}

func init() {
	rootCmd.AddCommand(serverlessCmd)
}

func executeCommand(commandType string, configFile map[string]interface{}, currentDirectory string, environment string, action string)  {
	websocketFilePath := configFile["websocketFilePath"]

	if websocketFilePath == nil {
		errorMessage.Println("ERROR: WebSocket file path could not be found")
		return
	}

	directory := fmt.Sprintf("%s%v", currentDirectory, websocketFilePath)

	command := getExecuteCommand(commandType, action, environment)
	command.Dir = directory
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Run()
}

func serverlessExecCommand(configFile map[string]interface{}, currentDirectory string, environment string, action string) {
	executeCommand("serverless", configFile, currentDirectory, environment, action)
}

func npmInstallExecCommand(configFile map[string]interface{}, currentDirectory string, environment string) {
	executeCommand("npm", configFile, currentDirectory, environment, "")
}

func DeployServerless(configFile map[string]interface{}, currentDirectory, environment string){
	serverlessExecCommand(configFile, currentDirectory, environment, "deploy")
}

func RemoveServerless(configFile map[string]interface{}, currentDirectory, environment string) {
	serverlessExecCommand(configFile, currentDirectory, environment, "remove")
}

func InstallNodePackages(configFile map[string]interface{}, currentDirectory, environment string)  {
	npmInstallExecCommand(configFile, currentDirectory, environment)
}

func getExecuteCommand(commandType string, action string, environment string) *exec.Cmd {
	switch commandType {
	case "serverless":
		return exec.Command("sls", action, "--stage", environment)
	case "npm":
		return exec.Command("npm", "install")
	}

	return exec.Command("sls", action, "--stage", environment);
}
