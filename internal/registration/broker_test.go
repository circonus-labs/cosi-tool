// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package registration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/circonus-labs/cosi-tool/internal/registration/options"
	"github.com/circonus-labs/go-apiclient"
	"github.com/rs/zerolog"
)

var (
	emptyBrokers = []apiclient.Broker{}
	emptyOptions = options.Options{
		Brokers: options.Brokers{},
		Checks:  options.Checks{},
	}
	validBrokerRx      = regexp.MustCompile(`^/broker/[123]$`)
	validOptBrokerList = []string{"/broker/1", "/broker/2"}
)

func TestSelectBroker(t *testing.T) {
	t.Log("Testing selectBroker")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	reg := &Registration{cliCirc: genMockCircAPI()}

	tests := []struct {
		name        string
		checkType   string
		shouldFail  bool
		expectedErr string
	}{
		{"invalid check type", "", true, "invalid check type (empty)"},
		{"no brokers returned", "invalid", true, "no brokers returned by API"},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			_, err := reg.selectBroker(tst.checkType)
			if tst.shouldFail {
				if err == nil {
					t.Fatal("expected error")
				}
				if err.Error() != tst.expectedErr {
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

func TestGetExplicit(t *testing.T) {
	t.Log("Testing getExplicit")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	broker := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "foo")
	}))
	defer broker.Close()
	bu, err := url.Parse(broker.URL)
	if err != nil {
		t.Fatalf("error parsing broker url (%s)", err)
	}
	bip := bu.Hostname()
	bport, err := strconv.ParseUint(bu.Port(), 10, 16)
	if err != nil {
		t.Fatalf("error parsing broker url port (%s)", err)
	}

	validBrokers := []apiclient.Broker{
		{
			CID: "/broker/1",
			Details: []apiclient.BrokerDetail{
				{
					Status:       "active",
					Modules:      []string{"json", "httptrap"},
					ExternalHost: &bip,
					ExternalPort: uint16(bport),
				},
			},
		},
		{
			CID: "/broker/2",
			Details: []apiclient.BrokerDetail{
				{
					Status:       "active",
					Modules:      []string{"json", "httptrap"},
					ExternalHost: &bip,
					ExternalPort: uint16(bport),
				},
			},
		},
	}

	reg := &Registration{
		cliCirc:               genMockCircAPI(),
		maxBrokerResponseTime: 500 * time.Millisecond,
	}

	tests := []struct {
		name        string
		checkType   string
		brokers     *[]apiclient.Broker
		options     *options.Checks
		expectValid bool
		expectBid   bool
		shouldFail  bool
		expectedErr string
	}{
		{"invalid check type", "", &emptyBrokers, &emptyOptions.Checks, false, false, true, "invalid check type (empty)"},
		{"invalid broker list", "foo", nil, &emptyOptions.Checks, false, false, true, "invalid broker list (nil)"},
		{"invalid options", "foo", &emptyBrokers, nil, false, false, true, "invalid options config (nil)"},
		{"unsupported check type", "foo", &emptyBrokers, &emptyOptions.Checks, false, false, true, "unsupported check type (foo)"},
		{"sys, pass-thru", "json", &emptyBrokers, &emptyOptions.Checks, false, false, false, ""},
		{"sys, /broker/1", "json", &validBrokers, &options.Checks{System: options.SystemCheck{BrokerID: "/broker/1"}}, true, true, false, ""},
		{"sys, 1", "json", &validBrokers, &options.Checks{System: options.SystemCheck{BrokerID: "1"}}, true, true, false, ""},
		{"grp, pass-thru", "httptrap", &emptyBrokers, &emptyOptions.Checks, false, false, false, ""},
		{"grp, /broker/1", "httptrap", &validBrokers, &options.Checks{Group: options.GroupCheck{BrokerID: "/broker/2"}}, true, true, false, ""},
		{"grp, 1", "httptrap", &validBrokers, &options.Checks{Group: options.GroupCheck{BrokerID: "2"}}, true, true, false, ""},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			valid, bid, err := reg.getExplicit(tst.checkType, tst.brokers, tst.options)
			if tst.shouldFail {
				if err == nil {
					t.Fatal("expected error")
				}
				if err.Error() != tst.expectedErr {
					t.Fatalf("unexpected error (%s)", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error (%s)", err)
				}
			}
			if tst.expectValid {
				if !valid {
					t.Fatal("expected true")
				}
			} else {
				if valid {
					t.Fatal("expected false")
				}
			}
			if tst.expectBid {
				if !validBrokerRx.MatchString(bid) {
					t.Fatalf("expected %s got (%s)", validBrokerRx.String(), bid)
				}
			} else {
				if bid != "" {
					t.Fatalf("expected empty bid got (%s)", bid)
				}
			}

		})
	}
}

func TestSelectFromConfigList(t *testing.T) {
	t.Log("Testing selectFromConfigList")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	broker := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "foo")
	}))
	defer broker.Close()
	bu, err := url.Parse(broker.URL)
	if err != nil {
		t.Fatalf("error parsing broker url (%s)", err)
	}
	bip := bu.Hostname()
	bport, err := strconv.ParseUint(bu.Port(), 10, 16)
	if err != nil {
		t.Fatalf("error parsing broker url port (%s)", err)
	}

	validBrokers := []apiclient.Broker{
		{
			CID: "/broker/1",
			Details: []apiclient.BrokerDetail{
				{
					Status:       "active",
					Modules:      []string{"json", "httptrap"},
					ExternalHost: &bip,
					ExternalPort: uint16(bport),
				},
			},
		},
		{
			CID: "/broker/2",
			Details: []apiclient.BrokerDetail{
				{
					Status:       "active",
					Modules:      []string{"json", "httptrap"},
					ExternalHost: &bip,
					ExternalPort: uint16(bport),
				},
			},
		},
	}

	reg := &Registration{
		cliCirc:               genMockCircAPI(),
		maxBrokerResponseTime: 500 * time.Millisecond,
	}

	tests := []struct {
		name        string
		checkType   string
		brokers     *[]apiclient.Broker
		options     *options.Brokers
		expectValid bool
		expectBid   bool
		shouldFail  bool
		expectedErr string
	}{
		{"invalid check type (empty)", "", &emptyBrokers, &emptyOptions.Brokers, false, false, true, "invalid check type (empty)"},
		{"invalid broker list", "foo", nil, &emptyOptions.Brokers, false, false, true, "invalid broker list (nil)"},
		{"invalid options", "foo", &emptyBrokers, nil, false, false, true, "invalid options config (nil)"},
		// system check
		{"sys no optcfg, pass-through", "json", &emptyBrokers, &emptyOptions.Brokers, false, false, false, ""},
		{"sys w/opt invalid default", "json", &validBrokers, &options.Brokers{System: options.SystemBrokers{List: validOptBrokerList, Default: -2}}, false, false, true, "invalid system check broker config in regconf (default invalid)"},
		{"sys w/opt default out of range", "json", &validBrokers, &options.Brokers{System: options.SystemBrokers{List: validOptBrokerList, Default: 3}}, false, false, true, "invalid system check broker config in regconf (default out of list range)"},
		{"sys w/opt 0", "json", &validBrokers, &options.Brokers{System: options.SystemBrokers{List: validOptBrokerList, Default: 0}}, true, true, false, ""},
		{"sys w/opt -1", "json", &validBrokers, &options.Brokers{System: options.SystemBrokers{List: validOptBrokerList, Default: -1}}, true, true, false, ""},
		// group check
		{"grp no optcfg, pass-through", "httptrap", &emptyBrokers, &emptyOptions.Brokers, false, false, false, ""},
		{"grp w/opt invalid default", "httptrap", &validBrokers, &options.Brokers{Group: options.GroupBrokers{List: validOptBrokerList, Default: -2}}, false, false, true, "invalid group check broker config in regconf (default invalid)"},
		{"grp w/opt default out of range", "httptrap", &validBrokers, &options.Brokers{Group: options.GroupBrokers{List: validOptBrokerList, Default: 3}}, false, false, true, "invalid group check broker config in regconf (default out of list range)"},
		{"grp w/opt 0", "httptrap", &validBrokers, &options.Brokers{Group: options.GroupBrokers{List: validOptBrokerList, Default: 0}}, true, true, false, ""},
		{"grp w/opt -1", "httptrap", &validBrokers, &options.Brokers{Group: options.GroupBrokers{List: validOptBrokerList, Default: -1}}, true, true, false, ""},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			valid, bid, err := reg.selectFromConfigList(tst.checkType, tst.brokers, tst.options)
			if tst.shouldFail {
				if err == nil {
					t.Fatal("expected error")
				}
				if err.Error() != tst.expectedErr {
					t.Fatalf("unexpected error (%s)", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error (%s)", err)
				}
			}
			if tst.expectValid {
				if !valid {
					t.Fatal("expected true")
				}
			} else {
				if valid {
					t.Fatal("expected false")
				}
			}
			if tst.expectBid {
				if !validBrokerRx.MatchString(bid) {
					t.Fatalf("invalid bid expected (%s) got (%s)", validBrokerRx.String(), bid)
				}
			} else {
				if bid != "" {
					t.Fatalf("expected empty bid, got (%s)", bid)
				}
			}
		})
	}
}

func TestSelectEnterprise(t *testing.T) {
	t.Log("Testing selectEnterprise")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	broker := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "foo")
	}))
	defer broker.Close()
	bu, err := url.Parse(broker.URL)
	if err != nil {
		t.Fatalf("error parsing broker url (%s)", err)
	}
	bip := bu.Hostname()
	bport, err := strconv.ParseUint(bu.Port(), 10, 16)
	if err != nil {
		t.Fatalf("error parsing broker url port (%s)", err)
	}

	noValidEntBrokers := []apiclient.Broker{
		{
			CID:  "/broker/1",
			Type: "enterprise",
			Details: []apiclient.BrokerDetail{
				{
					Status:       "active",
					Modules:      []string{"invalid"},
					ExternalHost: &bip,
					ExternalPort: uint16(bport),
				},
			},
		},
	}

	validBrokers := []apiclient.Broker{
		{
			CID:  "/broker/1",
			Type: "enterprise",
			Details: []apiclient.BrokerDetail{
				{
					Status:       "active",
					Modules:      []string{"httptrap"},
					ExternalHost: &bip,
					ExternalPort: uint16(bport),
				},
			},
		},
		{
			CID:  "/broker/2",
			Type: "enterprise",
			Details: []apiclient.BrokerDetail{
				{
					Status:       "active",
					Modules:      []string{"json"},
					ExternalHost: &bip,
					ExternalPort: uint16(bport),
				},
			},
		},
		{
			CID:  "/broker/3",
			Type: "enterprise",
			Details: []apiclient.BrokerDetail{
				{
					Status:       "active",
					Modules:      []string{"json"},
					ExternalHost: &bip,
					ExternalPort: uint16(bport),
				},
			},
		},
	}

	reg := &Registration{
		cliCirc:               genMockCircAPI(),
		maxBrokerResponseTime: 500 * time.Millisecond,
	}

	tests := []struct {
		name        string
		checkType   string
		brokers     *[]apiclient.Broker
		expectValid bool
		expectBid   bool
		shouldFail  bool
		expectedErr string
	}{
		{"invalid check type (empty)", "", &emptyBrokers, false, false, true, "invalid check type (empty)"},
		{"invalid broker list", "foo", nil, false, false, true, "invalid broker list (nil)"},
		{"no valid enterprise", "json", &noValidEntBrokers, false, false, true, "available enterprise brokers found, none valid"},
		// system check
		{"sys no enterprise, pass-through", "json", &emptyBrokers, false, false, false, ""},
		{"sys valid", "json", &validBrokers, true, true, false, ""},
		// group check
		{"grp no enterprise, pass-through", "httptrap", &emptyBrokers, false, false, false, ""},
		{"grp valid", "httptrap", &validBrokers, true, true, false, ""},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			valid, bid, err := reg.selectEnterprise(tst.checkType, tst.brokers)
			if tst.shouldFail {
				if err == nil {
					t.Fatal("expected error")
				}
				if err.Error() != tst.expectedErr {
					t.Fatalf("unexpected error (%s)", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error (%s)", err)
				}
			}
			if tst.expectValid {
				if !valid {
					t.Fatal("expected true")
				}
			} else {
				if valid {
					t.Fatal("expected false")
				}
			}
			if tst.expectBid {
				if !validBrokerRx.MatchString(bid) {
					t.Fatalf("invalid bid expected (%s) got (%s)", validBrokerRx.String(), bid)
				}
			} else {
				if bid != "" {
					t.Fatalf("expected empty bid, got (%s)", bid)
				}
			}
		})
	}
}

func TestGetCosiDefault(t *testing.T) {
	t.Log("Testing getCosiDefault")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	broker := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "foo")
	}))
	defer broker.Close()
	bu, err := url.Parse(broker.URL)
	if err != nil {
		t.Fatalf("error parsing broker url (%s)", err)
	}
	bip := bu.Hostname()
	bport, err := strconv.ParseUint(bu.Port(), 10, 16)
	if err != nil {
		t.Fatalf("error parsing broker url port (%s)", err)
	}

	validBrokers := []apiclient.Broker{
		{
			CID: "/broker/1",
			Details: []apiclient.BrokerDetail{
				{
					Status:       "active",
					Modules:      []string{"json", "httptrap"},
					ExternalHost: &bip,
					ExternalPort: uint16(bport),
				},
			},
		},
		{
			CID: "/broker/2",
			Details: []apiclient.BrokerDetail{
				{
					Status:       "active",
					Modules:      []string{"json", "httptrap"},
					ExternalHost: &bip,
					ExternalPort: uint16(bport),
				},
			},
		},
	}

	reg := &Registration{
		cliCirc:               genMockCircAPI(),
		cliCosi:               genMockCosiAPI(),
		maxBrokerResponseTime: 500 * time.Millisecond,
	}

	tests := []struct {
		name        string
		checkType   string
		brokers     *[]apiclient.Broker
		client      CosiAPI
		expectValid bool
		expectBid   bool
		shouldFail  bool
		expectedErr string
	}{
		{"invalid check type (empty)", "", &emptyBrokers, reg.cliCosi, false, false, true, "invalid check type (empty)"},
		{"invalid broker list", "foo", nil, reg.cliCosi, false, false, true, "invalid broker list (nil)"},
		{"invalid cosi api client", "foo", &emptyBrokers, nil, false, false, true, "invalid cosi API client (nil)"},
		{"invalid check type (unknown)", "foo", &emptyBrokers, reg.cliCosi, false, false, true, "unknown check type (foo)"},
		// system check
		{"sys valid", "json", &validBrokers, reg.cliCosi, true, true, false, ""},
		// group check
		{"grp valid", "httptrap", &validBrokers, reg.cliCosi, true, true, false, ""},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			valid, bid, err := reg.getCosiDefault(tst.checkType, tst.brokers, tst.client)
			if tst.shouldFail {
				if err == nil {
					t.Fatal("expected error")
				}
				if err.Error() != tst.expectedErr {
					t.Fatalf("unexpected error (%s)", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error (%s)", err)
				}
			}
			if tst.expectValid {
				if !valid {
					t.Fatal("expected true")
				}
			} else {
				if valid {
					t.Fatal("expected false")
				}
			}
			if tst.expectBid {
				if !validBrokerRx.MatchString(bid) {
					t.Fatalf("invalid bid expected (%s) got (%s)", validBrokerRx.String(), bid)
				}
			} else {
				if bid != "" {
					t.Fatalf("expected empty bid, got (%s)", bid)
				}
			}
		})
	}
}

func TestCheckBroker(t *testing.T) {
	t.Log("Testing checkBroker")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	broker := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "foo")
	}))
	defer broker.Close()
	bu, err := url.Parse(broker.URL)
	if err != nil {
		t.Fatalf("error parsing broker url (%s)", err)
	}
	bip := bu.Hostname()
	port, err := strconv.ParseUint(bu.Port(), 10, 16)
	if err != nil {
		t.Fatalf("error parsing broker url port (%s)", err)
	}
	bport := uint16(port)

	validBrokers := []apiclient.Broker{
		{
			CID: "/broker/1",
			Details: []apiclient.BrokerDetail{
				{
					Status:       "active",
					Modules:      []string{"json", "httptrap"},
					ExternalHost: &bip,
					ExternalPort: bport,
				},
			},
		},
		{
			CID: "/broker/2",
			Details: []apiclient.BrokerDetail{
				{
					Status:  "active",
					Modules: []string{"json", "httptrap"},
					IP:      &bip,
					Port:    &bport,
				},
			},
		},
	}

	reg := &Registration{
		cliCirc:               genMockCircAPI(),
		cliCosi:               genMockCosiAPI(),
		maxBrokerResponseTime: 500 * time.Millisecond,
	}

	tests := []struct {
		name        string
		checkType   string
		brokerID    string
		brokers     *[]apiclient.Broker
		expectValid bool
		expectBid   bool
		shouldFail  bool
		expectedErr string
	}{
		{"invalid check type (empty)", "", "", &emptyBrokers, false, false, true, "invalid check type (empty)"},
		{"invalid broker id (empty)", "foo", "", &emptyBrokers, false, false, true, "invalid broker id (empty)"},
		{"invalid broker list", "foo", "bar", nil, false, false, true, "invalid broker list (nil)"},
		{"invalid check type (unknown)", "foo", "1", &emptyBrokers, false, false, true, "unknown check type (foo)"},
		{"invalid broker id (foo)", "json", "foo", &emptyBrokers, false, false, true, "invalid broker id specified (foo) - format should be '#' or '/broker/#'"},
		// system check
		{"sys valid", "json", "1", &validBrokers, true, true, false, ""},
		// group check
		{"grp valid", "httptrap", "2", &validBrokers, true, true, false, ""},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			valid, bid, err := reg.checkBroker(tst.checkType, tst.brokerID, tst.brokers)
			if tst.shouldFail {
				if err == nil {
					t.Fatal("expected error")
				}
				if err.Error() != tst.expectedErr {
					t.Fatalf("unexpected error (%s)", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error (%s)", err)
				}
			}
			if tst.expectValid {
				if !valid {
					t.Fatal("expected true")
				}
			} else {
				if valid {
					t.Fatal("expected false")
				}
			}
			if tst.expectBid {
				if !validBrokerRx.MatchString(bid) {
					t.Fatalf("invalid bid expected (%s) got (%s)", validBrokerRx.String(), bid)
				}
			} else {
				if bid != "" {
					t.Fatalf("expected empty bid, got (%s)", bid)
				}
			}
		})
	}
}

func TestBrokerConnectionTest(t *testing.T) {
	t.Log("Testing brokerConnectionTest")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	broker := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "foo")
	}))
	defer broker.Close()
	bu, err := url.Parse(broker.URL)
	if err != nil {
		t.Fatalf("error parsing broker url (%s)", err)
	}
	bip := bu.Hostname()
	port, err := strconv.ParseUint(bu.Port(), 10, 16)
	if err != nil {
		t.Fatalf("error parsing broker url port (%s)", err)
	}
	bport := uint16(port)

	tests := []struct {
		name        string
		addr        string
		port        uint16
		timeout     time.Duration
		expectValid bool
		shouldFail  bool
		expectedErr string
	}{
		{"invalid address (empty)", "", 0, 500 * time.Millisecond, false, true, "invalid ip (empty)"},
		{"invalid port (0)", bip, 0, 500 * time.Millisecond, false, true, "invalid port (empty|0)"},
		{"invalid duration (0)", bip, bport, time.Duration(0), false, true, "invalid duration (0)"},
		{"valid", bip, bport, 500 * time.Millisecond, true, false, ""},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			valid, err := brokerConnectionTest(tst.addr, tst.port, tst.timeout)
			if tst.shouldFail {
				if err == nil {
					t.Fatal("expected error")
				}
				if err.Error() != tst.expectedErr {
					t.Fatalf("unexpected error (%s)", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error (%s)", err)
				}
			}
			if tst.expectValid {
				if !valid {
					t.Fatal("expected true")
				}
			} else {
				if valid {
					t.Fatal("expected false")
				}
			}
		})
	}
}
