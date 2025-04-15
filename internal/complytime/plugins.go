// SPDX-License-Identifier: Apache-2.0

package complytime

import (
	"errors"
	"reflect"
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
