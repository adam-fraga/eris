package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func showWelcome() {
	// Create a cyan color for the ASCII art
	cyan := color.New(color.FgCyan, color.Bold)
	// Green for the robot/instruction
	green := color.New(color.FgGreen, color.Bold)

	banner := `
  ______  _____   _____   _____  
 |  ____||  __ \ |_   _| / ____| 
 | |__   | |__) |  | |  | (___   
 |  __|  |  _  /   | |   \___ \  
 | |____ | | \ \  _| |_  ____) | 
 |______||_|  \_\|_____||_____/  
`

	cyan.Println(banner)
	green.Println("🤖 Hello! I am ERIS, your assistant.")
	green.Println("Use the 'eris prompt' command to interact with me.")
}

var rootCmd = &cobra.Command{
	Use:   "eris",
	Short: "Eris Agent",
	Long:  "Hi, Eris is a simple CLI tool that provide logics to llm to make your developer life easier.",
	Run: func(cmd *cobra.Command, args []string) {
		showWelcome()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func AddCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}
