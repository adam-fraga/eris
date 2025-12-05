package cmd

import (
	"github.com/adam-fraga/eris/cmd/handler"
	"github.com/spf13/cobra"
)

var thinkPomptCmd = &cobra.Command{
	Use:   "think",
	Short: "Ask the LLM a thinking prompt",
	Run: func(cmd *cobra.Command, args []string) {
		handler.RunThinkingPrompt()
	},
}

var codePomptCmd = &cobra.Command{
	Use:   "code",
	Short: "Ask the LLM a coding prompt",
	Run: func(cmd *cobra.Command, args []string) {
		handler.RunCodingPrompt()
	},
}

func init() {
	rootCmd.AddCommand(thinkPomptCmd)
	rootCmd.AddCommand(codePomptCmd)

}
