// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package worksheet

const (
	// KeyCID is the worksheet id
	KeyCID = "worksheet.id"
	// CIDDefault is the default value for the id fetch option
	CIDDefault = ""

	// KeyType is the cosi worksheet type
	KeyType = "worksheet.type"
	// TypeDefault is the default value for the type fetch option
	TypeDefault = ""

	// KeyTitle is the title of a worksheet
	KeyTitle = "worksheet.title"
	// TitleDefault is the default value for the name fetch option
	TitleDefault = ""

	// KeyOutFile is the output file for saving a fetched worksheet configuration
	KeyOutFile = "worksheet.out_file"
	// OutFileDefault is the default value for the output file option
	OutFileDefault = ""

	// KeyInFile is the input worksheet configuration file
	KeyInFile = "worksheet.in_file"
	// InFileDefault is the default value for the input file option
	InFileDefault = ""

	// KeyQuiet is a flag for limiting output
	KeyQuiet = "worksheet.quiet"
	// QuietDefault is the default value for the quiet flag
	QuietDefault = false

	// KeyForce is a flag to force overwritting files
	KeyForce = "worksheet.force"
	// ForceDefault is the default value for the force flag
	ForceDefault = false

	// KeyVerify is a flag to indicate list should verify worksheet has not been modified
	KeyVerify = "worksheet.verify"
	// VerifyDefault is the default value for the verify flag
	VerifyDefault = false

	// KeyLong is a flag indicating list should output long, more verbose, listings
	KeyLong = "worksheet.long"
	// LongDefault is the default value for the long flag
	LongDefault = false
)
