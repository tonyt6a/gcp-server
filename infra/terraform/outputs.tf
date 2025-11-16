output "image_path" {
  value = "${var.region}-docker.pkg.dev/${var.project_id}/${var.artifact_repo}/rt-gaas:${var.image_tag}"
}