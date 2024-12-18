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

// scanOptions defined options for the scan subcommand.
type scanOptions struct {
	*option.Common
	assessmentPlanPath string
}

// scanCmd creates a new cobra.Command for the version subcommand.
func scanCmd(common *option.Common) *cobra.Command {
	scanOpts := &scanOptions{Common: common}
	return &cobra.Command{
		Use:          "scan [flags] <assessment_plan_path>",
		Short:        "Scan environment with assessment plan",
		Example:      "complytime scan assessment-plan.json",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			scanOpts.assessmentPlanPath = filepath.Clean(args[0])
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runScan(cmd, scanOpts)
		},
	}
}

func runScan(cmd *cobra.Command, opts *scanOptions) error {
	// Adding this message to the user for now because assessment Plans are unused
	_, _ = fmt.Fprintf(opts.Out, "OSCAL Assessment Plans are not supported yet. The file %s will not be used\n", opts.assessmentPlanPath)

	// Create the application directory if it does not exist
	appDir, err := complytime.NewApplicationDirectory(true)
	if err != nil {
		return err
	}

	cfg, err := complytime.Config(appDir)
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
	// Ensure all the plugin launch above are cleaned up
	defer manager.Clean()

	allResults, err := manager.AggregateResults(cmd.Context(), plugins)
	if err != nil {
		return err
	}

	for _, result := range allResults {
		_, _ = fmt.Fprintf(opts.Out, "Result:\n%s", result)
	}
	return nil
}
