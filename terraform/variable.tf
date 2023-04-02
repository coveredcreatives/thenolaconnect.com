variable "config_bucket" {
  description = "The bucket to store the terraform configration data"
}

variable "organization_id" {
  description = "The organization ID in Google Cloud to use for these resources."
}

variable "project_id" {
  description = "The project ID in Google Cloud to use for these resources."
}

variable "region" {
  description = "The region in Google Cloud where the resources will be deployed."
}

variable "database_password" {
  description = "Password to connect to database instance"
  sensitive   = true
}

variable "domain" {
  description = "top level main domain of project"
}

variable "admin_email" {
  description = "The email for creating iam roles"
}


