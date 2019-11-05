# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# 	https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
resource "google_cloudfunctions_function" "remove-public-ip" {
  name                  = "RemovePublicIP"
  description           = "Removes all the external IP addresses of a GCE instance."
  runtime               = "go111"
  available_memory_mb   = 128
  source_archive_bucket = "${var.setup.gcf-bucket-name}"
  source_archive_object = "${var.setup.gcf-object-name}"
  timeout               = 60
  project               = "${var.setup.automation-project}"
  region                = "${var.setup.region}"
  entry_point           = "RemovePublicIP"

  event_trigger {
    event_type = "providers/cloud.pubsub/eventTypes/topic.publish"
    resource   = "${var.setup.cscc-notifications-topic-prefix}-topic"
  }
}

# Required to retrieve ancestry for projects within this folder.
resource "google_folder_iam_member" "roles-viewer" {
  count = length(var.folder-ids)

  folder = "folders/${var.folder-ids[count.index]}"
  role   = "roles/viewer"
  member = "serviceAccount:${var.setup.automation-service-account}"
}

# Required to delete the access config (IP) from the network interface of the GCE instance.
resource "google_folder_iam_member" "roles-instance-admin-v1" {
  count = length(var.folder-ids)

  folder = "folders/${var.folder-ids[count.index]}"
  role   = "roles/compute.instanceAdmin.v1"
  member = "serviceAccount:${var.setup.automation-service-account}"
}