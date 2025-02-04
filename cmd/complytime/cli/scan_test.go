// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/complytime/complytime/cmd/complytime/option"
)

func TestLoadAssessmentPlan(t *testing.T) {

	testOpts := &option.ComplyTime{
		UserWorkspace: "./testdata",
	}

	assessmentPlan, err := loadAssessmentPlan(testOpts)
	require.NoError(t, err)
	require.Equal(t, assessmentPlan.Metadata.Title, "REPLACE_ME")
}
