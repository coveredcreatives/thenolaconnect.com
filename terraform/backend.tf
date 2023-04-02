terraform {
  backend "gcs" {
    bucket = "the-new-orleans-connection-terraform-state" # Bucket is passed in via cli arg. Eg, terraform init -reconfigure -backend-configuration=dev.tfbackend
  }
}
