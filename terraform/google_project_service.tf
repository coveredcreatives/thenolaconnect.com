resource "google_project_service" "csql" {
  project                    = data.google_project.new_orleans_connection.project_id
  service                    = "sqladmin.googleapis.com"
  disable_dependent_services = true
  disable_on_destroy         = false
}

resource "google_project_service" "cf" {
  project = data.google_project.new_orleans_connection.project_id
  service = "cloudfunctions.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

resource "google_project_service" "cb" {
  project = data.google_project.new_orleans_connection.project_id
  service = "cloudbuild.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

resource "google_project_service" "sr" {
  project = data.google_project.new_orleans_connection.project_id
  service = "sourcerepo.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

resource "google_project_service" "sheets" {
  project                    = data.google_project.new_orleans_connection.project_id
  service                    = "sheets.googleapis.com"
  disable_dependent_services = true
  disable_on_destroy         = false
}

resource "google_project_service" "artifactregistry" {
  project                    = data.google_project.new_orleans_connection.project_id
  service                    = "artifactregistry.googleapis.com"
  disable_dependent_services = true
  disable_on_destroy         = false
}

resource "google_project_service" "container" {
  project                    = data.google_project.new_orleans_connection.project_id
  service                    = "container.googleapis.com"
  disable_dependent_services = true
  disable_on_destroy         = false
}

resource "google_project_service" "cm" {
  project = data.google_project.new_orleans_connection.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

resource "google_project_service" "sn" {
  project = data.google_project.new_orleans_connection.project_id
  service = "servicenetworking.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}
