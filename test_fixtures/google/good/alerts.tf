// DO NOT EDIT
// This file has been auto-generated by hedwig-terraform-generator {{GENERATOR_VERSION}}

module "alerts-dev-myapp" {
  source  = "standard-ai/hedwig-alerts/google"
  version = "~> {{TFGoogleAlertsModuleVersion}}"

  subscription_name     = module.consumer-dev-myapp.subscription_name
  dlq_subscription_name = module.consumer-dev-myapp.dlq_subscription_name

  alerting_project = var.alerting_project

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
}

module "alerts-dev-myapp-my-topic" {
  source  = "standard-ai/hedwig-alerts/google"
  version = "~> {{TFGoogleAlertsModuleVersion}}"

  subscription_name = module.sub-dev-myapp-my-topic.subscription_name

  alerting_project = var.alerting_project

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
}

module "alerts-dev-myapp-my-topic-2" {
  source  = "standard-ai/hedwig-alerts/google"
  version = "~> {{TFGoogleAlertsModuleVersion}}"

  subscription_name = module.sub-dev-myapp-my-topic-2.subscription_name

  alerting_project = var.alerting_project

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
}

module "alerts-dev-secondapp" {
  source  = "standard-ai/hedwig-alerts/google"
  version = "~> {{TFGoogleAlertsModuleVersion}}"

  subscription_name     = module.consumer-dev-secondapp.subscription_name
  dlq_subscription_name = module.consumer-dev-secondapp.dlq_subscription_name

  alerting_project = var.alerting_project

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
}

module "alerts-dev-secondapp-my-topic-2" {
  source  = "standard-ai/hedwig-alerts/google"
  version = "~> {{TFGoogleAlertsModuleVersion}}"

  subscription_name = module.sub-dev-secondapp-my-topic-2.subscription_name

  alerting_project = var.alerting_project

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
}

module "alerts-dev-secondapp-other-project-my-topic-3" {
  source  = "standard-ai/hedwig-alerts/google"
  version = "~> {{TFGoogleAlertsModuleVersion}}"

  subscription_name = module.sub-dev-secondapp-other-project-my-topic-3.subscription_name

  alerting_project = var.alerting_project

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
}
