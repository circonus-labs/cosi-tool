// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package defaults

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/circonus-labs/cosi-tool/internal/release"
	"github.com/rs/zerolog/log"
)

const (
	// COSIURL is the cosi-server URL
	COSIURL = "https://onestep.circonus.com/"

	// APIURL is the Circonus API URL
	APIURL = "https://api.circonus.com/v2/"

	// AgentMode defines the mode for the check connecting to the agent
	AgentMode = "reverse"
	// AgentURL is the URL the agent is listening on
	AgentURL = "http://localhost:2609/"

	// HostBrokerID is the id of a specific broker to use when creating checks (an id or 0 which means auto select)
	HostBrokerID = 0
	// HostBrokerType is a specific type of broker to limit to when selecting a broker automatically
	// (any|enterprise)
	HostBrokerType = "any"

	// Debug is false by default
	Debug = false

	// LogLevel set to info by default
	LogLevel = "info"

	// LogPretty colored/formatted output to stderr
	LogPretty = true
)

var (
	// BasePath is the "base" directory
	//
	// expected installation structure:
	// base        (e.g. /opt/circonus/cosi)
	//   /bin      (e.g. /opt/circonus/cosi/bin)
	//   /etc      (e.g. /opt/circonus/cosi/etc)
	//   /registration  (e.g. /opt/circonus/cosi/registration)
	//   /log     (e.g. /opt/circonus/cosi/log)
	BasePath = ""

	// EtcPath returns the default etc directory within base directory
	// (e.g. /opt/circonus/cosi/etc)
	EtcPath = ""

	// ConfigFile defines the default configuration file name
	ConfigFile = ""

	// RegPath defines the registration files path
	RegPath = ""

	// HostCheckTarget is the check target
	HostCheckTarget = ""

	// RegConf defines the registration options configuration file
	RegConf = ""
)

func init() {
	var exePath string
	var resolvedExePath string
	var err error

	exePath, err = os.Executable()
	if err == nil {
		resolvedExePath, err = filepath.EvalSymlinks(exePath)
		if err == nil {
			BasePath = filepath.Clean(filepath.Join(filepath.Dir(resolvedExePath), ".."))
		}
	}

	if err != nil {
		fmt.Printf("Unable to determine path to binary %v\n", err)
		os.Exit(1)
	}

	EtcPath = filepath.Join(BasePath, "etc")
	RegPath = filepath.Join(BasePath, "registration")

	ConfigFile = filepath.Join(EtcPath, release.NAME+".yaml")

	RegConf = filepath.Join(EtcPath, "regconf")

	hn, err := os.Hostname()
	if err != nil {
		log.Fatal().Err(err).Msg("obtaining hostname from OS")
	}
	HostCheckTarget = hn
}
