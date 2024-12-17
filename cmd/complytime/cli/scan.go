// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"path/filepath"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/spf13/cobra"

	"github.com/complytime/complytime/cmd/complytime/option"
	"github.com/complytime/complytime/internal/complytime"
)

// scanCmd creates a new cobra.Command for the version subcommand.
func scanCmd(common *option.Common) *cobra.Command {
	return &cobra.Command{
		Use:   "scan",
		Short: "",
		Long:  "",
		Args:  cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			// Validate and process args
			var assessmentPlanPath string
			if len(args) == 1 {
				assessmentPlanPath = filepath.Clean(args[0])
			}
			_, _ = fmt.Fprintf(common.Out, "Assessment Plans are not supported yet. The file %s will not be used\n", assessmentPlanPath)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runScan(cmd, common)
		},
	}
}

func runScan(cmd *cobra.Command, opts *option.Common) error {
	cfg, err := complytime.Config()
	if err != nil {
		return err
	}
	manager, err := framework.NewPluginManager(cfg)
	if err != nil {
		return fmt.Errorf("error initializing plugin manager: %w", err)
	}
	plugins, err := manager.LaunchPolicyPlugins()
	if err != nil {
		return err
	}
	allResults, err := manager.AggregateResults(cmd.Context(), plugins)
	if err != nil {
		return err
	}
	manager.Clean()
	for _, result := range allResults {
		_, _ = fmt.Fprintf(opts.Out, "Result:\n%s", result)
	}
	return nil
}
