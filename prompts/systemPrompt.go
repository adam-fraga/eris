package prompts

import (
	"fmt"
	"strings"

	"github.com/adam-fraga/eris/capabilities"
)

func BuildSystemPrompt() string {
	var b strings.Builder

	b.WriteString("You are a helpful assistant with access to the following capabilities:\n")

	for _, cap := range capabilities.Capabilities {
		fmt.Fprintf(&b, "- %s: %s\n", cap.Name, cap.Description)
	}

	b.WriteString("When appropriate, call the capability instead of answering directly.")

	return b.String()
}
