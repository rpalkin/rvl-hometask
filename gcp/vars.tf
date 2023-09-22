variable "gcloud_project_id" {
  type = string
  description = "The GCP project ID"
}

variable "gcloud_region" {
  type = string
  description = "The GCP region"
}

variable "db_tier" {
  type = string
  default = "db-f1-micro"
  description = "The GCP Cloud SQL tier"
}

variable "db_user" {
  type = string
  description = "The GCP Cloud SQL user"
}

variable "db_password" {
  type = string
  description = "The GCP Cloud SQL password"
}

variable "db_name" {
  type = string
  description = "The GCP Cloud SQL database name"
}

variable "dns_zone" {
  type = string
  description = "The DNS zone name"
}