/*
Copyright © 2024 Raznar Lab <xabhista19@raznar.id>
*/
package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"raznar.id/invoice-broker/config"
	"raznar.id/invoice-broker/internal/rest"
)

func startServer(configFile string) {
	if !strings.HasSuffix(configFile, ".yml") {
		log.Fatalf("The config file must be yml! instead of %s", configFile)
	}

	conf, err := config.New(configFile)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// fiber has built-in block, so we dont need any signal block
	if err = rest.Start(conf); err != nil {
		log.Fatalf("An error occured when starting the bot: %s", err.Error())
	}
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		configFile := cmd.Flag("config").Value.String()
		startServer(configFile)
	},
}

func init() {
	startCmd.Flags().String("config", "config.yml", "Configuration file")
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
