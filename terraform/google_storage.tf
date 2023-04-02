# zip our code so we can store it for deployment
data "archive_file" "primary_server_source" {
  type        = "zip"
  source_dir  = local.primary_server_dir
  output_path = "/tmp/primary_server.zip"
}

# data "archive_file" "pdf_gen_server_source" {
#   type        = "zip"
#   source_dir  = local.pdf_gen_server_dir
#   output_path = "/tmp/pdf_gen_server.zip"
# }

data "archive_file" "web_application_source" {
  type        = "zip"
  source_dir  = local.web_application_dir
  output_path = "/tmp/web.zip"
}

# Lets create a google storage bucket to store our applications functions
resource "google_storage_bucket" "functions" {
  project  = data.google_project.new_orleans_connection.project_id
  name     = "the_new_orleans_connection_functions"
  location = "US"
}

# add the zipped file to the bucket.
# resource "google_storage_bucket_object" "cloud_functions_zip" {
#   # Use an MD5 here. If there's no changes to the source code, this won't change either.
#   # We can avoid unnecessary redeployments by validating the code is unchanged, and forcing
#   # a redeployment when it has!
#   name   = "${data.archive_file.cloud_functions_source.output_md5}.zip"
#   bucket = google_storage_bucket.functions.name
#   source = data.archive_file.cloud_functions.output_path
# }

resource "google_storage_bucket_object" "primary_server_zip" {
  # Use an MD5 here. If there's no changes to the source code, this won't change either.
  # We can avoid unnecessary redeployments by validating the code is unchanged, and forcing
  # a redeployment when it has!
  name   = "${data.archive_file.primary_server_source.output_md5}.zip"
  bucket = google_storage_bucket.functions.name
  source = data.archive_file.primary_server_source.output_path
}

# resource "google_storage_bucket_object" "pdf_gen_server_zip" {
#   # Use an MD5 here. If there's no changes to the source code, this won't change either.
#   # We can avoid unnecessary redeployments by validating the code is unchanged, and forcing
#   # a redeployment when it has!
#   name   = "${data.archive_file.pdf_gen_server_source.output_md5}.zip"
#   bucket = google_storage_bucket.functions.name
#   source = data.archive_file.pdf_gen_server_source.output_path
# }

resource "google_storage_bucket" "web" {
  project  = data.google_project.new_orleans_connection.project_id
  name     = "the_new_orleans_connection_web"
  location = "US"
  website {
    main_page_suffix = "index.html"
    not_found_page   = "index.html"
  }
}

resource "google_storage_bucket" "web_http" {
  project  = data.google_project.new_orleans_connection.project_id
  name     = format("www.%s", var.domain)
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
  project  = data.google_project.new_orleans_connection.project_id
  name     = "the_new_orleans_connection_qr_codes"
  location = "US"
}

# storage bucket for received files (pdfs, pngs, etc)
data "google_storage_bucket" "company_assets" {
  name = "the_new_orleans_connection_company_assets"
}
