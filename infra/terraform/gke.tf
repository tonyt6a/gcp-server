resource "google_container_cluster" "gke" {
  name             = var.cluster_name
  location         = var.region
  enable_autopilot = true

  ip_allocation_policy {}
}