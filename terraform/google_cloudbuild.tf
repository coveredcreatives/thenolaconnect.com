resource "google_cloudbuild_trigger" "build_website" {
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
      name    = "gcr.io/cloud-builders/gsutil"
      args    = ["cp", google_storage_bucket_object.web_application_zip.self_link, google_storage_bucket_object.web_application_zip.name]
      timeout = "120s"
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
      args = ["-m", "cp", "-r", "build/*", "gs://www.thenolaconnect.com"]
    }
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["-m", "cp", "-r", "build/*", "gs://www.theneworleansseafoodconnection.com"]
    }
  }
}

resource "google_cloudbuild_trigger" "build_executable" {
  location = "us-central1"
  source_to_build {
    uri       = google_sourcerepo_repository.qr_code.url
    ref       = "refs/heads/main"
    repo_type = "CLOUD_SOURCE_REPOSITORIES"
  }

  build {
    source {
      storage_source {
        bucket = google_storage_bucket.functions.name
        object = google_storage_bucket_object.primary_server_zip.name
      }
    }
    step {
      name    = "gcr.io/cloud-builders/gsutil"
      args    = ["cp", google_storage_bucket_object.primary_server_zip.self_link, google_storage_bucket_object.primary_server_zip.name]
      timeout = "120s"
    }
    step {
      name       = "golang:1.20"
      entrypoint = "go"
      args       = ["build", "./cmd/cli"]
      timeout    = "120s"
    }
    step {
      name       = "ubuntu"
      entrypoint = "/bin/sh"
      args       = ["build", "./cmd/cli"]
      timeout    = "120s"
    }
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "cli", "cli@$SHORT_SHA"]
    }
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "cli", "cli@latest"]
    }
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "cli@latest", "gs://${google_storage_bucket.executables.name}"]
    }
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "cli@$SHORT_SHA", "gs://${google_storage_bucket.executables.name}"]
    }
  }
}
