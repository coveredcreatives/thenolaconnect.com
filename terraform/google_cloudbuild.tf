resource "google_cloudbuild_trigger" "build_website" {
  location = "us-central1"
  name     = "build-website"

  github {
    owner = var.github_username
    name  = "thenolaconnect.com"
    push {
      branch = "dev"
    }
  }

  build {
    source {
      repo_source {
        project_id  = var.project_id
        repo_name   = data.google_sourcerepo_repository.thenolaconnect.name
        dir         = "./web"
        branch_name = "dev"
      }
    }
    step {
      name       = "node:16.13.0"
      dir        = "web"
      entrypoint = "npm"
      args       = ["install"]
    }
    step {
      name       = "node:16.13.0"
      dir        = "web"
      entrypoint = "npm"
      args       = ["run", "build"]
    }
    step {
      name = "gcr.io/cloud-builders/gsutil"
      dir  = "web"
      args = ["-m", "cp", "-r", "build/*", "gs://www.thenolaconnect.com"]
    }
    step {
      name = "gcr.io/cloud-builders/gsutil"
      dir  = "web"
      args = ["-m", "cp", "-r", "build/*", "gs://www.theneworleansseafoodconnection.com"]
    }
  }
}

resource "google_cloudbuild_trigger" "build_executable" {
  location = "us-central1"
  name     = "build-executable"

  github {
    owner = var.github_username
    name  = "thenolaconnect.com"
    push {
      branch = "dev"
    }
  }

  build {
    source {
      repo_source {
        project_id  = var.project_id
        repo_name   = data.google_sourcerepo_repository.thenolaconnect.name
        dir         = "./pkg"
        branch_name = "dev"
      }
    }
    step {
      name       = "golang:1.20"
      dir        = "pkg"
      entrypoint = "go"
      args       = ["build", "-o", "cli-linux-amd64", "./cmd/app"]
      timeout    = "120s"
      env = [
        "GOOS=linux",
        "GOARCH=amd64"
      ]
    }
    step {
      name = "gcr.io/cloud-builders/gsutil"
      dir  = "pkg"
      args = ["cp", "cli-linux-amd64", "cli-linux-amd64@$SHORT_SHA"]
    }
    step {
      name = "gcr.io/cloud-builders/gsutil"
      dir  = "pkg"
      args = ["cp", "cli-linux-amd64", "cli-linux-amd64@latest"]
    }
    step {
      name = "gcr.io/cloud-builders/gsutil"
      dir  = "pkg"
      args = ["cp", "cli-linux-amd64@latest", "gs://${google_storage_bucket.executables.name}"]
    }
    step {
      name = "gcr.io/cloud-builders/gsutil"
      dir  = "pkg"
      args = ["cp", "cli-linux-amd64@$SHORT_SHA", "gs://${google_storage_bucket.executables.name}"]
    }
  }
}
