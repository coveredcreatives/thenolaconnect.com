variable "organization_id" {
  description = "The organization ID in Google Cloud to use for these resources."
}

variable "project_id" {
  description = "The project ID in Google Cloud to use for these resources."
}

variable "region" {
  description = "The region in Google Cloud where the resources will be deployed."
}

variable "database_password" {
  description = "Password to connect to database instance"
  sensitive   = true
}

variable "domain" {
  description = "top level main domain of project"
}

variable "admin_email" {
  description = "The email for creating iam roles"
}

variable "db_username" {
  description = "PGUSER of postgres database connection string"
}

variable "db_password" {
  description = "PGPASS of postgres database connection string"
  sensitive   = true
}

variable "db_name" {
  description = "PGDATABASE of postgres database connection string"
}

variable "db_port" {
  description = "PGPORT of postgres database connection string"
}

variable "http_port" {
  description = "assigns a number for the application to listen to TCP ports"
  default     = 8080
}

variable "dns_printer_ipv4_address" {
  description = "assign an ip address to deliver outgoing printer messages to"
}

variable "google_form_id_orders" {
  description = "define a form id for the application to READ form data from answer, form_response, form, and item records in google forms API"
}

variable "google_application_service_account_email" {
  description = "define the email for project application service account to execute on behalf of"
}

variable "google_storage_bucket_name" {
  description = "define the storage bucket to store generated QR and associated media objects"
}

variable "twilio_account_sid" {
  description = "assign twilio account sid for application to access conversations api"
}

variable "twilio_conversation_service_sid" {
  description = "define twilio conversation service application should list and create conversations to"
}

variable "twilio_account_auth_token" {
  description = "assign auth token that enabled the application to preform actions on the twilio account"
}
