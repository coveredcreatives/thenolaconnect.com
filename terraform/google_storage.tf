# Lets create a google storage bucket to store our applications functions
resource "google_storage_bucket" "functions" {
  project  = data.google_project.new_orleans_connection.project_id
  name     = "the_new_orleans_connection_functions"
  location = "US"
}

# Populate the bucket with an object for our zipped application code
resource "google_storage_bucket_object" "primary_server_zip" {
  name   = "${data.archive_file.primary_server_source.output_md5}.zip"
  bucket = google_storage_bucket.functions.name
  source = data.archive_file.primary_server_source.output_path
}

# Lets create a google storage bucket to store our applications execution
resource "google_storage_bucket" "executables" {
  project  = data.google_project.new_orleans_connection.project_id
  name     = "the_new_orleans_connection_executables"
  location = "US"
}

# Lets create a google storage bucket to store our web application build
resource "google_storage_bucket" "web" {
  project  = data.google_project.new_orleans_connection.project_id
  name     = "the_new_orleans_connection_web"
  location = "US"
  website {
    main_page_suffix = "index.html"
    not_found_page   = "index.html"
  }
}

# Lets create a google storage bucket to serve under our CDN
resource "google_storage_bucket" "web_http" {
  project  = data.google_project.new_orleans_connection.project_id
  name     = "www.thenolaconnect.com"
  location = "US"
  website {
    main_page_suffix = "index.html"
    not_found_page   = "index.html"
  }
}

# Lets create a google storage bucket to serve under our CDN
resource "google_storage_bucket" "web_http_alt" {
  project  = data.google_project.new_orleans_connection.project_id
  name     = "www.theneworleansseafoodconnection.com"
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

resource "google_storage_default_object_access_control" "web_http_alt_acl_read" {
  bucket = google_storage_bucket.web_http_alt.name
  role   = "READER"
  entity = "allUsers"
}

# Lets store our web application source code in a bucket.
resource "google_storage_bucket_object" "web_application_zip" {
  # Use an MD5 here. If there's no changes to the source code, this won't change either.
  # We can avoid unnecessary redeployments by validating the code is unchanged, and forcing
  # a redeployment when it has!
  name   = "${data.archive_file.web_application_source.output_md5}.zip"
  bucket = google_storage_bucket.web.name
  source = data.archive_file.web_application_source.output_path
}

# Lets create a storage bucket for created QR codes
resource "google_storage_bucket" "qr_codes" {
  project  = data.google_project.new_orleans_connection.project_id
  name     = "the_new_orleans_connection_qr_codes"
  location = "US"
}

# Lets create a storage bucket for received files (pdfs, pngs, etc)
data "google_storage_bucket" "company_assets" {
  name = "the_new_orleans_connection_company_assets"
}
