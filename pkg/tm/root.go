package tm

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"knative.dev/client/pkg/kn/flags"

	"github.com/triggermesh/tm/pkg/tm/commands"
	"github.com/triggermesh/tm/pkg/tm/commands/version"
)

// NewRootCommand creates the root command for tm
func NewRootCommand() *cobra.Command {
	return newRootCommand()
}

// newRootCommand returns a root command. Params are used for testing
func newRootCommand(params ...commands.TmParams) *cobra.Command {
	var p *commands.TmParams
	if len(params) == 0 {
		p = &commands.TmParams{}
	} else if len(params) == 1 {
		p = &params[0]
	} else {
		panic("Too many params objects to newRootCommand")
	}
	p.Initialize()

	rootCmd := &cobra.Command{
		Use:   "tm",
		Short: "Triggermesh CLI",
		Long: `Manage TriggerMesh services to ease serverless events.

* Manage your serverless functions.
* Build container images.
* Integrate functions at Git repos.`,

		SilenceUsage: true,

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return flags.ReconcileBoolFlags(cmd.Flags())
		},
	}

	if p.Output != nil {
		rootCmd.SetOutput(p.Output)
	}

	// persistent flags
	rootCmd.PersistentFlags().StringVar(&p.KubeCfgPath, "kubeconfig", "", "kubectl config file (default is ~/.kube/config)")

	// root child commands
	rootCmd.AddCommand(version.NewVersionCommand(p))

	// Initialize default `help` cmd early to prevent unknown command errors
	rootCmd.InitDefaultHelpCmd()

	// Deal with empty and unknown sub command groups
	EmptyAndUnknownSubCommands(rootCmd)

	// Wrap usage.
	w, err := width()
	if err == nil {
		newUsage := strings.ReplaceAll(rootCmd.UsageTemplate(), "FlagUsages ",
			fmt.Sprintf("FlagUsagesWrapped %d ", w))
		rootCmd.SetUsageTemplate(newUsage)
	}

	// For glog parse error.
	flag.CommandLine.Parse([]string{})

	return rootCmd
}

// EmptyAndUnknownSubCommands adds a RunE to all commands that are groups to
// deal with errors when called with empty or unknown sub command
func EmptyAndUnknownSubCommands(cmd *cobra.Command) {
	for _, childCmd := range cmd.Commands() {
		if childCmd.HasSubCommands() && childCmd.RunE == nil {
			childCmd.RunE = func(aCmd *cobra.Command, args []string) error {
				aCmd.Help()
				if len(args) == 0 {
					return fmt.Errorf("please provide a valid sub-command for \"%s\"", aCmd.Name())
				}
				return fmt.Errorf("unknown sub-command \"%s\" for \"%s\"", args[0], aCmd.Name())
			}
		}

		// recurse to deal with child commands that are themselves command groups
		EmptyAndUnknownSubCommands(childCmd)
	}
}

func width() (int, error) {
	width, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	return width, err
}
