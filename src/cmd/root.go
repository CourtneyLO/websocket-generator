package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "websocket-generator",
	Short: "WebSocket Generator will create everything you need to deploy a fully functional WebSocket API",
	Long: `WebSocket Generator adds the infrastructure and serverless code to your choosen directory and provides useful commands to deploying and deleting this code.
Run - websocket-generator init <projectName> - to get started
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please add --help to your command in order to find our what commands are avaliable to you.")
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".websocket-generator")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
