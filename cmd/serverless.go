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

func severlessExecCommand(action string, directory string, environment string) {
	command := exec.Command("sls", action, "--stage", environment)
	command.Dir = directory
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Run()
}

func DeploySeverless(directory string, environment string) {
	severlessExecCommand("deploy", directory, environment)
}

func RemoveSeverless(directory string, environment string) {
	severlessExecCommand("remove", directory, environment)
}
