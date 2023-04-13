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
sudo rm /var/lib/man-db/auto-update
sudo apt-get update
sudo apt-get install wget
# Install wkhtmltopdf
wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6.1-2/wkhtmltox_0.12.6.1-2.bullseye_amd64.deb
sudo apt-get -y install ./wkhtmltox_0.12.6.1-2.bullseye_amd64.deb
# Install app executable
sudo gsutil cp gs://${google_storage_bucket.executables.name}/cli-linux-amd64@latest /bin/thenolaconnect
sudo chmod +x /bin/thenolaconnect
# Create config file
echo "ENV: production
DB_USERNAME: ${var.db_username}
DB_PASSWORD: ${var.db_password}
DB_NAME: ${var.db_name}
DB_PORT: ${var.db_port}
DB_HOSTNAME: ${google_sql_database_instance.company_database_instance.first_ip_address}
HTTP_PORT: ${var.http_port}
DNS_PRINTER_IPV4_ADDRESS: ${var.dns_printer_ipv4_address}
DNS_RETRIEVE_TRIGGER_URL: ${google_cloudfunctions_function.qr_code_retrieve.https_trigger_url}
GOOGLE_APPLICATION_SERVICE_ACCOUNT_EMAIL: ${var.google_application_service_account_email}
REACT_APP_GOOGLE_API_KEY: ${var.google_api_key}
GOOGLE_STORAGE_BUCKET_NAME: ${data.google_storage_bucket.company_assets.name}
GOOGLE_API_KEY_ORDERS: ${var.google_api_key}
GOOGLE_FORM_ID_ORDERS: ${var.google_form_id_orders}
TWILIO_ACCOUNT_SID: ${var.twilio_account_sid}
TWILIO_ACCOUNT_AUTH_TOKEN: ${var.twilio_account_auth_token}
TWILIO_CONVERSATION_SERVICE_SID: ${var.twilio_conversation_service_sid}
APP_CONFIG_PATH: /home/devscrum" > /home/devscrum/production.yaml
# Create sysctl service
echo "
[Unit]
Description=TheNolaConnect Go Service
ConditionPathExists=/bin/thenolaconnect
After=network.target
[Service]
Type=simple
User=devscrum
Group=devscrum
WorkingDirectory=/home/devscrum
ExecStart=/bin/thenolaconnect
Restart=on-failure
RestartSec=10
SyslogIdentifier=thenolaconnect
Environment=NOLA_APP_CONFIG_PATH=/homes/devscrum
Environment=NOLA_ENV=production
[Install]
WantedBy=multi-user.target
" | sudo tee /etc/systemd/system/thenolaconnect.service
sudo systemctl daemon-reload
sudo service thenolaconnect start
sudo service thenolaconnect status
sudo systemctl enable thenolaconnect
sudo systemctl start thenolaconnect
# enable logging
mkdir /var/log/thenolaconnect
chown devscrum:devscrum /var/log/thenolaconnect
echo "
if $programname == 'thenolaconnect' then /var/log/thenolaconnect/output.log & stop
" | sudo tee /etc/rsyslog.d/thenolaconnect.conf
systemctl restart rsyslog.service
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
