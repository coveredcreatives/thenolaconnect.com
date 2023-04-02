data "google_project" "new_orleans_connection" {
  project_id = var.project_id
}

data "google_service_account" "new_orleans_connection" {
  account_id = format("%s@appspot.gserviceaccount.com", var.project_id)
}
