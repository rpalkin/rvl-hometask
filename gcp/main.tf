terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "4.83.0"
    }
  }
}

provider "google" {
  project = var.gcloud_project_id
  region  = var.gcloud_region
}




