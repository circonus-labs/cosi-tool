// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

// Package check provides the various supporting functions for the
// `cosi check *` commands
package check

import (
	"github.com/pkg/errors"
)

const (
	// KeyCID is the check bundle id
	KeyCID = "check.id"
	// DefaultCID is the default value for the id fetch option
	DefaultCID = ""

	// KeyType is the cosi check type
	KeyType = "check.type"
	// DefaultType is the default value for the type fetch option
	DefaultType = ""

	// KeyName is the display name of a check
	KeyName = "check.display_name"
	// DefaultName is the default value for the name fetch option
	DefaultName = ""

	// KeyTarget is the target of a check
	KeyTarget = "check.target"
	// DefaultTarget is the default value for the target fetch option
	DefaultTarget = ""

	// KeyOutFile is the output file for saving a fetched check configuration
	KeyOutFile = "check.out_file"
	// DefaultOutFile is the default value for the output file option
	DefaultOutFile = ""

	// KeyInFile is the input check configuration file
	KeyInFile = "check.in_file"
	// DefaultInFile is the default value for the input file option
	DefaultInFile = ""

	// KeyQuiet is a flag for limiting output
	KeyQuiet = "check.quiet"
	// DefaultQuiet is the default value for the quiet flag
	DefaultQuiet = false

	// KeyForce is a flag to force overwritting files
	KeyForce = "check.force"
	// DefaultForce is the default value for the force flag
	DefaultForce = false

	// KeyVerify is a flag to indicate list should verify check has not been modified
	KeyVerify = "check.verify"
	// DefaultVerify is the default value for the verify flag
	DefaultVerify = false

	// KeyLong is a flag indicating list should output long, more verbose, listings
	KeyLong = "check.long"
	// DefaultLong is the default value for the long flag
	DefaultLong = false
)

// isModfied fetches a check from the Circonus API, compares the last modified
// time passed to what is received, returns true if the received check's last
// modified time is different from the passed time, otherwise false
func isModified(client CircAPI, id string, lm uint) (bool, error) {
	b, err := FetchByID(client, id)
	if err != nil {
		return false, errors.Wrap(err, "fetching check from API")
	}
	if b.LastModified != lm {
		return true, nil
	}
	return false, nil
}
