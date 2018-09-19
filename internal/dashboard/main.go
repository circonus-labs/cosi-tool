// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

// Package dashboard provides the various supporting functions for the
// `cosi dashboard *` commands
package dashboard

import "github.com/pkg/errors"

const (
	// KeyCID is the dashboard id
	KeyCID = "dashboard.id"
	// CIDDefault is the default value for the id fetch option
	CIDDefault = ""

	// KeyType is the cosi dashboard type
	KeyType = "dashboard.type"
	// TypeDefault is the default value for the type fetch option
	TypeDefault = ""

	// KeyTitle is the title of a dashboard
	KeyTitle = "dashboard.title"
	// TitleDefault is the default value for the name fetch option
	TitleDefault = ""

	// KeyOutFile is the output file for saving a fetched dashboard configuration
	KeyOutFile = "dashboard.out_file"
	// OutFileDefault is the default value for the output file option
	OutFileDefault = ""

	// KeyInFile is the input dashboard configuration file
	KeyInFile = "dashboard.in_file"
	// InFileDefault is the default value for the input file option
	InFileDefault = ""

	// KeyQuiet is a flag for limiting output
	KeyQuiet = "dashboard.quiet"
	// QuietDefault is the default value for the quiet flag
	QuietDefault = false

	// KeyForce is a flag to force overwritting files
	KeyForce = "dashboard.force"
	// ForceDefault is the default value for the force flag
	ForceDefault = false

	// KeyVerify is a flag to indicate list should verify dashboard has not been modified
	KeyVerify = "dashboard.verify"
	// VerifyDefault is the default value for the verify flag
	VerifyDefault = false

	// KeyLong is a flag indicating list should output long, more verbose, listings
	KeyLong = "dashboard.long"
	// LongDefault is the default value for the long flag
	LongDefault = false
)

// isModfied fetches a check from the Circonus API, compares the last modified
// time passed to what is received, returns true if the received check's last
// modified time is different from the passed time, otherwise false
func isModified(client API, id string, lm uint) (bool, error) {
	b, err := FetchByID(client, id)
	if err != nil {
		return false, errors.Wrap(err, "fetching dashboard from API")
	}
	if b.LastModified != lm {
		return true, nil
	}
	return false, nil
}
