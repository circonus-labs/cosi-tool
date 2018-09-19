// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package graphs

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestCompileFilters(t *testing.T) {
	t.Log("Testing compileFilters")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	tests := []struct {
		name        string
		filters     []string
		expectedLen int
	}{
		{"no filters", []string{}, 0},
		{"1 filter", []string{`lo0`}, 1},
		{"1 filter (1 good, 1 bad)", []string{`lo0`, `foo(bar]`}, 1},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			f := compileFilters(tst.filters)
			if len(f) != tst.expectedLen {
				t.Fatalf("unexpected number of filters (%d)", len(f))
			}
		})
	}
}
