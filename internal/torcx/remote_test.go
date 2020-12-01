// Copyright 2018 CoreOS Inc.
// Copyright 2020 Kinvolk GmbH
//
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

package torcx

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// TODO(lucab): add more positive tests

func TestBasicEvaluateURL(t *testing.T) {
	var r *Remote
	url := "https://example.com/basepath"

	_, err := r.evaluateURL("")
	if err != errNilRemote {
		t.Fatalf("expected %s, got %s", errNilRemote, err)
	}
	r = &Remote{}
	_, err = r.evaluateURL("")
	if err != errEmptyUsrMountpoint {
		t.Fatalf("expected %s, got %s", errEmptyUsrMountpoint, err)
	}
	_, err = r.evaluateURL("/usr")
	if err != errEmptyTemplateURL {
		t.Fatalf("expected %s, got %s", errEmptyTemplateURL, err)
	}
	r.TemplateURL = url
	res, err := r.evaluateURL("/usr")
	if err != nil {
		t.Fatalf("got unexpected error %s", err)
	}
	if res.String() != url {
		t.Fatalf("expected %s, got %s", url, res)
	}
}

func TestEvaluateURLTemplating(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "torcx_remote_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	libdir := filepath.Join(tmpDir, "lib")
	if err := os.MkdirAll(libdir, 0755); err != nil {
		t.Fatal(err)
	}
	osReleasePath := filepath.Join(libdir, "os-release")
	osContentCoreos := `
ID="coreos"
VERSION_ID="1680.2.0"
COREOS_BOARD="amd64-usr"
`

	if err := ioutil.WriteFile(osReleasePath, []byte(osContentCoreos), 0755); err != nil {
		t.Fatal(err)
	}

	basURL := "https://example.com/baseurl/"
	testCases := []struct {
		template string
		result   string
	}{
		{
			"",
			basURL,
		},
		{
			"${ID}",
			basURL + "coreos",
		},
		{
			"${VERSION_ID}",
			basURL + "1680.2.0",
		},
		{
			"${COREOS_BOARD}",
			basURL + "amd64-usr",
		},
		{
			"${FLATCAR_BOARD}",
			basURL + "amd64-usr",
		},
		{
			"${COREOS_USR}",
			basURL + tmpDir,
		},
		{
			"${FLATCAR_USR}",
			basURL + tmpDir,
		},
		{
			"${ID}/${COREOS_BOARD}/${VERSION_ID}",
			basURL + "coreos/amd64-usr/1680.2.0",
		},
		{
			"${ID}/${FLATCAR_BOARD}/${VERSION_ID}",
			basURL + "coreos/amd64-usr/1680.2.0",
		},
	}

	for _, tt := range testCases {
		template := basURL + tt.template
		r := Remote{
			TemplateURL: template,
		}
		res, err := r.evaluateURL(tmpDir)
		if err != nil {
			t.Fatalf("got unexpected error %s", err)
		}
		if res.String() != tt.result {
			t.Fatalf("using %q: expected %s, got %s", tt.template, tt.result, res)
		}
	}

	osContentFlatcar := `
NAME="Flatcar Container Linux by Kinvolk"
ID=flatcar
ID_LIKE=coreos
VERSION=2705.0.0
VERSION_ID=2705.0.0
BUILD_ID=2020-11-26-2020
PRETTY_NAME="Flatcar Container Linux by Kinvolk 2705.0.0 (Oklo)"
ANSI_COLOR="38;5;75"
HOME_URL="https://flatcar-linux.org/"
BUG_REPORT_URL="https://issues.flatcar-linux.org"
FLATCAR_BOARD="amd64-usr"
`
	if err := ioutil.WriteFile(osReleasePath, []byte(osContentFlatcar), 0755); err != nil {
		t.Fatal(err)
	}

	basURL = "https://example.com/baseurl/"
	testCases = []struct {
		template string
		result   string
	}{
		{
			"",
			basURL,
		},
		{
			"${ID}",
			basURL + "flatcar",
		},
		{
			"${VERSION_ID}",
			basURL + "2705.0.0",
		},
		{
			"${COREOS_BOARD}",
			basURL + "amd64-usr",
		},
		{
			"${FLATCAR_BOARD}",
			basURL + "amd64-usr",
		},
		{
			"${COREOS_USR}",
			basURL + tmpDir,
		},
		{
			"${FLATCAR_USR}",
			basURL + tmpDir,
		},
		{
			"${ID}/${COREOS_BOARD}/${VERSION_ID}",
			basURL + "flatcar/amd64-usr/2705.0.0",
		},
		{
			"${ID}/${FLATCAR_BOARD}/${VERSION_ID}",
			basURL + "flatcar/amd64-usr/2705.0.0",
		},
	}

	for _, tt := range testCases {
		template := basURL + tt.template
		r := Remote{
			TemplateURL: template,
		}
		res, err := r.evaluateURL(tmpDir)
		if err != nil {
			t.Fatalf("got unexpected error %s", err)
		}
		if res.String() != tt.result {
			t.Fatalf("using %q: expected %s, got %s", tt.template, tt.result, res)
		}
	}

}
