terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
  }
  required_version = ">= 1.5.0"
}

# Initial provider (you only need gcloud auth the very first time)
provider "google" {
  credentials = file("${path.module}/sa-key.json")
  project     = var.project_id
  region      = "us-central1"
}

# Enable required APIs
resource "google_project_service" "gmail_api" {
  project = var.project_id
  service = "gmail.googleapis.com"
}

resource "google_project_service" "pubsub_api" {
  project = var.project_id
  service = "pubsub.googleapis.com"
}

# Create a Service Account for your app
resource "google_service_account" "app_sa" {
  account_id   = "gmail-pubsub-sa"
  display_name = "Service Account for Gmail PubSub"
}

# Create a JSON key for the Service Account
resource "google_service_account_key" "app_sa_key" {
  service_account_id = google_service_account.app_sa.name
}

# Save the key JSON locally (you’ll mount this into your container later)
resource "local_file" "service_account_key" {
  content  = base64decode(google_service_account_key.app_sa_key.private_key)
  filename = "${path.module}/sa-key.json"
}

# Pub/Sub Topic
resource "google_pubsub_topic" "gmail_notifications" {
  name = "gmail-notifications"
}

# Allow Gmail API service account to publish to the topic
resource "google_pubsub_topic_iam_member" "gmail_publisher" {
  topic  = google_pubsub_topic.gmail_notifications.name
  role   = "roles/pubsub.publisher"
  member = "serviceAccount:gmail-api-push@system.gserviceaccount.com"
}

# Pub/Sub Subscription
resource "google_pubsub_subscription" "gmail_subscription" {
  name  = "gmail-notifications-sub"
  topic = google_pubsub_topic.gmail_notifications.name

  ack_deadline_seconds = 10
}

# Allow your Service Account to consume messages
resource "google_pubsub_subscription_iam_member" "sub_reader" {
  subscription = google_pubsub_subscription.gmail_subscription.name
  role         = "roles/pubsub.subscriber"
  member       = "serviceAccount:${google_service_account.app_sa.email}"
}
