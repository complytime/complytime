// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"testing"
)

func TestSetOptsFromArgs(t *testing.T) {
	opts := &generateOptions{}
	args := []string{"assessment-plan.json"}
	setOptsFromArgs(args, opts)
	fmt.Println(opts.assessmentPlanPath)
}
