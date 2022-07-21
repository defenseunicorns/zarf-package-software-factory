package types

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"
)

// TestPlatform is the test "state" that allows for helper functions such as deferring the teardown step.
type TestPlatform struct {
	T          *testing.T
	TestFolder string
}

// NewTestPlatform generates the test "state" object that allows for helper functions such as deferring the teardown step.
func NewTestPlatform(t *testing.T) *TestPlatform {
	t.Helper()
	testPlatform := new(TestPlatform)
	testPlatform.T = t
	tempFolder := teststructure.CopyTerraformFolderToTemp(t, "..", "tf/public-ec2-instance")
	testPlatform.TestFolder = tempFolder

	// Since Terraform is going to be run with that temp folder as the CWD, we also need our .tool-versions file to be
	// in that directory so that the right version of Terraform is being run there. I can neither confirm nor deny that
	// this took me 2 days to figure out...
	// Since we can't be sure what the working directory is, we are going to walk up one directory at a time until we
	// find a .tool-versions file and then copy it into the temp folder
	found := false
	filePath := ".tool-versions"
	for !found {
		//nolint:gocritic
		if _, err := os.Stat(filePath); err == nil {
			// The file exists
			found = true
		} else if errors.Is(err, os.ErrNotExist) {
			// The file does *not* exist. Add a "../" and try again
			filePath = fmt.Sprintf("../%v", filePath)
		} else {
			// Schrodinger: file may or may not exist. See err for details.
			// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
			require.NoError(t, err)
		}
	}
	err := copyFile(filePath, fmt.Sprintf("%v/.tool-versions", testPlatform.TestFolder))
	require.NoError(t, err)

	return testPlatform
}

// RunSSHCommandAsSudo provides a simple way to run a shell command on the server that is created using Terraform.
func (platform *TestPlatform) RunSSHCommandAsSudo(command string) (string, error) {
	terraformOptions := teststructure.LoadTerraformOptions(platform.T, platform.TestFolder)
	keyPair := teststructure.LoadEc2KeyPair(platform.T, platform.TestFolder)
	host := ssh.Host{
		Hostname:    terraform.Output(platform.T, terraformOptions, "public_instance_ip"),
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: "ec2-user",
	}
	output, err := ssh.CheckSshCommandE(platform.T, host, fmt.Sprintf(`sudo bash -c "%v"`, command))
	if err != nil {
		logger.Default.Logf(platform.T, output)

		return "nil", fmt.Errorf("ssh command failed: %w", err)
	}

	logger.Default.Logf(platform.T, output)

	return output, nil
}

// Teardown brings down the Terraform infrastructure that was created.
func (platform *TestPlatform) Teardown() {
	teststructure.RunTestStage(platform.T, "TEARDOWN", func() {
		keyPair := teststructure.LoadEc2KeyPair(platform.T, platform.TestFolder)
		terraformOptions := teststructure.LoadTerraformOptions(platform.T, platform.TestFolder)
		terraform.Destroy(platform.T, terraformOptions)
		aws.DeleteEC2KeyPair(platform.T, keyPair)
	})
}

// copyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherwise, attempt to create a hard link
// between the two files. If that fails, copy the file contents from src to dst.
func copyFile(src string, dest string) error {
	sourceFileInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}
	if !sourceFileInfo.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("non-regular source file %s (%q)", sourceFileInfo.Name(), sourceFileInfo.Mode().String())
	}
	destFileInfo, err := os.Stat(dest)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("unknown error: %w", err)
		}
	} else {
		if !(destFileInfo.Mode().IsRegular()) {
			return fmt.Errorf("non-regular destination file %s (%q)", destFileInfo.Name(), destFileInfo.Mode().String())
		}
		if os.SameFile(sourceFileInfo, destFileInfo) {
			return nil
		}
	}
	err = os.Link(src, dest)
	if err == nil {
		return nil
	}
	err = copyFileContents(src, dest)
	if err != nil {
		return err
	}

	return nil
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src string, dest string) error {
	cleanSrc := filepath.Clean(src)
	cleanDst := filepath.Clean(dest)
	srcFile, err := os.Open(cleanSrc)
	if err != nil {
		return fmt.Errorf("unable to open src file: %w", err)
	}
	defer func(in *os.File) {
		_ = in.Close()
	}(srcFile)
	dstFile, err := os.Create(cleanDst)
	if err != nil {
		return fmt.Errorf("unable to create dest file: %w", err)
	}
	defer func() {
		cerr := dstFile.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("unable to copy file: %w", err)
	}
	err = dstFile.Sync()

	return nil
}
