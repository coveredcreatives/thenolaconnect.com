resource "google_compute_managed_ssl_certificate" "web" {
  provider = google-beta
  name     = "ssl-certificate"
  type     = "MANAGED"
  project  = data.google_project.new_orleans_connection.project_id
  managed {
    domains = [var.domain, format("www.%s", var.domain)]
  }
}

resource "google_compute_managed_ssl_certificate" "wild_card" {
  provider = google-beta
  name     = "alb-wildcard-ssl-cert-1"
  type     = "MANAGED"
  project  = data.google_project.new_orleans_connection.project_id
  managed {
    domains = [
      "app.thenolaconnect.com",
      "api.thenolaconnect.com",
      "spec.thenolaconnect.com",
      "theneworleansseafoodconnection.com",
      "www.theneworleansseafoodconnection.com",
    ]
  }
}
