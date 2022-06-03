package test_test

import (
	"context"
	"testing"
	"time"

	"github.com/defenseunicorns/zarf-package-software-factory/test/e2e/types"
	"github.com/defenseunicorns/zarf-package-software-factory/test/e2e/utils"
	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"
)

// TestAllServicesRunning waits until all services report that they are ready.
func TestAllServicesRunning(t *testing.T) {
	// BOILERPLATE, EXPECTED TO BE PRESENT AT THE BEGINNING OF EVERY TEST FUNCTION

	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	go utils.HoldYourDamnHorses(ctx, t, 10*time.Second)
	defer cancel()
	platform := types.NewTestPlatform(t)
	defer platform.Teardown()
	utils.SetupTestPlatform(t, platform)
	// The repo has now been downloaded to /root/app and the software factory package deployment has been initiated.
	teststructure.RunTestStage(platform.T, "TEST", func() {
		// END BOILERPLATE

		// TEST CODE STARTS HERE.

		// Just make sure we can hit the cluster
		output, err := platform.RunSSHCommandAsSudo("kubectl get nodes")
		require.NoError(t, err, output)
		// Wait up to 18 minutes for the "bigbang" kustomization to report "Ready==True". Our testing shows if everything goes right this should take 11-13 minutes.
		output, err = platform.RunSSHCommandAsSudo("kubectl wait --timeout=1080s -n flux-system --for=condition=Ready kustomization/bigbang")
		require.NoError(t, err, output)
		// Wait up to 4 additional minutes for the "softwarefactoryaddons-deps" kustomization to report "Ready==True".
		output, err = platform.RunSSHCommandAsSudo("kubectl wait --timeout=240s -n flux-system --for=condition=Ready kustomization/softwarefactoryaddons-deps")
		require.NoError(t, err, output)
		// Wait up to 4 additional minutes for the "softwarefactoryaddons" kustomization to report "Ready==True".
		output, err = platform.RunSSHCommandAsSudo("kubectl wait --timeout=240s -n flux-system --for=condition=Ready kustomization/softwarefactoryaddons")
		require.NoError(t, err, output)
	})
}
