// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package broker

import (
	"io/ioutil"
	"testing"

	"github.com/rs/zerolog"
)

func TestDefault(t *testing.T) {
	t.Log("Testing Default")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Log("\tinvalid state (nil dest)")
	if err := Default(testCosiBrokers.URL, nil, 0, ""); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "invalid destination (nil)" {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid state (missing custom options)")
	if err := Default(testCosiBrokers.URL, ioutil.Discard, 0, "testdata/missing.json"); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "cosi broker default: custom options file: custom options file: open testdata/missing.json: no such file or directory" {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid state (no json, no fallback)")
	if err := Default(testCosiBrokers.URL+"/nojson/nofallback/", ioutil.Discard, 0, ""); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "cosi broker default: no json (or fallback) default brokers defined" {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid state (no trap, no fallback)")
	if err := Default(testCosiBrokers.URL+"/notrap/nofallback/", ioutil.Discard, 0, ""); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "cosi broker default: no trap (or fallback) default brokers defined" {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid state (no json, use fallback)")
	if err := Default(testCosiBrokers.URL+"/nojson/", ioutil.Discard, 0, ""); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid state (no trap, use fallback)")
	if err := Default(testCosiBrokers.URL+"/notrap/", ioutil.Discard, 0, ""); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid state (custom options explicit broker)")
	if err := Default(testCosiBrokers.URL, ioutil.Discard, 0, "testdata/explicit.json"); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid state (custom options trap & json set)")
	if err := Default(testCosiBrokers.URL, ioutil.Discard, 0, "testdata/full.json"); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid state (custom options no trap)")
	if err := Default(testCosiBrokers.URL, ioutil.Discard, 0, "testdata/notrap.json"); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid state (custom options no trap, cosi no trap)")
	if err := Default(testCosiBrokers.URL+"/notrap/", ioutil.Discard, 0, "testdata/notrap.json"); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid state (custom options no json)")
	if err := Default(testCosiBrokers.URL, ioutil.Discard, 0, "testdata/nojson.json"); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid state (custom options no json, cosi no json)")
	if err := Default(testCosiBrokers.URL+"/nojson/", ioutil.Discard, 0, "testdata/nojson.json"); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid state (cosi-install --broker)")
	if err := Default(testCosiBrokers.URL, ioutil.Discard, 1, ""); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid state (trap and json from cosi api defaults)")
	if err := Default(testCosiBrokers.URL, ioutil.Discard, 0, ""); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	}
}
