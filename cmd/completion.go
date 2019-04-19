package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"io"
	"os"
)

const (
	example = `  # Load k8s-login completion code for bash into the current shell:
  source <(k8s-login completion bash)

  # Load k8s-login completion code for bash on startup:
  k8s-login completion bash > ~/.k8s-login/completion.bash.inc
  printf "
    # k8s-login shell completion
	source '$HOME/.k8s-login/completion.bash.inc'
  " >> $HOME/.bash_profile
  source $HOME/.bash_profile`
)

var (
	shellCompletions = map[string]func(writer io.Writer) error{
		"bash": rootCommand.GenBashCompletion,
	}
)

func init() {
	shells := []string{}
	for shell := range shellCompletions {
		shells = append(shells, shell)
	}

	rootCommand.AddCommand(&cobra.Command{
		Use:                   "completion bash",
		DisableFlagsInUseLine: true,
		Short:                 "Generates shell completion for the specified shell (bash only at the moment)",
		Example:               example,
		ValidArgs:             shells,
		RunE: func(command *cobra.Command, args []string) error {
			return completion(command, args)
		},
	})
}

func completion(command *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("shell not specified")
	}
	if len(args) > 1 {
		return errors.New("too many arguments, expected only the shell type")
	}

	function, found := shellCompletions[args[0]]
	if !found {
		return errors.New("unknown shell '" + args[0] + "'")
	}

	return function(os.Stdout)
}
