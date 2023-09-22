resource "google_compute_network" "vpc" {
  name = "vpc1"
  auto_create_subnetworks = false
}