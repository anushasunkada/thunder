/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package enginebridge

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func ensureServerTestAssets(t *testing.T, serverHome string) {
	t.Helper()

	securityDir := filepath.Join(serverHome, "repository", "resources", "security")
	if err := os.MkdirAll(securityDir, 0o750); err != nil {
		t.Fatalf("create security dir: %v", err)
	}

	cryptoKey := filepath.Join(securityDir, "crypto.key")
	if _, err := os.Stat(cryptoKey); err != nil {
		//nolint:gosec // test-only key generation with fixed arguments
		out, err := exec.Command("openssl", "rand", "-hex", "32").Output()
		if err != nil {
			t.Skipf("openssl not available for test crypto key: %v", err)
		}
		out = []byte(strings.TrimSpace(string(out)))
		if err := os.WriteFile(cryptoKey, out, 0o600); err != nil {
			t.Fatalf("write crypto key: %v", err)
		}
	}

	for _, prefix := range []string{"server", "signing"} {
		ensureTestCertificate(t, securityDir, prefix)
	}
}

func ensureTestCertificate(t *testing.T, dir, prefix string) {
	t.Helper()

	certPath := filepath.Join(dir, prefix+".cert")
	keyPath := filepath.Join(dir, prefix+".key")
	if _, err := os.Stat(certPath); err == nil {
		if _, err := os.Stat(keyPath); err == nil {
			return
		}
	}

	subject := fmt.Sprintf("/O=WSO2/OU=ThunderID/CN=localhost-%s", prefix)
	//nolint:gosec // test-only certificate generation with fixed arguments
	cmd := exec.Command("openssl", "req", "-x509", "-nodes", "-days", "365", "-newkey", "rsa:2048",
		"-keyout", keyPath, "-out", certPath, "-subj", subject)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("openssl not available for test certificates: %v (%s)", err, out)
	}
}
