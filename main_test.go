package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func argsForTestNoOptional(cloudProvider string, configFilepath string) []string {
	return []string{
		"./hedwig-terraform-generator",
		fmt.Sprintf("--%s", cloudProviderFlag),
		cloudProvider,
		"generate",
		configFilepath,
	}
}

func argsForTest(cloudProvider string, configFilepath string) []string {
	args := []string{
		"./hedwig-terraform-generator",
		fmt.Sprintf("--%s", cloudProviderFlag),
		cloudProvider,
		"generate",
		configFilepath,
	}
	if cloudProvider == cloudProviderAWS {
		args = append(
			args,
			"--alerting",
			fmt.Sprintf(`--%s=12345`, awsAccountIDFlag),
			fmt.Sprintf(`--%s=us-west-2`, awsRegionFlag),
			fmt.Sprintf(`--%s=pager_action`, queueAlertAlarmActionsFlag),
			fmt.Sprintf(`--%s=pager_action2`, queueAlertAlarmActionsFlag),
			fmt.Sprintf(`--%s=pager_action`, queueAlertOKActionsFlag),
			fmt.Sprintf(`--%s=pager_action2`, queueAlertOKActionsFlag),
			fmt.Sprintf(`--%s=pager_action`, dlqAlertAlarmActionsFlag),
			fmt.Sprintf(`--%s=pager_action2`, dlqAlertAlarmActionsFlag),
			fmt.Sprintf(`--%s=pager_action`, dlqAlertOKActionsFlag),
			fmt.Sprintf(`--%s=pager_action2`, dlqAlertOKActionsFlag),
		)
	} else if cloudProvider == cloudProviderGoogle {
		args = append(
			args,
			fmt.Sprintf(`--%s=gs://myBucket/tmp`, googleDataFlowTmpGCSLocationFlag),
			fmt.Sprintf(
				`--%s=gs://dataflow-templates/2019-04-03-00/Cloud_PubSub_to_Cloud_PubSub`,
				googleDataFlowTemplateGCSPathFlag,
			),
		)
	}
	return args
}

func TestGenerate(t *testing.T) {
	info, err := ioutil.ReadDir("test_fixtures")
	require.NoError(t, err)

	dmp := diffmatchpatch.New()

	for _, cloudDir := range info {
		if !cloudDir.IsDir() {
			continue
		}

		cloudProvider := cloudDir.Name()

		info, err := ioutil.ReadDir(filepath.Join("test_fixtures", cloudDir.Name()))
		require.NoError(t, err)

		for _, testDir := range info {
			testDirFullPath := filepath.Join("test_fixtures", cloudProvider, testDir.Name())

			infoTestDir, err := ioutil.ReadDir(testDirFullPath)
			require.NoError(t, err)

			err = os.RemoveAll("hedwig")
			require.NoError(t, err)

			fmt.Printf("Test: %s/%s\n", cloudProvider, testDir.Name())

			configFilepath := filepath.Join(testDirFullPath, "test_config.json")

			var args []string
			if strings.Contains(testDir.Name(), "no_optional_param") {
				args = argsForTestNoOptional(cloudProvider, configFilepath)
			} else {
				args = argsForTest(cloudProvider, configFilepath)
			}

			assert.NoError(t, runApp(args))

			info, err := ioutil.ReadDir("hedwig")
			assert.NoError(t, err)

			files := make([]string, len(info))
			for i, f := range info {
				files[i] = f.Name()
			}

			var testFiles []string
			for _, testOutputFile := range infoTestDir {
				if filepath.Ext(testOutputFile.Name()) != ".tf" {
					continue
				}
				testFiles = append(testFiles, testOutputFile.Name())
			}

			assert.Equal(t, testFiles, files)

			for _, testOutputFile := range infoTestDir {
				if filepath.Ext(testOutputFile.Name()) != ".tf" {
					continue
				}
				testOutputFileName := filepath.Join(testDirFullPath, testOutputFile.Name())
				expectedBytes, err := ioutil.ReadFile(testOutputFileName)
				require.NoError(t, err)

				expected := string(expectedBytes)

				// poor template engine
				expected = strings.Replace(
					expected, "{{TFAWSQueueModuleVersion}}", TFAWSQueueModuleVersion, -1)
				expected = strings.Replace(
					expected, "{{TFAWSQueueSubscriptionModuleVersion}}", TFAWSQueueSubscriptionModuleVersion, -1)
				expected = strings.Replace(
					expected, "{{TFAWSLambdaSubscriptionModuleVersion}}", TFAWSLambdaSubscriptionModuleVersion, -1)
				expected = strings.Replace(
					expected, "{{TFAWSTopicModuleVersion}}", TFAWSTopicModuleVersion, -1)
				expected = strings.Replace(
					expected, "{{TFGoogleTopicModuleVersion}}", TFGoogleTopicModuleVersion, -1)
				expected = strings.Replace(
					expected, "{{TFGoogleQueueModuleVersion}}", TFGoogleQueueModuleVersion, -1)
				expected = strings.Replace(
					expected, "{{TFGoogleSubscriptionModuleVersion}}", TFGoogleSubscriptionModuleVersion, -1)
				expected = strings.Replace(
					expected, "{{GENERATOR_VERSION}}", VERSION, -1)

				actualB, err := ioutil.ReadFile(filepath.Join("hedwig", testOutputFile.Name()))
				require.NoError(t, err)

				assert.Equal(
					t, expected, string(actualB),
					dmp.DiffPrettyText(dmp.DiffMain(expected, string(actualB), true)),
				)
			}

			if t.Failed() {
				// so we can investigate what went wrong
				break
			}
		}

		if t.Failed() {
			// so we can investigate what went wrong
			break
		}
	}
}
