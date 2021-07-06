// DO NOT EDIT
// This file has been auto-generated by hedwig-terraform-generator {{GENERATOR_VERSION}}

module "alerts-dev-myapp" {
  source  = "cloudchacho/hedwig-alerts/google"
  version = "~> {{TFGoogleAlertsModuleVersion}}"

  subscription_name     = module.consumer-dev-myapp.subscription_name
  dlq_subscription_name = module.consumer-dev-myapp.dlq_subscription_name

  labels = {
    app = "myapp"
    env = "dev"
  }

  queue_high_message_count_notification_channels = [
    "projects/myProject/notificationChannels/10357685029951383687",
    "projects/myProject/notificationChannels/95138368710357685029"
  ]
  dlq_high_message_count_notification_channels = [
    "projects/myProject/notificationChannels/10357685029951383687",
    "projects/myProject/notificationChannels/95138368710357685029"
  ]
  queue_alarm_high_message_count_threshold = 10000
}

module "alerts-dev-myapp-my-topic" {
  source  = "cloudchacho/hedwig-alerts/google"
  version = "~> {{TFGoogleAlertsModuleVersion}}"

  subscription_name = module.sub-dev-myapp-my-topic.subscription_name

  labels = {
    app = "myapp"
    env = "dev"
  }

  queue_high_message_count_notification_channels = [
    "projects/myProject/notificationChannels/10357685029951383687",
    "projects/myProject/notificationChannels/95138368710357685029"
  ]
  dlq_high_message_count_notification_channels = [
    "projects/myProject/notificationChannels/10357685029951383687",
    "projects/myProject/notificationChannels/95138368710357685029"
  ]
  queue_alarm_high_message_count_threshold = 10000
}

module "alerts-dev-myapp-my-topic-2" {
  source  = "cloudchacho/hedwig-alerts/google"
  version = "~> {{TFGoogleAlertsModuleVersion}}"

  subscription_name = module.sub-dev-myapp-my-topic-2.subscription_name

  labels = {
    app = "myapp"
    env = "dev"
  }

  queue_high_message_count_notification_channels = [
    "projects/myProject/notificationChannels/10357685029951383687",
    "projects/myProject/notificationChannels/95138368710357685029"
  ]
  dlq_high_message_count_notification_channels = [
    "projects/myProject/notificationChannels/10357685029951383687",
    "projects/myProject/notificationChannels/95138368710357685029"
  ]
  queue_alarm_high_message_count_threshold = 10000
}

module "alerts-dev-secondapp" {
  source  = "cloudchacho/hedwig-alerts/google"
  version = "~> {{TFGoogleAlertsModuleVersion}}"

  subscription_name     = module.consumer-dev-secondapp.subscription_name
  dlq_subscription_name = module.consumer-dev-secondapp.dlq_subscription_name

  labels = {
    app = "secondapp"
    env = "dev"
  }

  queue_high_message_count_notification_channels = [
    "projects/myProject/notificationChannels/10357685029951383687",
    "projects/myProject/notificationChannels/95138368710357685029"
  ]
  dlq_high_message_count_notification_channels = [
    "projects/myProject/notificationChannels/10357685029951383687",
    "projects/myProject/notificationChannels/95138368710357685029"
  ]
  queue_alarm_high_message_count_threshold = 10000
}

module "alerts-dev-secondapp-my-topic-2" {
  source  = "cloudchacho/hedwig-alerts/google"
  version = "~> {{TFGoogleAlertsModuleVersion}}"

  subscription_name = module.sub-dev-secondapp-my-topic-2.subscription_name

  labels = {
    app = "secondapp"
    env = "dev"
  }

  queue_high_message_count_notification_channels = [
    "projects/myProject/notificationChannels/10357685029951383687",
    "projects/myProject/notificationChannels/95138368710357685029"
  ]
  dlq_high_message_count_notification_channels = [
    "projects/myProject/notificationChannels/10357685029951383687",
    "projects/myProject/notificationChannels/95138368710357685029"
  ]
  queue_alarm_high_message_count_threshold = 10000
}
