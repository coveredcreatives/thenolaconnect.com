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
  metadata_startup_script = <<SCRIPT
#! /bin/bash
apt update
apt -y install apache2
sudo apt install wget
#! Install wkhtmltopdf
wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-1/wkhtmltox_0.12.6-1.focal_amd64.deb
sudo apt install ./wkhtmltox_0.12.6-1.focal_amd64.deb
#! Install app executable
gsutil cp gs://${google_storage_bucket.executables.name}/cli-linux-amd64@latest /bin/thenolaconnect
#! Update environment
export NOLA_ENV: production
export NOLA_DB_USERNAME: ${var.db_username}
export NOLA_DB_PASSWORD: ${var.db_password}
export NOLA_DB_NAME:${var.db_name}
export NOLA_DB_PORT: ${var.db_port}
export NOLA_DB_HOSTNAME: ${google_sql_database_instance.company_database_instance.first_ip_address}
export NOLA_HTTP_PORT: ${var.http_port}
export NOLA_DNS_PRINTER_IPV4_ADDRESS: ${var.dns_printer_ipv4_address}
export NOLA_DNS_RETRIEVE_TRIGGER_URL: ${google_cloudfunctions_function.qr_code_retrieve.https_trigger_url}
export NOLA_GOOGLE_APPLICATION_SERVICE_ACCOUNT_EMAIL: ${var.google_application_service_account_email}
export NOLA_REACT_APP_GOOGLE_API_KEY: ${var.google_api_key}
export NOLA_GOOGLE_STORAGE_BUCKET_NAME: ${data.google_storage_bucket.company_assets.name}
export NOLA_GOOGLE_API_KEY_ORDERS: ${var.google_api_key}
export NOLA_GOOGLE_FORM_ID_ORDERS: ${var.google_form_id_orders}
export NOLA_TWILIO_ACCOUNT_SID: ${var.twilio_account_sid}
export NOLA_TWILIO_ACCOUNT_AUTH_TOKEN: ${var.twilio_account_auth_token}
export NOLA_TWILIO_CONVERSATION_SERVICE_SID: ${var.twilio_conversation_service_sid}
SCRIPT
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
