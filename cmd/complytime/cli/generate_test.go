// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const assessmentPlan = "assessment-plan.json"

func TestSetOptsFromArgs(t *testing.T) {
	opts := &generateOptions{}
	args := []string{assessmentPlan}
	setOptsFromArgs(args, opts)

	require.Equal(t, opts.assessmentPlanPath, assessmentPlan)
}
