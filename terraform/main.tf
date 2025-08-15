terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
  }
  required_version = ">= 1.5.0"
}

provider "google" {
  project = var.project_id
  region  = "us-central1"
}

resource "google_project_service" "gmail_api" {
  project = var.project_id
  service = "gmail.googleapis.com"
}

resource "google_project_service" "pubsub_api" {
  project = var.project_id
  service = "pubsub.googleapis.com"
}

resource "google_pubsub_topic" "gmail_notifications" {
  name = "gmail-notifications"
}

resource "google_pubsub_topic_iam_member" "gmail_publisher" {
  topic  = google_pubsub_topic.gmail_notifications.name
  role   = "roles/pubsub.publisher"
  member = "serviceAccount:gmail-api-push@system.gserviceaccount.com"
}

resource "google_pubsub_subscription" "gmail_subscription" {
  name  = "gmail-notifications-sub"
  topic = google_pubsub_topic.gmail_notifications.name

  ack_deadline_seconds = 10
}
