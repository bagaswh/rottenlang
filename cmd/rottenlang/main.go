package main

import (
	"fmt"
	"io"
	"os"

	"github.com/bagaswh/rottenlang/pkg/rottenlang"
	"github.com/bagaswh/rottenlang/pkg/scanner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "app [file]",
	Short: "A simple application that processes a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]

		if _, err := os.Stat(filename); os.IsNotExist(err) {
			fmt.Printf("Error: File '%s' does not exist\n", filename)
			os.Exit(1)
		}

		f, err := os.OpenFile(filename, os.O_RDONLY, 0400)
		if err != nil {
			fmt.Printf("Error: Failed opening file '%s': %v", filename, err.Error())
			os.Exit(1)
		}
		source, err := io.ReadAll(f)
		if err != nil {
			fmt.Printf("Error: Failed reading file '%s': %v", filename, err.Error())
			os.Exit(1)
		}
		rottenlang := rottenlang.NewRottenlang(string(source), &scanner.StderrErrorReporter{})
		rottenlang.Scan()
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	// Example: rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.app.yaml)")
}

func initConfig() {
	viper.SetDefault("author", "Your Name")
	viper.SetDefault("license", "MIT")

	viper.AutomaticEnv()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
