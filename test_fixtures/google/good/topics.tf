// DO NOT EDIT
// This file has been auto-generated by hedwig-terraform-generator {{GENERATOR_VERSION}}

module "topic-my-topic" {
  source  = "standard-ai/hedwig-topic/google"
  version = "~> {{TFGoogleTopicModuleVersion}}"

  topic = "my-topic"

  iam_service_accounts = [
    "secondapp@project.iam.gserviceaccount.com",
  ]

  enable_firehose_all_messages = var.enable_firehose_all_topics
  dataflow_tmp_gcs_location    = var.dataflow_tmp_gcs_location
  dataflow_template_gcs_path   = var.dataflow_template_pubsub_to_storage_gcs_path
  dataflow_zone                = var.dataflow_zone
  dataflow_region              = var.dataflow_region
  dataflow_output_directory    = var.dataflow_output_directory

  enable_alerts    = var.enable_alerts
  alerting_project = var.alerting_project
}

module "topic-my-topic-2" {
  source  = "standard-ai/hedwig-topic/google"
  version = "~> {{TFGoogleTopicModuleVersion}}"

  topic = "my-topic2"

  iam_service_accounts = [
    "thirdapp@project.iam.gserviceaccount.com", "fourthapp@project.iam.gserviceaccount.com",
  ]

  enable_firehose_all_messages = var.enable_firehose_all_topics
  dataflow_tmp_gcs_location    = var.dataflow_tmp_gcs_location
  dataflow_template_gcs_path   = var.dataflow_template_pubsub_to_storage_gcs_path
  dataflow_zone                = var.dataflow_zone
  dataflow_region              = var.dataflow_region
  dataflow_output_directory    = var.dataflow_output_directory

  enable_alerts    = var.enable_alerts
  alerting_project = var.alerting_project
}
