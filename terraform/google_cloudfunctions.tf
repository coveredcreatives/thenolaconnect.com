resource "google_cloudfunctions_function" "qr_code_generate" {
  project               = data.google_project.new_orleans_connection.project_id
  name                  = "generate"
  description           = "receives a name + file via http and returns a QR code rendering that file."
  region                = var.region
  runtime               = "go119"
  trigger_http          = true
  entry_point           = "GenerateQRHandler" # Set the entry point 
  source_archive_bucket = google_storage_bucket.functions.name
  source_archive_object = google_storage_bucket_object.primary_server_zip.name
  service_account_email = local.service_account_email
  environment_variables = {
    DB_USERNAME                = "postgres"
    DB_PASSWORD                = var.database_password
    DB_NAME                    = google_sql_database.company_database.name
    DB_PORT                    = 5432
    DB_HOSTNAME                = google_sql_database_instance.company_database_instance.public_ip_address
    INSTANCE_CONNECTION_NAME   = google_sql_database_instance.company_database_instance.connection_name
    UNIX_SOCKET_PATH           = format("/cloudsql/%s", google_sql_database_instance.company_database_instance.connection_name)
    RETRIEVE_HTTPS_TRIGGER_URL = google_cloudfunctions_function.qr_code_retrieve.https_trigger_url
  }
}

resource "google_cloudfunctions_function" "qr_code_retrieve" {
  project               = data.google_project.new_orleans_connection.project_id
  name                  = "retrieve"
  description           = "lists metadata of stored qr codes or returns image for single lookup"
  region                = var.region
  runtime               = "go119"
  trigger_http          = true
  entry_point           = "RetrieveQRHandler" # Set the entry point 
  source_archive_bucket = google_storage_bucket.functions.name
  source_archive_object = google_storage_bucket_object.primary_server_zip.name
  service_account_email = local.service_account_email
  environment_variables = {
    DB_USERNAME              = "postgres"
    DB_PASSWORD              = var.database_password
    DB_NAME                  = google_sql_database.company_database.name
    DB_PORT                  = 5432
    DB_HOSTNAME              = google_sql_database_instance.company_database_instance.public_ip_address
    INSTANCE_CONNECTION_NAME = google_sql_database_instance.company_database_instance.connection_name
    UNIX_SOCKET_PATH         = format("/cloudsql/%s", google_sql_database_instance.company_database_instance.connection_name)
  }
}

resource "google_cloudfunctions_function_iam_member" "invoker_generate" {
  project        = data.google_project.new_orleans_connection.project_id
  cloud_function = google_cloudfunctions_function.qr_code_generate.name

  region = var.region
  role   = "roles/cloudfunctions.invoker"
  member = "allUsers"
}

resource "google_cloudfunctions_function_iam_member" "invoker_retrieve" {
  project        = data.google_project.new_orleans_connection.project_id
  cloud_function = google_cloudfunctions_function.qr_code_retrieve.name

  region = var.region
  role   = "roles/cloudfunctions.invoker"
  member = "allUsers"
}
