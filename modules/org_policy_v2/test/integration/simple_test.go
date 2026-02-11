// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/stretchr/testify/assert"
)

func TestSimpleExample(t *testing.T) {
	asserts := assert.New(t)

	// 1. Initialize and APPLY the Setup stage (creates the project)
	setup := tft.NewTFBlueprintTest(t, tft.WithTFDir("../setup"))
	setup.Apply(asserts)

	// SAFETY: Ensure the dynamic project is destroyed even if the test fails
	defer setup.Teardown(asserts)

	// 2. Initialize the Example stage, passing the new project_id
	orgPolicyTest := tft.NewTFBlueprintTest(t,
		tft.WithTFDir("../../examples/simple_example"),
		tft.WithVars(map[string]interface{}{
			"project_id": setup.GetJsonOutput("project_id").String(),
		}),
	)

	// 3. Define the Verification logic
	orgPolicyTest.DefineVerify(func(verifyAsserts *assert.Assertions) {
		orgPolicyTest.DefaultVerify(verifyAsserts)

		projectID := orgPolicyTest.GetStringOutput("project_id")
		
		// Run gcloud command to verify enforcement
		op := gcloud.Run(t, fmt.Sprintf("resource-manager org-policies describe compute.vmExternalIpAccess --project=%s --effective --format=json", projectID))
		
		// Updated Path: Look for 'listPolicy.allValues'
		// Updated Expectation: The value is "DENY" for a Deny-All policy
		allValues := op.Get("listPolicy.allValues").String()
		verifyAsserts.Equal("DENY", allValues, "The Policy should be enforcing a DENY allValues rule")
	})
	// 4. Run the test (Apply Example -> Verify -> Destroy Example)
	orgPolicyTest.Test()
}