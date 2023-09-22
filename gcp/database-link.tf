#provider "kubernetes" {
#  host                   = "https://${google_container_cluster.birthday_cluster.endpoint}"
#  token                  = data.google_client_config.current.access_token
#  cluster_ca_certificate = base64decode(google_container_cluster.birthday_cluster.master_auth.0.cluster_ca_certificate)
#}
#
#resource "kubernetes_namespace" "postgres" {
#  metadata {
#    name = var.db_namespace
#  }
#}
#
#resource "kubernetes_service" "postgres" {
#    metadata {
#        name = var.db_svc_name
#        namespace = kubernetes_namespace.postgres.metadata[0].name
#    }
#    spec {
#        type = "ExternalName"
#        external_name = google_sql_database_instance.primary.private_ip_address
#        port {
#          port = 5432
#        }
#    }
#}