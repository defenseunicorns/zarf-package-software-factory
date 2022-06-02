package test

import (
	"context"
	"testing"

	"github.com/defenseunicorns/zarf-package-software-factory/test/e2e/utils"
	"github.com/stretchr/testify/require"
)

func TestBasicSmoke(t *testing.T) {
	// BOILERPLATE, EXPECTED TO BE PRESENT AT THE BEGINNING OF EVERY TEST FUNCTION

	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	go utils.HoldYourDamnHorses(ctx, t)
	defer cancel()
	platform := utils.InitTestPlatform(t)
	defer platform.Teardown()
	utils.SetupTestPlatform(t, platform)
	// The repo has now been downloaded to /root/app and the software factory package deployment has been initiated.

	// TEST CODE STARTS HERE.

	// Just make sure we can hit the cluster
	output, err := platform.RunSSHCommandAsSudo("kubectl get nodes")
	require.NoError(t, err, output)
	// Wait up to 16 minutes for the "bigbang" kustomization to report "Ready==True". Our testing shows if everything goes right this should take 11-13 minutes.
	output, err = platform.RunSSHCommandAsSudoWithTimeout(`while [[ $(kubectl get kustomization bigbang -n flux-system -o json | jq -r '.status.conditions[] | select(.type == "Ready") | .status') != "True" ]]; do sleep 3; done`, 960)
	require.NoError(t, err, output)
	// Wait up to 2 additional minutes for the "softwarefactoryaddons-deps" kustomization to report "Ready==True".
	output, err = platform.RunSSHCommandAsSudoWithTimeout(`while [[ $(kubectl get kustomization softwarefactoryaddons-deps -n flux-system -o json | jq -r '.status.conditions[] | select(.type == "Ready") | .status') != "True" ]]; do sleep 3; done`, 120)
	require.NoError(t, err, output)
	// Wait up to 2 additional minutes for the "softwarefactoryaddons" kustomization to report "Ready==True".
	output, err = platform.RunSSHCommandAsSudoWithTimeout(`while [[ $(kubectl get kustomization softwarefactoryaddons -n flux-system -o json | jq -r '.status.conditions[] | select(.type == "Ready") | .status') != "True" ]]; do sleep 3; done`, 120)
	require.NoError(t, err, output)
}
