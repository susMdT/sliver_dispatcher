package tui

import (
	"strings"

	"sliver-dispatch/tui/cmd"
	"sliver-dispatch/utils"

	"github.com/c-bata/go-prompt"
	cobraprompt "github.com/stromland/cobra-prompt"
)

var advancedPrompt = &cobraprompt.CobraPrompt{
	RootCmd:                  cmd.RootCmd,
	PersistFlagValues:        false,
	ShowHelpCommandAndFlags:  true,
	DisableCompletionCommand: true,
	AddDefaultExitCommand:    true,
	GoPromptOptions: []prompt.Option{
		prompt.OptionTitle("sliver-dispatcher"),
		prompt.OptionPrefix("[>] "),
		prompt.OptionMaxSuggestion(10),
	},
	DynamicSuggestionsFunc: func(main_cmd string, document *prompt.Document) []prompt.Suggest {
		if suggestions := cmd.GetMainCmd(main_cmd); suggestions != nil {
			return suggestions
		}
		return []prompt.Suggest{}
	},
	OnErrorFunc: func(err error) {
		if strings.Contains(err.Error(), "unknown command") {
			utils.Eprint(err.Error())
			return
		}

		utils.Eprint(err.Error())
	},
}

func Main() {
	advancedPrompt.Run()
}
