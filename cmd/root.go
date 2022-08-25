package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var Verbose bool = false

var RootCmd = &cobra.Command{
	Use:   "n2t",
	Short: "n2t is a helper application for Nand2Tetris",
	Long:  `See https://github.com/trammell/n2t`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func setUpLogging(cmd *cobra.Command, args []string) {
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: `>>>`}
	output.FormatLevel = func(i interface{}) string {
		return ``
	}
	log.Logger = log.Output(output)
	zerolog.SetGlobalLevel(zerolog.Disabled)
	if Verbose {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}
