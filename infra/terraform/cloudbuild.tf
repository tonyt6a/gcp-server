# # Bucket for Cloud Build logs
# resource "google_storage_bucket" "build_logs" {
#   name                        = "${var.project_id}-cloudbuild-logs"
#   location                    = var.region
#   force_destroy               = true
#   uniform_bucket_level_access = true
# }

# # Cloud Source Repository that will mirror your gcp-server code
# resource "google_cloud_source_repository" "gcp_server_source" {
#   name = "rt-gaas-source"
# }

# # Cloud Build that builds and pushes the image
# resource "google_cloudbuild_build" "gcp_server_build" {
#   timeout    = "1200s"
#   logs_bucket = google_storage_bucket.build_logs.name

#   step {
#     name = "gcr.io/cloud-builders/docker"
#     args = [
#       "build",
#       "-t",
#       "${var.region}-docker.pkg.dev/${var.project_id}/${var.artifact_repo}/rt-gaas:${var.image_tag}",
#       ".",
#     ]
#   }

#   step {
#     name = "gcr.io/cloud-builders/docker"
#     args = [
#       "push",
#       "${var.region}-docker.pkg.dev/${var.project_id}/${var.artifact_repo}/rt-gaas:${var.image_tag}",
#     ]
#   }

#   images = [
#     "${var.region}-docker.pkg.dev/${var.project_id}/${var.artifact_repo}/rt-gaas:${var.image_tag}",
#   ]

#   source {
#     repo_source {
#       project_id  = var.project_id
#       repo_name   = google_cloud_source_repository.gcp_server_source.name
#       branch_name = "main"
#     }
#   }
# }