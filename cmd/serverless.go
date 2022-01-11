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

func serverlessExecCommand(action string, configFile map[string]interface{}, currentDirectory string, environment string) {
	websocketFilePath := configFile["websocketFilePath"]

	if websocketFilePath == nil {
		errorMessage.Println("ERROR: WebSocket file path could not be found")
		return
	}

	directory := fmt.Sprintf("%s%v", currentDirectory, websocketFilePath)

	command := exec.Command("sls", action, "--stage", environment)
	command.Dir = directory
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Run()
}

func npmInstallExecCommand(configFile map[string]interface{}, currentDirectory string, environment string) {
	websocketFilePath := configFile["websocketFilePath"]

	if websocketFilePath == nil {
		errorMessage.Println("ERROR: WebSocket file path could not be found")
		return
	}

	directory := fmt.Sprintf("%s%v", currentDirectory, websocketFilePath)

	command := exec.Command("nmp", "install")
	command.Dir = directory
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Run()
}

func DeployServerless(configFile map[string]interface{}, currentDirectory, environment string){
	serverlessExecCommand("deploy", configFile, currentDirectory, environment)
}

func RemoveServerless(configFile map[string]interface{}, currentDirectory, environment string) {
	serverlessExecCommand("remove", configFile, currentDirectory, environment)
}

func InstallNodePackages(configFile map[string]interface{}, currentDirectory, environment string)  {
	npmInstallExecCommand(configFile, currentDirectory, environment)
}
