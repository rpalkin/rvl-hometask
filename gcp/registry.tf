resource "google_artifact_registry_repository" "registry" {
  format        = "DOCKER"
  repository_id = "docker-registry"
  location = var.gcloud_region
}