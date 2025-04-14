package main

import (
	"fmt"

	elkcollector "nhandd/bego/cmd/elkCollector"
	gormgenerator "nhandd/bego/cmd/gormGenerator"
	"nhandd/bego/cmd/simulator"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "This is a CLI app",
		Long:  "A simple CLI application built with Cobra",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello from the root command!")
		},
	}

	var runElkWorkerCmd = &cobra.Command{
		Use:   "run-elk-worker",
		Short: "Run the monitor",
		Long:  "This command runs the monitor for the application",
		Run: func(cmd *cobra.Command, args []string) {
			elkcollector.RunCollector()
		},
	}

	var runSimulatorCmd = &cobra.Command{
		Use:   "run-simulation",
		Short: "Run the web server",
		Long:  "This command runs the web server for the application",
		Run: func(cmd *cobra.Command, args []string) {
			simulator.RunSimulation()
		},
	}

	var runGormGenerator = &cobra.Command{
		Use:   "run-generator",
		Short: "Run the web server",
		Long:  "This command runs the web server for the application",
		Run: func(cmd *cobra.Command, args []string) {
			gormgenerator.RunDbModelsGenerators()
		},
	}

	rootCmd.AddCommand(runElkWorkerCmd)
	rootCmd.AddCommand(runSimulatorCmd)
	rootCmd.AddCommand(runGormGenerator)
	rootCmd.Execute()
}
