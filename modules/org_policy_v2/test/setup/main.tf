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

resource "random_id" "suffix" {
  byte_length = 4
}

locals {
  project_id = "tmp-policy-test-${random_id.suffix.hex}"
}

variable "folder_id" {
  type        = string
  description = "The folder ID where the test project will be created."
}

variable "billing_account" {
  type        = string
  description = "The billing account to associate with the test project."
}

resource "google_project" "main" {
  name            = local.project_id
  project_id      = local.project_id
  folder_id       = var.folder_id
  billing_account = var.billing_account
  # This prevents the project from being deleted if the API call fails
  deletion_policy = "DELETE" 
}

resource "google_project_service" "main" {
  for_each = toset([
    "orgpolicy.googleapis.com",
    "compute.googleapis.com",
  ])
  project            = google_project.main.project_id
  service            = each.key
  disable_on_destroy = false
}

output "project_id" {
  value = google_project.main.project_id
}