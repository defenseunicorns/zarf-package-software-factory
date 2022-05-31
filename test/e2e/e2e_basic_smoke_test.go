package test

import (
	"testing"

	"github.com/defenseunicorns/zarf-package-software-factory/test/e2e/utils"
	"github.com/stretchr/testify/require"
)

func TestBasicSmoke(t *testing.T) {
	// BOILERPLATE, EXPECTED TO BE PRESENT AT THE BEGINNING OF EVERY TEST
	t.Parallel()
	platform := utils.InitTestPlatform(t)
	defer platform.Teardown()
	utils.SetupTestPlatform(t, platform)

	// TEST CODE STARTS HERE
	output, err := platform.RunSSHCommand("ls -la ~/app")
	require.NoError(t, err, output)
	t.Logf("%v", output)
}
