resource "google_compute_global_address" "peering_address" {
  name          = "default-vpc-sql"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  project       = var.gcloud_project_id
  prefix_length = 24
  network       = google_compute_network.vpc.self_link
}

resource "google_service_networking_connection" "vpc_peering" {
  network                 = google_compute_network.vpc.self_link
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.peering_address.name]
}

resource "google_sql_database_instance" "primary" {
  name             = "primary-db"
  database_version = "POSTGRES_14"
  region           = var.gcloud_region
  settings {
    tier              = var.db_tier
    availability_type = "ZONAL"
    activation_policy = "ALWAYS"
    ip_configuration {
      private_network = google_compute_network.vpc.self_link

    }
  }
  depends_on = [
    google_service_networking_connection.vpc_peering
  ]
}

resource "google_sql_user" "user" {
  instance = google_sql_database_instance.primary.name
  name     = var.db_user
  password = var.db_password
}

resource "google_sql_database" "main" {
  instance = google_sql_database_instance.primary.name
  name     = var.db_name
}