package commands

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"bitfan/cmd/bitfanUI/server"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "bitfan-ui",
	Short: "Bitfan Web UI",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("dev", cmd.Flags().Lookup("dev"))
		viper.BindPFlag("host", cmd.Flags().Lookup("host"))
		viper.BindPFlag("api", cmd.Flags().Lookup("api"))
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		httpServerMux := http.NewServeMux()
		httpServerMux.Handle("/", server.Handler(
			viper.GetString("api"),
			viper.GetBool("dev"),
		))

		addr := viper.GetString("host")
		fmt.Printf("serving on http://%s\n", addr)
		http.ListenAndServe(addr, httpServerMux)
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default : /etc/bitfan/bitfan-ui.toml, $HOME/.bitfan/bitfan-ui.toml, ./bitfan-ui.toml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	RootCmd.PersistentFlags().Bool("dev", false, "dev mode (serve asset and templates from disk")
	RootCmd.PersistentFlags().StringP("host", "H", "127.0.0.1:8081", "Serve UI on Host")
	RootCmd.PersistentFlags().StringP("api", "a", "127.0.0.1:5123", "Bitfan API to connect to")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigName("bitfan-ui") // name of config file (without extension)

	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("/etc/bitfan/")  // path to look for the config file in
		viper.AddConfigPath("$HOME/.bitfan") // call multiple times to add many search paths
		viper.AddConfigPath(".")             // optionally look for config in the working directory
	}

	viper.SetEnvPrefix("bitfanui")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
