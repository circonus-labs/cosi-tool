// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package graph

import (
	"io/ioutil"
	"testing"

	"github.com/spf13/viper"
)

func TestList(t *testing.T) {
	t.Log("Test List")

	t.Log("\tempty regdir")
	if err := List(nil, ioutil.Discard, "", "", false, false); err == nil {
		t.Fatal("expected error")
	}

	client := genMockClient()

	t.Log("\tinvalid regdir")
	if err := List(client, ioutil.Discard, "", "testdata/invalid/", false, false); err == nil {
		t.Fatal("expected error")
	}

	t.Log("\tvalid regdir (short)")
	if err := List(client, ioutil.Discard, "", "testdata/", false, false); err != nil {
		t.Fatalf("expected NO error, got %v", err)
	}

	t.Log("\tvalid regdir (long)")
	if err := List(client, ioutil.Discard, "http://example.com/", "testdata/", false, true); err != nil {
		t.Fatalf("expected NO error, got %v", err)
	}

	viper.Reset()
}

func TestGetDetail(t *testing.T) {
	t.Log("Testing getDetail")

	uiURL := "http://example.com/"

	t.Log("\tempty regdir (empty)")
	if _, err := getDetail("", "", uiURL); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "invalid registration directory (empty)" {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tinvalid regfile (empty)")
	if _, err := getDetail("testdata/", "", uiURL); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "invalid graph registration file (empty)" {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tinvalid regfile (json parse)")
	if _, err := getDetail("testdata/", "bad.json", uiURL); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "parsing graph registration json: unexpected end of JSON input" {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid regfile")
	if _, err := getDetail("testdata/", "registration-graph-test.json", uiURL); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	}
}
