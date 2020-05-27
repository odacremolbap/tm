package tm

import (
	"github.com/spf13/cobra"
	"github.com/triggermesh/tm/pkg/tm/commands"
	"github.com/triggermesh/tm/pkg/tm/commands/version"
)

// NewRootCommand creates the root command for tm
func NewRootCommand() *cobra.Command {
	return newRootCommand()
	// root := newRootCommand() {

	// }
}

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

	p.Initialize()

	rootCmd := &cobra.Command{
		Use:   "tm",
		Short: "Triggermesh CLI",
		Long: `⛵ Manage TriggerMesh services to ease serverless events.

* Manage your serverless functions.
* Build container images.
* Integrate functions at Git repos.`,
		SilenceUsage: true,
	}

	rootCmd.AddCommand(version.NewVersionCommand())

	return rootCmd
}
