# Create External Static IP Address
resource "google_compute_global_address" "external_ip" {
  name         = "app-lb-ip"
  provider     = google
  address_type = "EXTERNAL"
  ip_version   = "IPV4"
  project      = data.google_project.new_orleans_connection.project_id
  description  = "External IPV4 address"
}

# Output network value
output "external_ip" {
  value       = google_compute_global_address.external_ip.address
  description = "External static IP address for React app"
}

resource "google_compute_network" "vpc" {
  name                    = "app-vpc"
  auto_create_subnetworks = false
}

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

resource "google_compute_global_address" "private_ip_address" {
  name          = "private-ip-address"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.vpc.id
}

resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = google_compute_network.vpc.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]
}

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

# Create HTTPS certificate
resource "google_compute_managed_ssl_certificate" "web" {
  provider = google-beta
  name     = "ssl-certificate"
  project  = data.google_project.new_orleans_connection.project_id
  managed {
    domains = [var.domain, format("www.%s", var.domain)]
  }
}

# GCP URL MAP HTTPS
resource "google_compute_url_map" "web_https" {
  provider        = google
  name            = "web-url-map-https"
  default_service = google_compute_backend_bucket.static_site.id

  host_rule {
    hosts        = ["thenolaconnect.com"]
    path_matcher = "qrmapping"
  }

  path_matcher {
    name            = "qrmapping"
    default_service = google_compute_backend_bucket.static_site.id

    path_rule {
      paths   = ["/"]
      service = google_compute_backend_bucket.static_site.id
    }
  }
}

# GCP target proxy HTTPS
resource "google_compute_target_https_proxy" "web_https" {
  provider         = google
  name             = "web-target-proxy-https"
  url_map          = google_compute_url_map.web_https.self_link
  ssl_certificates = [google_compute_managed_ssl_certificate.web.self_link]
}

# GCP forwarding rule HTTPS
resource "google_compute_global_forwarding_rule" "web_https" {
  provider              = google
  name                  = "web-forwarding-rule-https"
  load_balancing_scheme = "EXTERNAL"
  ip_address            = google_compute_global_address.external_ip.address
  ip_protocol           = "TCP"
  port_range            = "443"
  target                = google_compute_target_https_proxy.web_https.self_link
}

resource "google_compute_url_map" "web_http" {
  name        = "web-url-map-http"
  description = "Web HTTP load balancer"

  default_url_redirect {
    https_redirect = true
    strip_query    = true
  }
}

resource "google_compute_global_forwarding_rule" "web_http" {
  name                  = "web-forwarding-rule-http"
  load_balancing_scheme = "EXTERNAL"
  ip_address            = google_compute_global_address.external_ip.address
  ip_protocol           = "TCP"
  target                = google_compute_target_http_proxy.web_http.id
  port_range            = "80"
}

resource "google_compute_target_http_proxy" "web_http" {
  name        = "web-target-proxy-http"
  description = "HTTP target proxy"
  url_map     = google_compute_url_map.web_http.id
}
