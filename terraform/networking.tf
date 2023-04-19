// Let's create a vpc for application resources to connect within
resource "google_compute_network" "vpc" {
  name                    = "app-vpc"
  auto_create_subnetworks = false
}

// Let's create an external IP address for global connectivity
resource "google_compute_global_address" "external_ip" {
  name         = "app-lb-ip"
  provider     = google
  address_type = "EXTERNAL"
  ip_version   = "IPV4"
  project      = data.google_project.new_orleans_connection.project_id
  description  = "External IPV4 address"
}

// Then create a private ip address for connecting resources within vpc
resource "google_compute_global_address" "private_ip_address" {
  name          = "private-ip-address"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.vpc.id
}

resource "google_compute_address" "private_ip_address" {
  name         = "private-ip-address"
  address_type = "INTERNAL"
  subnetwork   = google_compute_subnetwork.subnet.self_link
  address      = "10.8.0.2"
}

// Create a subnet for connecting resources within vpc
resource "google_compute_subnetwork" "subnet" {
  name                     = "app-subnet"
  ip_cidr_range            = "10.8.0.0/28"
  region                   = var.region
  network                  = google_compute_network.vpc.id
  private_ip_google_access = true

  secondary_ip_range {
    range_name    = "pod"
    ip_cidr_range = "10.0.16.0/20"
  }

  secondary_ip_range {
    range_name    = "svc"
    ip_cidr_range = "10.0.32.0/20"
  }
}
// TODO
resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = google_compute_network.vpc.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]
}

//  TODO
resource "google_compute_project_default_network_tier" "default" {
  project      = data.google_project.new_orleans_connection.project_id
  network_tier = "PREMIUM"
}

# Add the bucket as a CDN backend
resource "google_compute_backend_bucket" "static_site" {
  provider    = google
  name        = "qr-code-react-storage-bucket"
  description = "Connect source code within bucket to CDN to serve"
  bucket_name = google_storage_bucket.web_http.name
  enable_cdn  = true
}

# GCP forwarding rule HTTPS
resource "google_compute_global_forwarding_rule" "web_https" {
  provider              = google
  name                  = "web-forwarding-rule-https"
  load_balancing_scheme = "EXTERNAL"
  ip_address            = google_compute_global_address.external_ip.address
  ip_protocol           = "TCP"
  port_range            = "443"
  target                = google_compute_target_https_proxy.default_02.id
}

resource "google_compute_global_forwarding_rule" "web_http_8080" {
  provider              = google
  name                  = "web-forwarding-rule-http-8080"
  load_balancing_scheme = "EXTERNAL"
  ip_address            = google_compute_global_address.external_ip.address
  ip_protocol           = "TCP"
  target                = google_compute_target_http_proxy.default.id
  port_range            = "8080"
}


resource "google_compute_target_http_proxy" "default" {
  name        = "api-target-proxy-http"
  description = "HTTP target proxy"
  url_map     = google_compute_url_map.https_redirect.id
}

resource "google_compute_target_https_proxy" "default_02" {
  name    = "default-target-proxy-https-02"
  url_map = google_compute_url_map.default.id
  depends_on = [
    google_compute_managed_ssl_certificate.wild_card_2
  ]
  ssl_certificates = [
    google_compute_managed_ssl_certificate.wild_card_2.id
  ]
}

resource "google_compute_health_check" "http" {
  name = "health-check"
  http_health_check {
    request_path = "/_ah/warmup"
    port_name    = "http"
  }
}

resource "google_compute_router" "router" {
  project = var.project_id
  name    = "nat-router"
  network = google_compute_network.vpc.name
  region  = var.region
}

module "cloud-nat" {
  source                             = "terraform-google-modules/cloud-nat/google"
  version                            = "~> 2.0"
  project_id                         = var.project_id
  region                             = var.region
  router                             = google_compute_router.router.name
  name                               = "nat-config"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"
}
