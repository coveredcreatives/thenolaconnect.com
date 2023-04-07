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
    NOLA_ENV                                      = "production"
    NOLA_DB_USERNAME                              = var.db_username
    NOLA_DB_PASSWORD                              = var.db_password
    NOLA_DB_NAME                                  = var.db_name
    NOLA_DB_PORT                                  = var.db_port
    NOLA_DB_HOSTNAME                              = format("/cloudsql/%s", google_sql_database_instance.company_database_instance.connection_name)
    NOLA_HTTP_PORT                                = var.http_port
    NOLA_DNS_PRINTER_IPV4_ADDRESS                 = var.dns_printer_ipv4_address
    NOLA_GOOGLE_FORM_ID_ORDERS                    = var.google_form_id_orders
    NOLA_GOOGLE_APPLICATION_SERVICE_ACCOUNT_EMAIL = var.google_application_service_account_email
    NOLA_GOOGLE_STORAGE_BUCKET_NAME               = var.google_storage_bucket_name
    NOLA_TWILIO_ACCOUNT_SID                       = var.twilio_account_sid
    NOLA_TWILIO_CONVERSATION_SERVICE_SID          = var.twilio_conversation_service_sid
    NOLA_TWILIO_ACCOUNT_AUTH_TOKEN                = var.twilio_account_auth_token
    NOLA_RETRIEVE_HTTPS_TRIGGER_URL               = google_cloudfunctions_function.qr_code_retrieve.https_trigger_url
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
    NOLA_ENV                                      = "production"
    NOLA_DB_USERNAME                              = var.db_username
    NOLA_DB_PASSWORD                              = var.db_password
    NOLA_DB_NAME                                  = var.db_name
    NOLA_DB_PORT                                  = var.db_port
    NOLA_DB_HOSTNAME                              = format("/cloudsql/%s", google_sql_database_instance.company_database_instance.connection_name)
    NOLA_HTTP_PORT                                = var.http_port
    NOLA_DNS_PRINTER_IPV4_ADDRESS                 = var.dns_printer_ipv4_address
    NOLA_GOOGLE_FORM_ID_ORDERS                    = var.google_form_id_orders
    NOLA_GOOGLE_APPLICATION_SERVICE_ACCOUNT_EMAIL = var.google_application_service_account_email
    NOLA_GOOGLE_STORAGE_BUCKET_NAME               = var.google_storage_bucket_name
    NOLA_TWILIO_ACCOUNT_SID                       = var.twilio_account_sid
    NOLA_TWILIO_CONVERSATION_SERVICE_SID          = var.twilio_conversation_service_sid
    NOLA_TWILIO_ACCOUNT_AUTH_TOKEN                = var.twilio_account_auth_token
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
