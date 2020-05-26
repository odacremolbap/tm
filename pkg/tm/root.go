package tm

import "github.com/spf13/cobra"

// NewRootCommand creates the root command for tm
func NewRootCommand() {
	// root := newRootCommand() {

	// }
}

func newRootCommand() {

	rootCmd := &cobra.Command{
		Use:   "tm",
		Short: "Triggermesh CLI",
		Long: `Manage TriggerMesh ⛵ services to ease serverless events.

* Manage your serverless functions.
* Build container images.
* Integrate functions at Git repos.`,
		SilenceUsage: true,
		Version:      version,
	}

}
