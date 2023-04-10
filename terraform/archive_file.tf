data "archive_file" "primary_server_source" {
  type        = "zip"
  source_dir  = abspath("../pkg")
  output_path = "/tmp/primary_server.zip"
}

# data "archive_file" "pdf_gen_server_source" {
#   type        = "zip"
#   source_dir  = local.pdf_gen_server_dir
#   output_path = "/tmp/pdf_gen_server.zip"
# }

data "archive_file" "web_application_source" {
  type        = "zip"
  source_dir  = local.web_application_dir
  output_path = "/tmp/web.zip"
}
