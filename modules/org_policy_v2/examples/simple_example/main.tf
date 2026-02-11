/**
 * Copyright 2026 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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