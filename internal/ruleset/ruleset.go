// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

// Package ruleset handles interacting with api for rulesets
package ruleset

const (
	// KeyCID is the ruleset id
	KeyCID = "ruleset.id"
	// DefaultCID is the default value for the id fetch option
	DefaultCID = ""

	// KeyOutFile is the output file for saving a fetched ruleset configuration
	KeyOutFile = "ruleset.out_file"
	// DefaultOutFile is the default value for the output file option
	DefaultOutFile = ""

	// KeyInFile is the input ruleset configuration file
	KeyInFile = "ruleset.in_file"
	// DefaultInFile is the default value for the input file option
	DefaultInFile = ""

	// KeyForce is a flag to force overwritting files
	KeyForce = "ruleset.force"
	// DefaultForce is the default value for the force flag
	DefaultForce = false

	// KeyQuiet is a flag for limiting output
	KeyQuiet = "ruleset.quiet"
	// DefaultQuiet is the default value for the quiet flag
	DefaultQuiet = false

	// KeyLong is a flag indicating list should output long, more verbose, listings
	KeyLong = "ruleset.long"
	// DefaultLong is the default value for the long flag
	DefaultLong = false
)
