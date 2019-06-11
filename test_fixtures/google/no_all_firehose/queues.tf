// DO NOT EDIT
// This file has been auto-generated by hedwig-terraform-generator {{GENERATOR_VERSION}}

module "consumer-dev-myapp" {
  source  = "standard-ai/hedwig-queue/google"
  version = "~> {{TFGoogleQueueModuleVersion}}"

  queue    = "dev-myapp"
  alerting = "true"

  labels = {
    app = "myapp"
    env = "dev"
  }

  queue_high_message_count_notification_channels = [
    "projects/myProject/notificationChannels/10357685029951383687",
    "projects/myProject/notificationChannels/95138368710357685029",
  ]

  dlq_high_message_count_notification_channels = [
    "projects/myProject/notificationChannels/10357685029951383687",
    "projects/myProject/notificationChannels/95138368710357685029",
  ]
}

module "consumer-dev-secondapp" {
  source  = "standard-ai/hedwig-queue/google"
  version = "~> {{TFGoogleQueueModuleVersion}}"

  queue    = "dev-secondapp"
  alerting = "true"

  labels = {
    app = "secondapp"
    env = "dev"
  }

  queue_high_message_count_notification_channels = [
    "projects/myProject/notificationChannels/10357685029951383687",
    "projects/myProject/notificationChannels/95138368710357685029",
  ]

  dlq_high_message_count_notification_channels = [
    "projects/myProject/notificationChannels/10357685029951383687",
    "projects/myProject/notificationChannels/95138368710357685029",
  ]
}