resource "google_compute_managed_ssl_certificate" "wild_card_2" {
  provider = google-beta
  name     = "alb-wildcard-ssl-cert-2"
  type     = "MANAGED"
  project  = data.google_project.new_orleans_connection.project_id
  managed {
    domains = [
      "api.thenolaconnect.com",
      "thenolaconnect.com",
      "www.theneworleansseafoodconnection.com",
    ]
  }
}
