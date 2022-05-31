package types

import (
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"testing"
)

type TestPlatform struct {
	T          *testing.T
	TestFolder string
	//username         string
	//TerraformOptions *terraform.Options
	//KeyPair          *aws.Ec2Keypair
	//publicIP         string
	//publicHost       ssh.Host
}

func NewTestPlatform(t *testing.T, testFolder string) *TestPlatform {
	tp := new(TestPlatform)
	tp.T = t
	tp.TestFolder = testFolder
	return tp
}

func (platform *TestPlatform) Teardown() {
	keyPair := teststructure.LoadEc2KeyPair(platform.T, platform.TestFolder)
	terraformOptions := teststructure.LoadTerraformOptions(platform.T, platform.TestFolder)
	terraform.Destroy(platform.T, terraformOptions)
	aws.DeleteEC2KeyPair(platform.T, keyPair)
}

func (platform *TestPlatform) RunSSHCommand(command string) (string, error) {
	terraformOptions := teststructure.LoadTerraformOptions(platform.T, platform.TestFolder)
	keyPair := teststructure.LoadEc2KeyPair(platform.T, platform.TestFolder)
	host := ssh.Host{
		Hostname:    terraform.Output(platform.T, terraformOptions, "public_instance_ip"),
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: "ubuntu",
	}
	return ssh.CheckSshCommandE(platform.T, host, command)
}

//func (e2e *E2ETest) runSSHCommand(format string, a ...interface{}) (string, error) {
//	command := fmt.Sprintf(format, a...)
//	return ssh.CheckSshCommandE(e2e.testing, e2e.publicHost, command)
//}
