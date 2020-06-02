package service

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/triggermesh/tm/pkg/client"
	"github.com/triggermesh/tm/pkg/tm/commands"
)

func NewServiceCommand(p *commands.TmParams) *cobra.Command {
	serviceCmd := &cobra.Command{
		Use:   "service",
		Short: "Service command group",
	}
	serviceCmd.AddCommand(cmdListService(p))

	// serviceCmd.AddCommand(NewServiceDescribeCommand(p))
	// serviceCmd.AddCommand(NewServiceCreateCommand(p))
	// serviceCmd.AddCommand(NewServiceDeleteCommand(p))
	// serviceCmd.AddCommand(NewServiceUpdateCommand(p))
	// serviceCmd.AddCommand(NewServiceExportCommand(p))
	return serviceCmd
}

func cmdListService(p *commands.TmParams) *cobra.Command {

	return &cobra.Command{
		Use:     "service",
		Aliases: []string{"services"},
		Short:   "List of knative service resources",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			s.Namespace = client.Namespace
			if len(args) == 0 {
				list, err := s.List(clientset)
				if err != nil {
					clientset.Log.Fatalln(err)
				}
				if len(list.Items) == 0 {
					fmt.Fprintf(cmd.OutOrStdout(), "No services found\n")
					return
				}
				clientset.Printer.PrintTable(s.GetTable(list))
				return
			}
			s.Name = args[0]
			service, err := s.Get(clientset)
			if err != nil {
				clientset.Log.Fatalln(err)
			}
			clientset.Printer.PrintObject(s.GetObject(service))
		},
	}
}
