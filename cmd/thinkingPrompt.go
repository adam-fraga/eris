package cmd

import (
	"github.com/adam-fraga/eris/cmd/handler"
	"github.com/spf13/cobra"
)

var promptCmd = &cobra.Command{
	Use:   "think",
	Short: "Ask the LLM a prompt",
	Run: func(cmd *cobra.Command, args []string) {
		handler.RunThinkingPrompt()
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)
}
