# begin variables
variable "config_bucket" {
  description = "The bucket to store the terraform configration data"
}

variable "organization_id" {
  description = "The organization ID in Google Cloud to use for these resources."
}

variable "project_id" {
  description = "The project ID in Google Cloud to use for these resources."
}

variable "region" {
  description = "The region in Google Cloud where the resources will be deployed."
}

# begin config
terraform {
  backend "gcs" {
    bucket = "the-new-orleans-connection-terraform-state" # Bucket is passed in via cli arg. Eg, terraform init -reconfigure -backend-configuration=dev.tfbackend
  }
  required_providers {
    twilio = {
      source  = "twilio/twilio"
      version = "0.18.12"
    }
  }
}

provider "google-beta" {
  project = var.project_id
  region  = var.region
}

# initialize google provider and establish project + region
provider "google" {
  project     = var.project_id
  region      = var.region
}

provider "twilio" {
  username = var.twilio_api_key
  password = var.twilio_api_secret
}

# local directory where source is stored
locals {
  primary_server_dir = abspath("../pkg")
  pdf_gen_server_dir = abspath("../node")
  web_application_dir = abspath("../web")
  admin_email = "darius.calliet@coveredcreatives.com"
  service_account_email = format("%s@appspot.gserviceaccount.com", data.google_project.new_orleans_connection.project_id)
  default_pw = "eHon#X33S5P9XeLq"
  domain = "thenolaconnect.com"
}

# zip our code so we can store it for deployment
data "archive_file" "primary_server_source" {
  type        = "zip"
  source_dir  = local.primary_server_dir
  output_path = "/tmp/primary_server.zip"
}

data "archive_file" "pdf_gen_server_source" {
  type        = "zip"
  source_dir  = local.pdf_gen_server_dir
  output_path = "/tmp/pdf_gen_server.zip"
}

data "archive_file" "web_application_source" {
  type        = "zip"
  source_dir  = local.web_application_dir
  output_path = "/tmp/web.zip"
}


data "google_project" "new_orleans_connection" {
  project_id = var.project_id
}

data "google_service_account" "new_orleans_connection" {
  account_id = local.service_account_email
}

resource "google_project_iam_member" "admin_cloud_sql_instance_user" {
  project = data.google_project.new_orleans_connection.project_id
  role = "roles/cloudsql.instanceUser"
  member = format("user:%s", local.admin_email)
}

resource "google_project_iam_member" "service_account_cloud_sql_instance_user" {
  project = data.google_project.new_orleans_connection.project_id
  role = "roles/cloudsql.instanceUser"
  member = format("serviceAccount:%s", local.service_account_email)
}


# Lets create a google storage bucket to store our applications functions
resource "google_storage_bucket" "functions" {
    project = data.google_project.new_orleans_connection.project_id
    name = "the_new_orleans_connection_functions"
    location = "US"
}

# add the zipped file to the bucket.
resource "google_storage_bucket_object" "primary_server_zip" {
  # Use an MD5 here. If there's no changes to the source code, this won't change either.
  # We can avoid unnecessary redeployments by validating the code is unchanged, and forcing
  # a redeployment when it has!
  name   = "${data.archive_file.primary_server_source.output_md5}.zip"
  bucket = google_storage_bucket.functions.name
  source = data.archive_file.primary_server_source.output_path
}

resource "google_storage_bucket_object" "pdf_gen_server_zip" {
  # Use an MD5 here. If there's no changes to the source code, this won't change either.
  # We can avoid unnecessary redeployments by validating the code is unchanged, and forcing
  # a redeployment when it has!
  name   = "${data.archive_file.pdf_gen_server_source.output_md5}.zip"
  bucket = google_storage_bucket.functions.name
  source = data.archive_file.pdf_gen_server_source.output_path
}

resource "google_storage_bucket" "web" {
    project = data.google_project.new_orleans_connection.project_id
    name = "the_new_orleans_connection_web"
    location = "US"
    website {
      main_page_suffix = "index.html"
      not_found_page   = "index.html"
    }
}

resource "google_storage_bucket" "web_http" {
    project = data.google_project.new_orleans_connection.project_id
    name = format("www.%s", local.domain)
    location = "US"
    website {
      main_page_suffix = "index.html"
      not_found_page   = "index.html"
    }
}

resource "google_storage_default_object_access_control" "web_http_acl_read" {
  bucket = google_storage_bucket.web_http.name
  role   = "READER"
  entity = "allUsers"
}

resource "google_storage_bucket_object" "web_application_zip" {
  # Use an MD5 here. If there's no changes to the source code, this won't change either.
  # We can avoid unnecessary redeployments by validating the code is unchanged, and forcing
  # a redeployment when it has!
  name   = "${data.archive_file.web_application_source.output_md5}.zip"
  bucket = google_storage_bucket.web.name
  source = data.archive_file.web_application_source.output_path
}

# storage bucket for created QR codes
resource "google_storage_bucket" "qr_codes" {
    project = data.google_project.new_orleans_connection.project_id
    name = "the_new_orleans_connection_qr_codes"
    location = "US"
}

# storage bucket for received files (pdfs, pngs, etc)
data "google_storage_bucket" "company_assets" {
    name = "the_new_orleans_connection_company_assets"
}

# enabled sql admin
resource "google_project_service" "csql" {
  project = data.google_project.new_orleans_connection.project_id
  service = "sqladmin.googleapis.com"
  disable_dependent_services = true
  disable_on_destroy         = false
}

# primary business operations database instance
resource "google_sql_database_instance" "company_database_instance" {
  project = data.google_project.new_orleans_connection.project_id
  name             = "qr-code-instance"
  region           = var.region
  database_version = "POSTGRES_14"
  root_password = local.default_pw
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    database_flags {
      name  = "cloudsql.iam_authentication"
      value = "on"
    }
  }
}

# qr mappings database
resource "google_sql_database" "company_database" {
  project = data.google_project.new_orleans_connection.project_id
  instance = google_sql_database_instance.company_database_instance.name
  name     = "qr_mappings"
}

resource "google_sql_user" "default" {
  project = data.google_project.new_orleans_connection.project_id
  instance = google_sql_database_instance.company_database_instance.name
  name     = "postgres"
  password = local.default_pw
  type = "BUILT_IN"
}

# sql user admin uses iam for login access
resource "google_sql_user" "admin" {
  project = data.google_project.new_orleans_connection.project_id
  instance = google_sql_database_instance.company_database_instance.name
  name     = local.admin_email
  type = "CLOUD_IAM_USER"
}

# sql user service account for cloud function access
resource "google_sql_user" "service_account" {
  project = data.google_project.new_orleans_connection.project_id
  instance = google_sql_database_instance.company_database_instance.name
  name     = trimsuffix(local.service_account_email, "gserviceaccount.com")
  type = "CLOUD_IAM_SERVICE_ACCOUNT"
}

# enable cloudfunctions service
resource "google_project_service" "cf" {
  project = data.google_project.new_orleans_connection.project_id
  service = "cloudfunctions.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# enable cloud build service
resource "google_project_service" "cb" {
  project = data.google_project.new_orleans_connection.project_id
  service = "cloudbuild.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The cloud function resource.
resource "google_cloudfunctions_function" "qr_code_generate" {
  project = data.google_project.new_orleans_connection.project_id
  name                  = "generate"
  description = "receives a name + file via http and returns a QR code rendering that file."
  region                = var.region
  runtime = "go116"
  trigger_http = true
  entry_point = "GenerateQRHandler"  # Set the entry point 
  source_archive_bucket = google_storage_bucket.functions.name
  source_archive_object = google_storage_bucket_object.cloud_functions_zip.name
  service_account_email = local.service_account_email
  environment_variables = {
    DB_USERNAME = "postgres"
    DB_PASSWORD = local.default_pw
    DB_NAME = google_sql_database.company_database.name
    DB_PORT = 5432
    DB_HOSTNAME = google_sql_database_instance.company_database_instance.public_ip_address
    INSTANCE_CONNECTION_NAME = google_sql_database_instance.company_database_instance.connection_name
    UNIX_SOCKET_PATH = format("/cloudsql/%s", google_sql_database_instance.company_database_instance.connection_name)
    RETRIEVE_HTTPS_TRIGGER_URL = google_cloudfunctions_function.qr_code_retrieve.https_trigger_url
  }
}

resource "google_cloudfunctions_function" "qr_code_retrieve" {
  project = data.google_project.new_orleans_connection.project_id
  name                  = "retrieve"
  description = "lists metadata of stored qr codes or returns image for single lookup"
  region                = var.region
  runtime = "go116"
  trigger_http = true
  entry_point = "RetrieveQRHandler"  # Set the entry point 
  source_archive_bucket = google_storage_bucket.functions.name
  source_archive_object = google_storage_bucket_object.cloud_functions_zip.name
  service_account_email = local.service_account_email
  environment_variables = {
    DB_USERNAME = "postgres"
    DB_PASSWORD = local.default_pw
    DB_NAME = google_sql_database.company_database.name
    DB_PORT = 5432
    DB_HOSTNAME = google_sql_database_instance.company_database_instance.public_ip_address
    INSTANCE_CONNECTION_NAME = google_sql_database_instance.company_database_instance.connection_name
    UNIX_SOCKET_PATH = format("/cloudsql/%s", google_sql_database_instance.company_database_instance.connection_name)
  }
}

# IAM Configuration. This allows unauthenticated, public access to the function.
# Change this if you require more control here.
resource "google_cloudfunctions_function_iam_member" "invoker_generate" {
  project = data.google_project.new_orleans_connection.project_id
  cloud_function = google_cloudfunctions_function.qr_code_generate.name

  region         = var.region
  role   = "roles/cloudfunctions.invoker"
  member = "allUsers"
}

resource "google_cloudfunctions_function_iam_member" "invoker_retrieve" {
  project = data.google_project.new_orleans_connection.project_id
  cloud_function = google_cloudfunctions_function.qr_code_retrieve.name

  region         = var.region
  role   = "roles/cloudfunctions.invoker"
  member = "allUsers"
}

// Cloud Build React Web Application

resource "google_project_iam_member" "admin_cloud_build_builds_editor" {
  project = data.google_project.new_orleans_connection.project_id
  role = "roles/cloudbuild.builds.editor"
  member = format("user:%s", local.admin_email)
}

# enable source repo service
resource "google_project_service" "sr" {
  project = data.google_project.new_orleans_connection.project_id
  service = "sourcerepo.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

resource "google_sourcerepo_repository" "qr_code" {
  project = data.google_project.new_orleans_connection.project_id
  name = "qr-code"
}

resource "google_sourcerepo_repository_iam_member" "member" {
  project = data.google_project.new_orleans_connection.project_id
  repository = google_sourcerepo_repository.qr_code.name
  role = "roles/viewer"
  member = format("user:%s", local.admin_email)
}

resource "google_cloudbuild_trigger" "github_push_trigger" {
  location = "us-central1"
  source_to_build {
    uri = google_sourcerepo_repository.qr_code.url
    ref = "refs/heads/main"
    repo_type = "CLOUD_SOURCE_REPOSITORIES"
  }

  build {
    source {
      storage_source {
        bucket = google_storage_bucket.web.name
        object = google_storage_bucket_object.web_application_zip.name
      }
    }

    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", google_storage_bucket_object.web_application_zip.self_link, google_storage_bucket_object.web_application_zip.name]
      timeout = "120s"
      secret_env = ["MY_SECRET"]
    }
    step {
      name = "node:16.13.0"
      entrypoint = "npm"
      args = ["install"]
    }
    step {
      name = "node:16.13.0"
      entrypoint = "npm"
      args = ["run", "build"]
    }
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["-m", "cp", "-r", "build/*", "gs://${format("www.%s", local.domain)}"]
    }
  }
}

# enabled google sheets api
resource "google_project_service" "sheets" {
  project = data.google_project.new_orleans_connection.project_id
  service = "sheets.googleapis.com"
  disable_dependent_services = true
  disable_on_destroy         = false
}

# enable google cloud run service


# google cloud run, html to pdf service
