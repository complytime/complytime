// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/action"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/models/components"
	"github.com/oscal-compass/oscal-sdk-go/validation"
	"github.com/spf13/cobra"

	"github.com/complytime/complytime/cmd/complytime/option"
	"github.com/complytime/complytime/internal/complytime"
)

// generateOptions defines options for the "generate" subcommand
type generateOptions struct {
	*option.Common
	complyTimeOpts *option.ComplyTime
}

// generateCmd creates a new cobra.Command for the "generate" subcommand
func generateCmd(common *option.Common) *cobra.Command {
	generateOpts := &generateOptions{
		Common:         common,
		complyTimeOpts: &option.ComplyTime{},
	}
	cmd := &cobra.Command{
		Use:     "generate [flags]",
		Short:   "Generate PVP policy from an assessment plan",
		Example: "complytime generate",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runGenerate(cmd, generateOpts)
		},
	}
	generateOpts.complyTimeOpts.BindFlags(cmd.Flags())
	return cmd
}

func runGenerate(cmd *cobra.Command, opts *generateOptions) error {
	validator := validation.NewSchemaValidator()
	ap, _, err := loadPlan(opts.complyTimeOpts, validator)
	if err != nil {
		return err
	}

	planSettings, err := getPlanSettings(opts.complyTimeOpts, ap)
	if err != nil {
		return err
	}

	// Set the framework ID from state (assessment plan)
	frameworkProp, valid := extensions.GetTrestleProp(extensions.FrameworkProp, *ap.Metadata.Props)
	if !valid {
		return fmt.Errorf("error reading framework property from assessment plan")
	}
	opts.complyTimeOpts.FrameworkID = frameworkProp.Value

	// Create the application directory if it does not exist
	appDir, err := complytime.NewApplicationDirectory(true)
	if err != nil {
		return err
	}
	logger.Debug(fmt.Sprintf("Using application directory: %s", appDir.AppDir()))
	cfg, err := complytime.Config(appDir)
	if err != nil {
		return err
	}
	logger.Debug("The configuration from the C2PConfig was successfully loaded.")

	// set config logger to CLI charm logger
	cfg.Logger = logger

	// Complete pre-processing of the components in Component Definition for use with actions
	// TODO: Replace with AP logic when full Assessment Plan support is in place upstream.
	compDefBundles, err := complytime.FindComponentDefinitions(appDir.BundleDir(), validator)
	if err != nil {
		return err
	}
	var allComponents []components.Component
	for _, compDef := range compDefBundles {
		if compDef.Components == nil {
			continue
		}
		for _, comp := range *compDef.Components {
			allComponents = append(allComponents, components.NewDefinedComponentAdapter(comp))
		}
	}

	inputContext, err := action.NewContext(allComponents)
	if err != nil {
		return fmt.Errorf("error generating context from directory %s: %w", appDir.AppDir(), err)
	}

	manager, err := framework.NewPluginManager(cfg)
	if err != nil {
		return fmt.Errorf("error initializing plugin manager: %w", err)
	}

	plugins, err := getGeneratorPlugins(inputContext, manager, opts.complyTimeOpts)
	defer manager.Clean()
	if err != nil {
		return err
	}

	inputContext.Settings = planSettings
	err = action.GeneratePolicy(cmd.Context(), plugins, inputContext)
	if err != nil {
		return err
	}
	logger.Info("Policy generation process completed for available plugins.")
	return nil
}

func getGeneratorPlugins(inputCtx *action.InputContext, manager *framework.PluginManager, opts *option.ComplyTime) (map[plugin.ID]policy.Generator, error) {
	manifests, err := manager.FindRequestedPlugins(inputCtx.RequestedProviders(), plugin.GenerationPluginName)
	if err != nil {
		return nil, err
	}
	pluginOptions := opts.ToPluginOptions()
	getSelections := func(_ plugin.ID) map[string]string {
		return pluginOptions.ToMap()
	}
	if err := pluginOptions.Validate(); err != nil {
		return nil, fmt.Errorf("failed plugin config validation: %w", err)
	}
	plugins, err := manager.LaunchGeneratorPlugins(manifests, getSelections)
	if err != nil {
		return nil, fmt.Errorf("errors launching plugins: %w", err)
	}
	logger.Info(fmt.Sprintf("Successfully loaded %v plugin(s).", len(plugins)))
	return plugins, nil
}
