variable "project_id" {
  type        = string
  description = "Project ID from setup"
}

module "org_policy" {
  source         = "../../"
  policy_root    = "project"
  policy_root_id = var.project_id
  constraint     = "compute.vmExternalIpAccess"
  rules = [
    {
      enforcement = true
      allow       = []
      deny        = []
    }
  ]
}

# Add this output so the test framework can "see" the project_id during verification
output "project_id" {
  value = var.project_id
}