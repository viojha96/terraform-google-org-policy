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