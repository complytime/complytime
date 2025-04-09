// SPDX-License-Identifier: Apache-2.0

package option

import (
	"io"

	"github.com/spf13/pflag"
)

// Common options for the ComplyTime CLI.
type Common struct {
	Debug bool
	Output
}

// Output options for
type Output struct {
	// Out think, os.Stdout
	Out io.Writer
	// ErrOut think, os.Stderr
	ErrOut io.Writer
}

// BindFlags populate Common options from user-specified flags.
func (o *Common) BindFlags(fs *pflag.FlagSet) {
	fs.BoolVarP(&o.Debug, "debug", "d", false, "output debug logs")
}

// ComplyTime options are configurations needed for the ComplyTime CLI to run.
// They are less generic the Common options and would only be used in a subset of
// commands.
type ComplyTime struct {
	UserWorkspace string
}

type UpdatePlan struct {
	Controls   map[string]string
	Rules      map[string]string
	Parameters map[string]string
	Config     string
}

// BindFlags populate ComplyTime options from user-specified flags.
func (o *ComplyTime) BindFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.UserWorkspace, "workspace", "w", ".", "workspace to use for artifact generation")
}
func (o *UpdatePlan) BindFlags(fs *pflag.FlagSet) {
	fs.StringP("config", "c", "./config/complytime", "location of complytime bundles and controls in home dir")
	fs.StringSliceP("update-controls", "u", []string{}, "update controls from 'ready to assess' to waived")
	fs.StringSlice("exclude-controls", []string{}, "controls to be excluded from assessment-plan.json")
	fs.StringSlice("update-rules", []string{}, "update rules from 'ready to assess' to waived")
	fs.StringSlice("exclude-rules", []string{}, "rules to be excluded from assessment-plan.json")
	fs.StringSlice("update-parameters", []string{}, "update parameters from 'ready to assess' to waived")
}
