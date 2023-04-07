resource "google_compute_backend_service" "default" {
  provider    = google-beta
  name        = "application-backend-service"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks         = [google_compute_health_check.default.id]
  load_balancing_scheme = "EXTERNAL"

  backend {
    group = google_compute_instance_group.application.id
  }
}
