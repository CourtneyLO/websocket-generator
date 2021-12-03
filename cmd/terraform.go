package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Handles Terraform commands only",
	Long: `If you wish to only apply/reapply/detroy infrastructure code without deploying/removing serveless infrastructure use this command with either apply or delete.
Note: if you choose to destroy the infrastructure but keep the serveless functionlity you may experience some issues.
It is recommended to destory the infrastructure first.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please add apply, destroy or --help to your command")
	},
}

func init() {
	rootCmd.AddCommand(terraformCmd)
}

func terraformExecCommand(action string, directory string, environment string) {
	command := exec.Command("terraform", action, "-var-file", "./config/" + environment + ".json")
	command.Dir = directory
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Run()
}

func workspaceCommands(action string, directory string, environment string)  {
	command := exec.Command("terraform", "workspace", action, environment)
	command.Dir = directory
	command.Run()
}

func ApplyTerraform(directory string, environment string) {
	terraformExecCommand("apply", directory, environment)
}

func DestroyTerraform(directory string, environment string) {
	terraformExecCommand("destroy", directory, environment)
}

func InitTerraform(directory string, environment string)  {
	terraformExecCommand("init", directory, environment)
}

func CreateWorkSpaceTerraform(directory string, environment string)  {
	workspaceCommands("new", directory, environment)
}

func SelectWorkSpaceTerraform(directory string, environment string)  {
	workspaceCommands("select", directory, environment)
}

func DeleteWorkSpaceTerraform(directory string, environment string)  {
	// Move to default directory before trying to delete current environments workspace
	SelectWorkSpaceTerraform(directory, "default")
	workspaceCommands("delete", directory, environment)
}

func PlanTerraform(directory string, environment string) {
	terraformExecCommand("plan", directory, environment)
}
