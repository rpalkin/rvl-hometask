output "registry_location" {
  value = "${google_artifact_registry_repository.registry.location}-docker.pkg.dev/${google_artifact_registry_repository.registry.project}/${google_artifact_registry_repository.registry.name}/"
}

output "dns_servers" {
  value = google_dns_managed_zone.birthday_zone.name_servers
}

output "dns_sa" {
  value = google_service_account.dns_editor.email
}

output "cloudsql_host" {
  value = google_sql_database_instance.primary.private_ip_address
}