output "google_compute_global_address__external_ip_address" {
  value       = google_compute_global_address.external_ip.address
  description = "External static IP address for React app"
}

output "google_sql_database_instance__company_database_instance" {
  value       = google_sql_database_instance.company_database_instance.first_ip_address
  description = "First IP address for associated company database instance"
}
