resource "google_compute_instance" "application_server" {
  name         = "application-server"
  machine_type = "f1_micro"
  zone         = format("%s-a", var.region)
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }
  network_interface {
    network = google_compute_subnetwork.subnet.self_link
  }
  service_account {
    email = local.service_account_email
    scopes = [
      "default",
      "logging-write",
      "storage-rw",
      "cloud-platform"
    ]
  }
}
