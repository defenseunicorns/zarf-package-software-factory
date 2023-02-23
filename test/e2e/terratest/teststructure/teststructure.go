// Package teststructure is customized test-structure code based on https://github.com/gruntwork-io/terratest/tree/5913a2925623d3998841cb25de7b26731af9ab13 that fixes the issue identified in https://github.com/gruntwork-io/terratest/issues/1135
package teststructure

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/logger"
	terratesting "github.com/gruntwork-io/terratest/modules/testing"
)

// SaveEc2KeyPair serializes and saves an Ec2KeyPair into the given folder. This allows you to create an Ec2KeyPair during setup
// and to reuse that Ec2KeyPair later during validation and teardown.
// This function is directly copied from https://github.com/gruntwork-io/terratest/tree/5913a2925623d3998841cb25de7b26731af9ab13
// due to this issue: https://github.com/gruntwork-io/terratest/issues/1135
func SaveEc2KeyPair(t terratesting.TestingT, testFolder string, keyPair *aws.Ec2Keypair) {
	saveTestData(t, formatEc2KeyPairPath(testFolder), keyPair)
}

// SaveTestData serializes and saves a value used at test time to the given path. This allows you to create some sort of test data
// (e.g., TerraformOptions) during setup and to reuse this data later during validation and teardown.
// This function is directly copied from https://github.com/gruntwork-io/terratest/tree/5913a2925623d3998841cb25de7b26731af9ab13
// due to this issue: https://github.com/gruntwork-io/terratest/issues/1135
func saveTestData(t terratesting.TestingT, path string, value interface{}) {
	logger.Default.Logf(t, "Storing test data in %s so it can be reused later", path)

	if IsTestDataPresent(t, path) {
		logger.Default.Logf(t, "[WARNING] The named test data at path %s is non-empty. Save operation will overwrite existing value with \"%v\".\n.", path, value)
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("Failed to convert value %s to JSON: %v", path, err)
	}

	// Don't log this data, it exposes the EC2 Key Pair's private key in the logs, which are public on GitHub Actions
	// logger.Logf(t, "Marshalled JSON: %s", string(bytes))

	parentDir := filepath.Dir(path)
	if err := os.MkdirAll(parentDir, 0750); err != nil { //nolint:gomnd
		t.Fatalf("Failed to create folder %s: %v", parentDir, err)
	}

	if err := os.WriteFile(path, bytes, 0600); err != nil { //nolint:gomnd
		t.Fatalf("Failed to save value %s: %v", path, err)
	}
}

// formatEc2KeyPairPath formats a path to save an Ec2KeyPair in the given folder.
// This function is directly copied from https://github.com/gruntwork-io/terratest/tree/5913a2925623d3998841cb25de7b26731af9ab13
// due to this issue: https://github.com/gruntwork-io/terratest/issues/1135
func formatEc2KeyPairPath(testFolder string) string {
	return formatTestDataPath(testFolder, "Ec2KeyPair.json")
}

// FormatTestDataPath formats a path to save test data.
// This function is directly copied from https://github.com/gruntwork-io/terratest/tree/5913a2925623d3998841cb25de7b26731af9ab13
// due to this issue: https://github.com/gruntwork-io/terratest/issues/1135
func formatTestDataPath(testFolder string, filename string) string {
	return filepath.Join(testFolder, ".test-data", filename)
}

// IsTestDataPresent returns true if a file exists at $path and the test data there is non-empty.
// This function is directly copied from https://github.com/gruntwork-io/terratest/tree/5913a2925623d3998841cb25de7b26731af9ab13
// due to this issue: https://github.com/gruntwork-io/terratest/issues/1135
func IsTestDataPresent(t terratesting.TestingT, path string) bool {
	exists, err := files.FileExistsE(path)
	if err != nil {
		t.Fatalf("Failed to load test data from %s due to unexpected error: %v", path, err)
	}
	if !exists {
		return false
	}

	bytes, err := os.ReadFile(path)

	if err != nil {
		t.Fatalf("Failed to load test data from %s due to unexpected error: %v", path, err)
	}

	if isEmptyJSON(t, bytes) {
		return false
	}

	return true
}

// isEmptyJSON returns true if the given bytes are empty, or in a valid JSON format that can reasonably be considered empty.
// The types used are based on the type possibilities listed at https://golang.org/src/encoding/json/decode.go?s=4062:4110#L51
// This function is directly copied from https://github.com/gruntwork-io/terratest/tree/5913a2925623d3998841cb25de7b26731af9ab13
// due to this issue: https://github.com/gruntwork-io/terratest/issues/1135
//
//nolint:cyclop
func isEmptyJSON(t terratesting.TestingT, bytes []byte) bool {
	var value interface{}

	if len(bytes) == 0 {
		return true
	}

	if err := json.Unmarshal(bytes, &value); err != nil {
		t.Fatalf("Failed to parse JSON while testing whether it is empty: %v", err)
	}

	if value == nil {
		return true
	}

	valueBool, ok := value.(bool)
	if ok && !valueBool {
		return true
	}

	valueFloat64, ok := value.(float64)
	if ok && valueFloat64 == 0 {
		return true
	}

	valueString, ok := value.(string)
	if ok && valueString == "" {
		return true
	}

	valueSlice, ok := value.([]interface{})
	if ok && len(valueSlice) == 0 {
		return true
	}

	valueMap, ok := value.(map[string]interface{})
	if ok && len(valueMap) == 0 {
		return true
	}

	return false
}
