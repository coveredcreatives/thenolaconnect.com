locals {
  primary_server_dir    = abspath("../pkg")
  pdf_gen_server_dir    = abspath("../node")
  web_application_dir   = abspath("../web")
  startup_scripts_dir   = abspath("../scripts")
  service_account_email = format("%s@appspot.gserviceaccount.com", var.project_id)
}
