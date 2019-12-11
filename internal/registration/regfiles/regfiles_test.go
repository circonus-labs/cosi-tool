// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package regfiles

import (
	"os"
	"path"
	"strings"
	"testing"

	circapi "github.com/circonus-labs/go-apiclient"
	"github.com/rs/zerolog"
)

func TestSave(t *testing.T) {
	t.Log("Testing Save")

	file := "testdata/will_be_overwritten.json"
	tests := []struct {
		name       string
		file       string
		obj        interface{}
		force      bool
		clean      bool
		shouldFail bool
		expected   string
	}{
		// NOTE: an empty file emits the object to the screen (iow, a valid use case...)
		// {"invalid file empty", "", nil, false, false, true, "invalid file name (empty)"},
		{"invalid api object", file, nil, false, false, true, "invalid configuration (nil)"},
		{"invalid file", "testdata/isdir", &circapi.CheckBundle{CID: "test"}, false, false, true, "testdata/isdir is a directory"},
		{"valid create", file, &circapi.CheckBundle{CID: "test"}, false, true, false, ""},
		{"valid overwrite", file, &circapi.CheckBundle{CID: "test"}, true, false, false, ""},
		{"valid no overwrite", file, &circapi.CheckBundle{CID: "test"}, false, false, true, "testdata/will_be_overwritten.json already exists, see --force"},
		{"valid create (yaml)", strings.Replace(file, ".json", ".yaml", 1), &circapi.CheckBundle{CID: "test"}, false, true, false, ""},
		// NOTE: on toml version, pass struct not ptr to struct
		{"valid create (toml)", strings.Replace(file, ".json", ".toml", 1), circapi.CheckBundle{CID: "test"}, false, true, false, ""},
	}

	// NOTE: these are sequential, not parallel
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if test.clean {
				os.Remove(test.file)
			}
			err := Save(test.file, test.obj, test.force)
			if test.shouldFail {
				if err == nil {
					t.Fatal("expected error")
				} else if err.Error() != test.expected {
					t.Fatalf("unexpected error (%s)", err)
				}
			} else if err != nil {
				t.Fatalf("unexpected error (%s)", err)
			}
		})
	}
}

func TestFind(t *testing.T) {
	t.Log("Test Find")

	regDir := "testdata"
	tests := []struct {
		name       string
		regDir     string
		regType    string
		shouldFail bool
		expected   string
	}{
		{"invalid dir empty", "", "", true, "invalid registration directory (empty)"},
		{"invalid type empty", regDir, "", true, "invalid registration type (empty)"},
		{"invalid type (foo)", regDir, "foo", true, "invalid registration type (foo)"},
		{"invalid dir", path.Join(regDir, "invalid"), "check", true, "reading registration directory: open testdata/invalid: no such file or directory"},
		{"valid", regDir, "check", false, ""},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			files, err := Find(tst.regDir, tst.regType)
			if tst.shouldFail {
				if err == nil {
					t.Fatal("expected error")
				} else if err.Error() != tst.expected {
					t.Fatalf("unexpected error (%s)", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error (%s)", err)
				} else if len(*files) == 0 {
					t.Fatalf("expected at least 1 reg file")
				}
			}
		})
	}
}

func TestLoad(t *testing.T) {
	t.Log("Testing Load")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	{
		t.Log("invalid (nil)")
		_, err := Load("testdata/reg-valid.json", nil)
		if err == nil {
			t.Fatal("expected error")
		} else if err.Error() != "invalid interface (<nil>) - nil or not pointer to struct" {
			t.Fatalf("unexpected error (%s)", err)
		}
	}
	{
		t.Log("invalid (non-ptr)")
		var v circapi.CheckBundle
		_, err := Load("testdata/reg-valid.json", v)
		if err == nil {
			t.Fatal("expected error")
		} else if err.Error() != "invalid interface (apiclient.CheckBundle) - nil or not pointer to struct" {
			t.Fatalf("unexpected error (%s)", err)
		}
	}

	tests := []struct {
		name       string
		file       string
		exists     bool
		shouldFail bool
		expected   string
	}{
		{"invalid (empty)", "", false, true, "invalid registration file (empty)"},
		{"invalid (dir)", "testdata/registration-isdir.json", true, true, "read testdata/registration-isdir.json: is a directory"},
		{"invalid (error)", "testdata/reg-error.json", true, true, "parsing registration (testdata/reg-error.json): unexpected end of JSON input"},
		{"invalid (missing)", "testdata/registration-missing.json", false, false, ""},
		{"valid", "testdata/reg-valid.json", true, false, ""},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			var v circapi.CheckBundle
			exists, err := Load(tst.file, &v)
			if tst.shouldFail {
				if err == nil {
					t.Fatal("expected error")
				} else if err.Error() != tst.expected {
					t.Fatalf("unexpected error (%s)", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error (%s)", err)
				}
			}
			if tst.exists != exists {
				t.Fatalf("expected %v got %v", tst.exists, exists)
			}
		})
	}
}
