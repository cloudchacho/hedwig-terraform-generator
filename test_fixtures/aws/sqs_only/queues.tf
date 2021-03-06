// DO NOT EDIT
// This file has been auto-generated by hedwig-terraform-generator {{GENERATOR_VERSION}}

module "consumer-dev-myapp" {
  source  = "Automatic/hedwig-queue/aws"
  version = "~> {{TFAWSQueueModuleVersion}}"

  queue          = "DEV-MYAPP"
  aws_region     = "${var.aws_region}"
  aws_account_id = "${var.aws_account_id}"

  tags = {
    App = "myapp"
    Env = "dev"
  }
}
