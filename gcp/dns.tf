resource "google_dns_managed_zone" "birthday_zone" {
  name        = "birthday-zone"
  dns_name    = var.dns_zone
  visibility = "public"
}

resource "google_service_account" "dns_editor" {
    account_id   = "dns-editor"
}

resource "google_project_iam_member" "dns_editor" {
  project = var.gcloud_project_id
  role    = "roles/dns.admin"
  member = "serviceAccount:${google_service_account.dns_editor.email}"
}

resource "google_service_account_iam_member" "dns_editor" {
  service_account_id = google_service_account.dns_editor.name
  role               = "roles/iam.workloadIdentityUser"
  member = "serviceAccount:${var.gcloud_project_id}.svc.id.goog[external-dns/dns-editor]"
}