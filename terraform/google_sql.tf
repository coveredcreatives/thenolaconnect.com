# primary business operations database instance
resource "google_sql_database_instance" "company_database_instance" {
  project             = data.google_project.new_orleans_connection.project_id
  name                = "qr-code-instance"
  region              = var.region
  database_version    = "POSTGRES_14"
  root_password       = var.db_password
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    database_flags {
      name  = "cloudsql.iam_authentication"
      value = "on"
    }

    ip_configuration {
      private_network                               = google_compute_network.vpc.id
      enable_private_path_for_google_cloud_services = true
      authorized_networks {
        name  = "onprem-devscrum"
        value = "98.163.211.86"
      }
    }


  }

}

# qr mappings database
resource "google_sql_database" "company_database" {
  project  = data.google_project.new_orleans_connection.project_id
  instance = google_sql_database_instance.company_database_instance.name
  name     = "qr_mappings"
}

resource "google_sql_user" "default" {
  project  = data.google_project.new_orleans_connection.project_id
  instance = google_sql_database_instance.company_database_instance.name
  name     = "postgres"
  password = var.db_password
  type     = "BUILT_IN"
}

# sql user admin uses iam for login access
resource "google_sql_user" "admin" {
  project  = data.google_project.new_orleans_connection.project_id
  instance = google_sql_database_instance.company_database_instance.name
  name     = var.admin_email
  type     = "CLOUD_IAM_USER"
}

# sql user service account for cloud function access
resource "google_sql_user" "service_account" {
  project  = data.google_project.new_orleans_connection.project_id
  instance = google_sql_database_instance.company_database_instance.name
  name     = trimsuffix(local.service_account_email, "gserviceaccount.com")
  type     = "CLOUD_IAM_SERVICE_ACCOUNT"
}
