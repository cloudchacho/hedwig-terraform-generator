// DO NOT EDIT
// This file has been auto-generated by hedwig-terraform-generator {{GENERATOR_VERSION}}

module "consumer-dev-myapp" {
  source  = "Automatic/hedwig-queue/aws"
  version = "~> {{TFAWSQueueModuleVersion}}"

  queue          = "DEV-MYAPP"
  aws_region     = "${var.aws_region}"
  aws_account_id = "${var.aws_account_id}"
  alerting       = "true"

  tags = {
    App = "myapp"
    Env = "dev"
  }

  dlq_alarm_high_message_count_actions = [
    "pager_action",
    "pager_action2",
  ]

  dlq_ok_high_message_count_actions = [
    "pager_action",
    "pager_action2",
  ]

  queue_alarm_high_message_count_actions = [
    "pager_action",
    "pager_action2",
  ]

  queue_ok_high_message_count_actions = [
    "pager_action",
    "pager_action2",
  ]

  queue_alarm_high_message_count_threshold = 100000
}

module "consumer-dev-secondapp" {
  source  = "Automatic/hedwig-queue/aws"
  version = "~> {{TFAWSQueueModuleVersion}}"

  queue          = "DEV-SECONDAPP"
  aws_region     = "${var.aws_region}"
  aws_account_id = "${var.aws_account_id}"
  alerting       = "true"

  tags = {
    App = "secondapp"
    Env = "dev"
  }

  dlq_alarm_high_message_count_actions = [
    "pager_action",
    "pager_action2",
  ]

  dlq_ok_high_message_count_actions = [
    "pager_action",
    "pager_action2",
  ]

  queue_alarm_high_message_count_actions = [
    "pager_action",
    "pager_action2",
  ]

  queue_ok_high_message_count_actions = [
    "pager_action",
    "pager_action2",
  ]

  queue_alarm_high_message_count_threshold = 10000
}
