resource "google_sourcerepo_repository" "qr_code" {
  project = data.google_project.new_orleans_connection.project_id
  name    = "qr-code"
}

resource "google_sourcerepo_repository_iam_member" "member" {
  project    = data.google_project.new_orleans_connection.project_id
  repository = google_sourcerepo_repository.qr_code.name
  role       = "roles/viewer"
  member     = format("user:%s", var.admin_email)
}
