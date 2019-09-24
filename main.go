package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1" // imports as package "cli"
)

const (
	cloudProviderGoogle = "google"
	cloudProviderAWS    = "aws"
)

const (
	// VERSION represents the version of the generator tool
	VERSION = "v3.2.0"

	// TFAWSQueueModuleVersion represents the version of the AWS hedwig-queue module
	TFAWSQueueModuleVersion = "1.0.0"

	// TFAWSQueueSubscriptionModuleVersion represents the version of the AWS hedwig-queue-subscription module
	TFAWSQueueSubscriptionModuleVersion = "1.0.0"

	// TFAWSLambdaSubscriptionModuleVersion represents the version of the AWS hedwig-lambda-subscription module
	TFAWSLambdaSubscriptionModuleVersion = "1.0.0"

	// TFAWSTopicModuleVersion represents the version of the AWS hedwig-topic module
	TFAWSTopicModuleVersion = "1.0.0"

	// TFGoogleTopicModuleVersion represents the version of the Google hedwig-topic module
	TFGoogleTopicModuleVersion = "1.2.1"

	// TFGoogleQueueModuleVersion represents the version of the Google hedwig-queue module
	TFGoogleQueueModuleVersion = "2.0.0"

	// TFGoogleAlertsModuleVersion represents the version of the Google hedwig-alerts module
	TFGoogleAlertsModuleVersion = "1.1.2"

	// TFGoogleSubscriptionModuleVersion represents the version of the Google hedwig-subscription module
	TFGoogleSubscriptionModuleVersion = "2.0.1"

	tfDoNotEditStamp = `// DO NOT EDIT
// This file has been auto-generated by hedwig-terraform-generator ` + VERSION
)

const (
	// alertingFlag represents the cli flag that indicates if alerting should be generated
	alertingFlag = "alerting"

	// awsAccountIDFlag represents the cli flag for aws account id (AWS only)
	awsAccountIDFlag = "aws-account-id"

	// awsRegionFlag represents the cli flag for aws region (AWS only)
	awsRegionFlag = "aws-region"

	// cloudProviderFlag represents the cli flag for cloud provider name
	cloudProviderFlag = "cloud"

	// dlqAlertAlarmActionsFlag represents the cli flag for DLQ alert actions on ALARM (AWS only)
	dlqAlertAlarmActionsFlag = "dlq-alert-alarm-actions"

	// dlqAlertOKActionsFlag represents the cli flag for DLQ alert notification channels (Google only)
	dlqAlertNotificationChannelsFlag = "dlq-alert-notification-channels"

	// dlqAlertOKActionsFlag represents the cli flag for DLQ alert actions on OK (AWS only)
	dlqAlertOKActionsFlag = "dlq-alert-ok-actions"

	// dataflowTmpGCSLocationFlag represents the cli flag for Dataflow temporary GCS location (Google only)
	dataflowTmpGCSLocationFlag = "dataflow-tmp-gcs-location"

	// dataflowPubSubToPubSubTemplateGCSPathFlag represents the cli flag for Dataflow template GCS path
	// for pub sub to pub sub dataflow (Google only)
	dataflowPubSubToPubSubTemplateGCSPathFlag = "dataflow-template-pubsub-to-pubsub-gcs-path"

	// dataflowPubSubToStorageGCSPathFlag represents the cli flag for Dataflow template GCS path
	// for pub sub to Storage dataflow (Google only)
	dataflowPubSubToStorageGCSPathFlag = "dataflow-template-pubsub-to-storage-gcs-path"

	// googleDataflowZoneFlag represents the cli flag for Dataflow template GCS zone (Google only)
	googleDataflowZoneFlag = "dataflow-zone"

	// googleDataflowRegionFlag represents the cli flag for Dataflow template GCS region (Google only)
	googleDataflowRegionFlag = "dataflow-region"

	// googleFirehoseDataflowOutputDirectoryFlag represents the cli flag for Firehose Dataflow output directory
	// (Google only)
	googleFirehoseDataflowOutputDirectoryFlag = "google-firehose-dataflow-output-dir"

	// googleProjectAlerting represents the cli flag that indicates the google project for alerting resources
	googleProjectAlerting = "google-project-alerting"

	// enableFirehoseAllTopics represents the cli flag to enable Google Firehose for all hedwig topics (Google only
	// for now)
	enableFirehoseAllTopics = "enable-firehose-all-topics"

	// moduleFlag represents the cli flag for output module name
	moduleFlag = "module"

	// queueAlertAlarmActionsFlag represents the cli flag for queue alert actions on ALARM (AWS only)
	queueAlertAlarmActionsFlag = "queue-alert-alarm-actions"

	// queueAlertNotificationChannelsFlag represents the cli flag for queue alert notification channels (Google only)
	queueAlertNotificationChannelsFlag = "queue-alert-notification-channels"

	// queueAlertOKActionsFlag represents the cli flag for queue alert actions on OK
	queueAlertOKActionsFlag = "queue-alert-ok-actions"
)

var providerSpecificFlags = map[string][]string{
	cloudProviderAWS: {
		awsAccountIDFlag,
		awsRegionFlag,
		dlqAlertAlarmActionsFlag,
		dlqAlertOKActionsFlag,
		queueAlertAlarmActionsFlag,
		queueAlertOKActionsFlag,
	},
	cloudProviderGoogle: {
		dlqAlertNotificationChannelsFlag,
		dataflowPubSubToPubSubTemplateGCSPathFlag,
		dataflowPubSubToStorageGCSPathFlag,
		dataflowTmpGCSLocationFlag,
		enableFirehoseAllTopics,
		googleDataflowZoneFlag,
		googleDataflowRegionFlag,
		googleFirehoseDataflowOutputDirectoryFlag,
		googleProjectAlerting,
		queueAlertNotificationChannelsFlag,
	},
}

var providerRequiredFlags = map[string][]string{
	cloudProviderAWS: {
		awsAccountIDFlag,
		awsRegionFlag,
	},
	cloudProviderGoogle: {},
}

var providerAlertingFlags = map[string][]string{
	cloudProviderAWS: {
		queueAlertAlarmActionsFlag,
		queueAlertOKActionsFlag,
		dlqAlertAlarmActionsFlag,
		dlqAlertOKActionsFlag,
	},
	cloudProviderGoogle: {
		queueAlertNotificationChannelsFlag,
		dlqAlertNotificationChannelsFlag,
		googleProjectAlerting,
	},
}

var providerAlertingRequiredFlags = map[string][]string{
	cloudProviderAWS: {
		queueAlertAlarmActionsFlag,
		queueAlertOKActionsFlag,
		dlqAlertAlarmActionsFlag,
		dlqAlertOKActionsFlag,
	},
	cloudProviderGoogle: {
		queueAlertNotificationChannelsFlag,
		dlqAlertNotificationChannelsFlag,
	},
}

func validateArgs(c *cli.Context) *cli.ExitError {
	cloudProvider := c.GlobalString(cloudProviderFlag)
	if cloudProvider == "" {
		return cli.NewExitError(fmt.Sprintf("--%s is required", cloudProviderFlag), 1)
	}
	if cloudProvider != cloudProviderAWS && cloudProvider != cloudProviderGoogle {
		return cli.NewExitError(fmt.Sprintf("invalid cloud provider: %s", cloudProvider), 1)
	}

	if c.NArg() != 1 {
		return cli.NewExitError("<config-file> is required", 1)
	}

	// verify provider flags are used correctly
	for provider, flags := range providerSpecificFlags {
		if provider == cloudProvider {
			continue
		}
		for _, flag := range flags {
			if c.IsSet(flag) {
				return cli.NewExitError(
					fmt.Sprintf("flag --%s disallowed for provider: %s", flag, cloudProvider),
					1,
				)
			}
		}
	}

	// verify required flags are provided
	for _, flag := range providerRequiredFlags[cloudProvider] {
		if !c.IsSet(flag) {
			return cli.NewExitError(
				fmt.Sprintf("flag --%s is required for provider: %s", flag, cloudProvider),
				1,
			)
		}
	}

	// verify alerting flags are used correctly
	alertingFlagsOkay := true
	if c.Bool(alertingFlag) {
		for _, f := range providerAlertingRequiredFlags[cloudProvider] {
			if !c.IsSet(f) {
				alertingFlagsOkay = false
				msg := fmt.Sprintf("--%s is required\n", f)
				if _, err := fmt.Fprint(cli.ErrWriter, msg); err != nil {
					return cli.NewExitError(msg, 1)
				}
			}
		}
		if !alertingFlagsOkay {
			return cli.NewExitError("missing required flags for --alerting", 1)
		}
	} else {
		for _, f := range providerAlertingFlags[cloudProvider] {
			if c.IsSet(f) {
				alertingFlagsOkay = false
				msg := fmt.Sprintf("--%s is disallowed\n", f)
				if _, err := fmt.Fprint(cli.ErrWriter, msg); err != nil {
					return cli.NewExitError(msg, 1)
				}
			}
		}
		if !alertingFlagsOkay {
			return cli.NewExitError("disallowed flags specified with missing --alerting", 1)
		}
	}

	return nil
}

func generateModule(c *cli.Context) error {
	if err := validateArgs(c); err != nil {
		return err
	}

	configFile := c.Args().Get(0)

	config, err := newConfig(c, configFile)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	writer := newConfigWriter(c, config)
	err = writer.writeTerraform()
	if err != nil {
		return cli.NewExitError(errors.Wrap(err, "failed to generate terraform module"), 1)
	}

	fmt.Println("Created Terraform Hedwig module successfully!")
	return nil
}

func generateConfigFileStructure(c *cli.Context) error {
	cloudProvider := c.GlobalString(cloudProviderFlag)
	if cloudProvider == "" {
		return cli.NewExitError(fmt.Sprintf("--%s is required", cloudProviderFlag), 1)
	}
	if cloudProvider != cloudProviderAWS && cloudProvider != cloudProviderGoogle {
		return cli.NewExitError(fmt.Sprintf("invalid cloud provider: %s", cloudProvider), 1)
	}

	var structure interface{}
	if cloudProvider == cloudProviderAWS {
		structure = AWSConfig{
			Topics: []string{
				"my-topic",
			},
			QueueConsumers: []*AWSQueueConsumer{
				{
					"DEV-MYAPP",
					map[string]string{
						"App": "myapp",
						"Env": "dev",
					},
					[]string{"my-topic"},
				},
			},
			LambdaConsumers: []*AWSLambdaConsumer{
				{
					"arn:aws:lambda:us-west-2:12345:function:my_function:deployed",
					"{optional - this value is inferred from FunctionARN if that's not an interpolated value}",
					"{optional - this value is inferred from FunctionARN if that's not an interpolated value}",
					[]string{"my-topic"},
				},
			},
		}
	} else if cloudProvider == cloudProviderGoogle {
		structure = GoogleConfig{
			Topics: []*GoogleTopic{
				{
					"my-topic",
					false,
				},
			},
			PullConsumers: []*GooglePullConsumer{
				{
					"dev-myapp",
					[]string{"my-topic"},
					map[string]string{
						"App": "myapp",
						"Env": "dev",
					},
				},
			},
		}
	}
	structureAsJSON, err := json.MarshalIndent(structure, "", "    ")
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	fmt.Println(string(structureAsJSON))
	return nil
}

func runApp(args []string) error {
	cli.VersionFlag = cli.BoolFlag{Name: "version, V"}

	app := cli.NewApp()
	app.Name = "Hedwig Terraform"
	app.Usage = "Manage Terraform configuration for Hedwig apps"
	app.Version = VERSION
	app.HelpName = "hedwig-terraform"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  cloudProviderFlag,
			Usage: "Cloud provider - either aws or google",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "generate",
			Usage:     "Generates Terraform module for Hedwig apps",
			ArgsUsage: "<config-file>",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  moduleFlag,
					Usage: "Terraform module name to generate",
					Value: "hedwig",
				},
				cli.BoolFlag{
					Name:  alertingFlag,
					Usage: "Should alerting be generated?",
				},
				cli.StringSliceFlag{
					Name:  queueAlertAlarmActionsFlag,
					Usage: "Cloudwatch Action ARNs for high message count in queue when in ALARM (AWS only)",
				},
				cli.StringSliceFlag{
					Name:  queueAlertNotificationChannelsFlag,
					Usage: "Stackdriver Notification Channels for high message count in queue (Google only)",
				},
				cli.StringSliceFlag{
					Name:  queueAlertOKActionsFlag,
					Usage: "Cloudwatch Action ARNs for high message count in queue when OK (AWS only)",
				},
				cli.StringSliceFlag{
					Name:  dlqAlertAlarmActionsFlag,
					Usage: "Cloudwatch Action ARNs for high message count in dead-letter queue when in ALARM (AWS only)",
				},
				cli.StringSliceFlag{
					Name:  dlqAlertNotificationChannelsFlag,
					Usage: "Stackdriver Notification Channels for high message count in dead-letter queue (Google only)",
				},
				cli.StringSliceFlag{
					Name:  dlqAlertOKActionsFlag,
					Usage: "Cloudwatch Action ARNs for high message count in dead-letter queue when OK (AWS only)",
				},
				cli.StringFlag{
					Name:  dataflowTmpGCSLocationFlag,
					Usage: "Dataflow tmp GCS location (Google only) (required for firehose)",
				},
				cli.StringFlag{
					Name:  dataflowPubSubToPubSubTemplateGCSPathFlag,
					Usage: "Dataflow template for pubsub to pubsub GCS location (Google only) (required for firehose)",
				},
				cli.StringFlag{
					Name:  dataflowPubSubToStorageGCSPathFlag,
					Usage: "Dataflow template for pubsub to storage GCS location (Google only) (required for firehose)",
				},
				cli.BoolFlag{
					Name:  enableFirehoseAllTopics,
					Usage: "Enable Google Firehose for all hedwig messages (Google only for now)",
				},
				cli.StringFlag{
					Name:  googleFirehoseDataflowOutputDirectoryFlag,
					Usage: "Google Firehose Dataflow output directory. Must end with /. (Google only)",
				},
				cli.StringFlag{
					Name: googleDataflowZoneFlag,
					Usage: "Dataflow zone (Google only) (required for firehose if it's not set at the provider level, " +
						"or that zone doesn't support Dataflow regional endpoints (see " +
						"https://cloud.google.com/dataflow/docs/concepts/regional-endpoints)",
				},
				cli.StringFlag{
					Name: googleDataflowRegionFlag,
					Usage: "Dataflow zone (Google only) (required for firehose if it's not set at the provider level, " +
						"or you want to use a region different from the zone (see " +
						"https://cloud.google.com/dataflow/docs/concepts/regional-endpoints)",
				},
				cli.StringFlag{
					Name: googleProjectAlerting,
					Usage: "Google project to use for alerting resources. This may be different from your main" +
						"app environment (Google only)",
				},
				cli.StringFlag{
					Name:  awsAccountIDFlag,
					Usage: "AWS Account ID (AWS only) (required)",
				},
				cli.StringFlag{
					Name:  awsRegionFlag,
					Usage: "AWS Region (AWS only) (required)",
				},
			},
			Action: generateModule,
		},
		{
			Name:   "config-file-structure",
			Usage:  "Outputs the structure for config file required for generate command",
			Action: generateConfigFileStructure,
		},
	}

	return app.Run(args)
}

func main() {
	if err := runApp(os.Args); err != nil {
		log.Fatal(err)
	}
}
