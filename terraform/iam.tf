resource "google_project_iam_member" "admin_cloud_sql_instance_user" {
  project = data.google_project.new_orleans_connection.project_id
  role    = "roles/cloudsql.instanceUser"
  member  = format("user:%s", var.admin_email)
}

resource "google_project_iam_member" "service_account_cloud_sql_instance_user" {
  project = data.google_project.new_orleans_connection.project_id
  role    = "roles/cloudsql.instanceUser"
  member  = format("serviceAccount:%s", local.service_account_email)
}

resource "google_project_iam_member" "admin_cloud_build_builds_editor" {
  project = data.google_project.new_orleans_connection.project_id
  role    = "roles/cloudbuild.builds.editor"
  member  = format("user:%s", var.admin_email)
}
