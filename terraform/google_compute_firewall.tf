resource "google_compute_firewall" "http" {
  name        = "${google_compute_network.vpc.name}-firewall-http"
  network     = google_compute_network.vpc.self_link
  description = "Allow web applications traffic"

  direction = "INGRESS"
  allow {
    protocol = "tcp"
    ports = [
      "80", # HTTP
      "8080"
    ]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["web"]
}

resource "google_compute_firewall" "https" {
  name        = "${google_compute_network.vpc.name}-firewall-https"
  network     = google_compute_network.vpc.self_link
  description = "Allow web applications traffic"

  direction = "INGRESS"
  allow {
    protocol = "tcp"
    ports = [
      "443", # HTTPS
    ]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["web"]
}

resource "google_compute_firewall" "ssh" {
  name        = "${google_compute_network.vpc.name}-firewall-ssh"
  network     = google_compute_network.vpc.self_link
  description = "Allow web applications traffic"

  direction = "INGRESS"
  allow {
    protocol = "tcp"
    ports = [
      "22", // SSH
    ]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["web"]
}

resource "google_compute_firewall" "postgresql" {
  name    = "${google_compute_network.vpc.name}-firewall-postgresql"
  network = google_compute_network.vpc.name

  direction = "EGRESS"
  allow {
    protocol = "tcp"
    ports    = ["5432"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["web"]
}
