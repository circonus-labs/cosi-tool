// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package options

import (
	"path"
	"testing"

	"github.com/circonus-labs/cosi-tool/internal/config"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func TestLoadConfigFile(t *testing.T) {
	t.Log("Testing LoadConfigFile")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	viper.Set(config.KeyCosiID, "abc123")
	tests := []struct {
		name       string
		file       string
		shouldFail bool
		expected   string
	}{
		{"invalid readerr", path.Join("testdata", "cust-config-readerr"), true, "reading config file: read testdata/cust-config-readerr.toml: is a directory"},
		{"valid missing", path.Join("testdata", "missing"), false, ""}, // no error, just empty config
		{"valid (no config)", "", false, ""}, // should just get an empty config back
		{"valid", path.Join("..", "..", "..", "etc", "example-reg-conf"), false, ""},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			_, err := LoadConfigFile(tst.file)
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
		})
	}
}
