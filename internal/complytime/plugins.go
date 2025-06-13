// SPDX-License-Identifier: Apache-2.0

package complytime

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/actions"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

// PluginOptions defines global options all complytime plugins should
// support.
type PluginOptions struct {
	// Workspace is the location where all
	// plugin outputs should be written.
	Workspace string `config:"workspace"`
	// Profile is the compliance profile that the plugin should use for
	// pre-defined policy groups.
	Profile string `config:"profile"`
	// UserConfigRoot is the root directory where users customize
	// plugin configuration options
	UserConfigRoot string `config:"userconfigroot"`
}

// NewPluginOptions created a new PluginOptions struct.
func NewPluginOptions() PluginOptions {
	return PluginOptions{}
}

// Validate ensure the required plugin options are set.
func (p PluginOptions) Validate() error {
	// TODO[jpower432]: If these options grow, using third party
	// validation through struct tags could be simpler if the validation
	// logic gets more complex.
	if p.Workspace == "" {
		return errors.New("workspace must be set")
	}
	if p.Profile == "" {
		return errors.New("profile must be set")
	}
	if p.UserConfigRoot != "" {
		if _, err := os.Stat(p.UserConfigRoot); os.IsNotExist(err) {
			return errors.New("user config root does not exist")
		}
	}
	return nil
}

// ToMap transforms the PluginOption struct into a map that can be consumed
// by the C2P Plugin Manager.
func (p PluginOptions) ToMap() map[string]string {
	selections := make(map[string]string)
	val := reflect.ValueOf(p)
	t := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := t.Field(i)
		key := fieldType.Tag.Get("config")
		selections[key] = field.String()
	}
	return selections
}

func (p PluginOptions) updateOptions(manifest plugin.Manifest) error {
	configPath := filepath.Join(p.UserConfigRoot, manifest.Metadata.ID.String()+".json")
	configFile, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to open plugin config file: %w", err)
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	var options map[string]string
	err = jsonParser.Decode(&options)
	if err != nil {
		return fmt.Errorf("failed to parse plugin config file: %w", err)
	}
	// Update manifest configuration
	for index, configurationOption := range manifest.Configuration {
		optionConfig, ok := options[configurationOption.Name]
		if ok {
			configurationOption.Default = &optionConfig
			manifest.Configuration[index] = configurationOption
		}
	}
	return nil
}

// Plugins launches and configures plugins with the given complytime global options. This function returns the plugin map with the
// launched plugins, a plugin cleanup function, and an error. The cleanup function should be used if it is not nil.
func Plugins(manager *framework.PluginManager, inputs *actions.InputContext, selections PluginOptions) (map[plugin.ID]policy.Provider, func(), error) {
	manifests, err := manager.FindRequestedPlugins(inputs.RequestedProviders())
	if err != nil {
		return nil, nil, err
	}

	if err := selections.Validate(); err != nil {
		return nil, nil, fmt.Errorf("failed plugin config validation: %w", err)
	}

	selectionsMap := selections.ToMap()
	getSelections := func(_ plugin.ID) map[string]string {
		return selectionsMap
	}

	// Update the manifests to overide options default
	if selections.UserConfigRoot != "" {
		for _, manifest := range manifests {
			if len(manifest.Configuration) > 0 {
				if err := selections.updateOptions(manifest); err != nil {
					return nil, nil, fmt.Errorf("failed to update plugin configuration: %w", err)
				}
			}
		}
	}
	plugins, err := manager.LaunchPolicyPlugins(manifests, getSelections)
	// Plugin subprocess has now been launched; cleanup always required below
	if err != nil {
		return nil, manager.Clean, err
	}
	return plugins, manager.Clean, nil
}
