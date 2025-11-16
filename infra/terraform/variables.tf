variable "project_id" {
  type    = string
  default = "gcp-server-474714"
}

variable "region" {
  type    = string
  default = "us-central1"
}

variable "cluster_name" {
  type    = string
  default = "gcp-server-gke"
}

variable "artifact_repo" {
  type    = string
  default = "gcp-server"
}

variable "image_tag" {
  type    = string
  default = "v1"
}