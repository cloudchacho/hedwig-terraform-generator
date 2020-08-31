// DO NOT EDIT
// This file has been auto-generated by hedwig-terraform-generator {{GENERATOR_VERSION}}

module "sub-dev-myapp-my-topic" {
  source  = "standard-ai/hedwig-subscription/google"
  version = "~> {{TFGoogleSubscriptionModuleVersion}}"

  queue = "dev-myapp"
  topic = module.topic-my-topic.name

  iam_service_account = "myapp@project.iam.gserviceaccount.com"

  labels = {
    app = "myapp"
    env = "dev"
  }
}

module "sub-dev-myapp-my-topic-2" {
  source  = "standard-ai/hedwig-subscription/google"
  version = "~> {{TFGoogleSubscriptionModuleVersion}}"

  queue = "dev-myapp"
  topic = module.topic-my-topic-2.name

  iam_service_account = "myapp@project.iam.gserviceaccount.com"

  labels = {
    app = "myapp"
    env = "dev"
  }
}

module "sub-dev-secondapp-my-topic-2" {
  source  = "standard-ai/hedwig-subscription/google"
  version = "~> {{TFGoogleSubscriptionModuleVersion}}"

  queue = "dev-secondapp"
  topic = module.topic-my-topic-2.name

  iam_service_account = "secondapp@project.iam.gserviceaccount.com"

  labels = {
    app = "secondapp"
    env = "dev"
  }
}

module "sub-dev-secondapp-other-project-my-topic-3" {
  source  = "standard-ai/hedwig-subscription/google"
  version = "~> {{TFGoogleSubscriptionModuleVersion}}"

  queue = "dev-secondapp"
  topic = "projects/other-project/topics/hedwig-my-topic3"

  iam_service_account = "secondapp@project.iam.gserviceaccount.com"

  labels = {
    app = "secondapp"
    env = "dev"
  }
}
