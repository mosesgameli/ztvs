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

package main

import (
	"context"
	"github.com/mosesgameli/ztvs-sdk-go/sdk"
)

type GoCheck struct{}

func (c *GoCheck) ID() string   { return "go_test_check" }
func (c *GoCheck) Name() string { return "Go Fixture Check" }

func (c *GoCheck) Run(ctx context.Context) (*sdk.Finding, error) {
	return &sdk.Finding{
		ID:          "F-GO-001",
		Severity:    "info",
		Title:       "Go Fixture Executed",
		Description: "This is a successful result from the Go polyglot fixture.",
	}, nil
}

func main() {
	RunPlugin()
}

func RunPlugin() {
	sdk.Run(sdk.Metadata{
		Name:       "go-fixture",
		Version:    "1.0.0",
		APIVersion: 1,
	}, []sdk.Check{&GoCheck{}})
}
