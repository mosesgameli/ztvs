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
	"os"
	"testing"
)

func TestRunPlugin(t *testing.T) {
	// We mock stdin to simulate a handshake request
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	
	r, w, _ := os.Pipe()
	os.Stdin = r
	
	// Write handshake request and then close writing end
	go func() {
		w.Write([]byte(`{"jsonrpc":"2.0","id":"1","method":"handshake","params":{}}`))
		w.Close()
	}()
	
	// Redirect stdout to avoid polluting test output
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()
	_, wOut, _ := os.Pipe()
	os.Stdout = wOut
	
	// Invoke RunPlugin (this will process the handshake and then return because we closed stdin)
	RunPlugin()
}

func TestGoCheck_Metadata(t *testing.T) {
	c := &GoCheck{}
	if c.ID() != "go_test_check" {
		t.Errorf("expected ID go_test_check, got %s", c.ID())
	}
	if c.Name() != "Go Fixture Check" {
		t.Errorf("expected Name Go Fixture Check, got %s", c.Name())
	}
}

func TestGoCheck_Run(t *testing.T) {
	c := &GoCheck{}
	finding, err := c.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if finding.ID != "F-GO-001" {
		t.Errorf("expected finding ID F-GO-001, got %s", finding.ID)
	}
}
