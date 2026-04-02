// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

type PluginMetadata struct {
	Name          string `json:"name"`
	LatestVersion string `json:"latest_version"`
	Repo          string `json:"repo"`
	Checksum      string `json:"checksum"`
	Signature     string `json:"signature"`
	AuditStatus   string `json:"audit_status"`
}

type Index struct {
	Version string           `json:"version"`
	Plugins []PluginMetadata `json:"plugins"`
}
