resource "google_compute_instance" "application_server" {
  name         = "application-server"
  machine_type = "e2-small"
  zone         = "${var.region}-a"
  tags         = ["web", "http-server", "https-server"]
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }
  network_interface {
    subnetwork = google_compute_subnetwork.subnet.self_link
    network_ip = google_compute_address.private_ip_address.address
  }
  service_account {
    email = local.service_account_email
    scopes = [
      "logging-write",
      "storage-rw",
      "cloud-platform"
    ]
  }
}

resource "google_compute_instance_group" "application" {
  name        = "application-instance-group"
  description = "VMs related to serving application traffic"
  zone        = "${var.region}-a"
  network     = google_compute_network.vpc.id

  instances = [
    google_compute_instance.application_server.self_link,
  ]

  named_port {
    name = "http"
    port = "8080"
  }

  named_port {
    name = "https"
    port = "8443"
  }
}
