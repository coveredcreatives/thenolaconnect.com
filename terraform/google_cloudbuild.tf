resource "google_cloudbuild_trigger" "build_website" {
  location = "us-central1"
  source_to_build {
    uri       = data.google_sourcerepo_repository.thenolaconnect.url
    ref       = "refs/heads/dev"
    repo_type = "GITHUB"
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
    uri       = data.google_sourcerepo_repository.thenolaconnect.url
    ref       = "refs/heads/dev"
    repo_type = "GITHUB"
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
      args       = ["build", "-o", "cli-linux-amd64", "./cmd/cli"]
      timeout    = "120s"
      env = [
        "GOOS=linux",
        "GOARCH=amd64"
      ]
    }
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "cli-linux-amd64", "cli-linux-amd64@$SHORT_SHA"]
    }
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "cli-linux-amd64", "cli-linux-amd64@latest"]
    }
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "cli-linux-amd64@latest", "gs://${google_storage_bucket.executables.name}"]
    }
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "cli-linux-amd64@$SHORT_SHA", "gs://${google_storage_bucket.executables.name}"]
    }
  }
}
