resource "google_cloudbuild_trigger" "github_push_trigger" {
  location = "us-central1"
  source_to_build {
    uri       = google_sourcerepo_repository.qr_code.url
    ref       = "refs/heads/main"
    repo_type = "CLOUD_SOURCE_REPOSITORIES"
  }

  build {
    source {
      storage_source {
        bucket = google_storage_bucket.web.name
        object = google_storage_bucket_object.web_application_zip.name
      }
    }

    step {
      name       = "gcr.io/cloud-builders/gsutil"
      args       = ["cp", google_storage_bucket_object.web_application_zip.self_link, google_storage_bucket_object.web_application_zip.name]
      timeout    = "120s"
      secret_env = ["MY_SECRET"]
    }
    step {
      name       = "node:16.13.0"
      entrypoint = "npm"
      args       = ["install"]
    }
    step {
      name       = "node:16.13.0"
      entrypoint = "npm"
      args       = ["run", "build"]
    }
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["-m", "cp", "-r", "build/*", "gs://${format("www.%s", var.domain)}"]
    }
  }
}
