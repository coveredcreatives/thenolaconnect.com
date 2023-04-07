resource "google_compute_url_map" "default" {
  provider        = google
  name            = "web-url-map"
  default_service = google_compute_backend_bucket.static_site.id

  host_rule {
    hosts = [
      "app.thenolaconnect.com",
      "theneworleansseafoodconnection.com",
      "www.theneworleansseafoodconnection.com"
    ]
    path_matcher = "web"
  }

  path_matcher {
    name            = "web"
    default_service = google_compute_backend_bucket.static_site.id

    path_rule {
      paths   = ["/"]
      service = google_compute_backend_bucket.static_site.id
    }
  }

  host_rule {
    hosts = [
      "spec.thenolaconnect.com",
    ]
    path_matcher = "openapi"
  }

  host_rule {
    hosts = [
      "api.thenolaconnect.com"
    ]
    path_matcher = "api"
  }

  path_matcher {
    name            = "openapi"
    default_service = google_compute_backend_service.default.self_link

    path_rule {
      paths   = ["/"]
      service = google_compute_backend_service.default.self_link
    }
  }

  path_matcher {
    name            = "api"
    default_service = google_compute_backend_service.default.self_link

    path_rule {
      paths   = ["/"]
      service = google_compute_backend_service.default.self_link
    }
  }
}
