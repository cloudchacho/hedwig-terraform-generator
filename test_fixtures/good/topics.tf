// DO NOT EDIT
// This file has been auto-generated by hedwig-terraform-generator {{GENERATOR_VERSION}}

module "topic-my-topic" {
  source  = "Automatic/hedwig-topic/aws"
  version = "~> {{TFTopicModuleVersion}}"

  topic = "my-topic"
}

module "topic-my-topic-2" {
  source  = "Automatic/hedwig-topic/aws"
  version = "~> {{TFTopicModuleVersion}}"

  topic = "my-topic2"
}
