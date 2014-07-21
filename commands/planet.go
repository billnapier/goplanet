package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "goplanet",
	Short: `planet style RSS aggregator`,
	Long:  `Provides a planet style RSS aggregator`,
	Run:   rootRun,
}

func rootRun(cmd *cobra.Command, args []string) {
	fmt.Println(viper.GetString("appname"))
	fmt.Println(viper.Get("feeds"))
}

func Execute() {
	addCommands()
	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

var CfgFile string

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&CfgFile, "config", "",
		"config file (default is $HOME/.goplanet/config.yaml)")
}

func initConfig() {
	if CfgFile != "" {
		viper.SetConfigFile(CfgFile)
	}
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/goplanet")
	viper.AddConfigPath("$HOME/.goplanet")
	viper.ReadInConfig()
}

func addCommands() {
	RootCmd.AddCommand(fetchCmd)
}
