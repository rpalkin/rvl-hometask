resource "google_compute_subnetwork" "gke_subnet" {
  name          = "gke-subnet"
  ip_cidr_range = "10.0.0.0/15"
  network       = google_compute_network.vpc.name
}

resource "google_service_account" "default" {
  account_id   = "nodes-sa"
}

resource "google_container_cluster" "birthday_cluster" {
  name = "birthday-cluster"
  location = var.gcloud_region

  initial_node_count       = 1
  remove_default_node_pool = true

  networking_mode = "VPC_NATIVE"
  datapath_provider = "ADVANCED_DATAPATH"

  network = google_compute_network.vpc.name
  subnetwork = google_compute_subnetwork.gke_subnet.name

  ip_allocation_policy {
    cluster_ipv4_cidr_block  = "/16"
    services_ipv4_cidr_block = "/22"
  }

  workload_identity_config {
    workload_pool = "${var.gcloud_project_id}.svc.id.goog"
  }

  release_channel {
    channel = "STABLE"
  }

}

data "google_client_config" "current" {}

resource "google_container_node_pool" "birthday-pool-one" {
  name = "birthday-pool-one"
  cluster = google_container_cluster.birthday_cluster.name
  location = var.gcloud_region
  initial_node_count = 1
  node_config {
    workload_metadata_config {
      mode = "GKE_METADATA"
    }
    disk_size_gb = 30
    machine_type = "e2-small"
    service_account = google_service_account.default.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}

resource "google_artifact_registry_repository_iam_binding" "node_iam_binding" {
  repository = google_artifact_registry_repository.registry.name
  location   = var.gcloud_region
  role       = "roles/artifactregistry.reader"
  project    = var.gcloud_project_id

  members = [
    "serviceAccount:${google_service_account.default.email}"
  ]
}