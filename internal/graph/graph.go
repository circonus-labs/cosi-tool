// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

// Package graph handles interacting with API graph endpoint
package graph

const (
	// KeyCID is the graph id
	KeyCID = "graph.id"
	// CIDDefault is the default value for the id fetch option
	CIDDefault = ""

	// KeyType is the cosi graph type
	KeyType = "graph.type"
	// TypeDefault is the default value for the type fetch option
	TypeDefault = ""

	// KeyTitle is the title of a graph
	KeyTitle = "graph.title"
	// TitleDefault is the default value for the name fetch option
	TitleDefault = ""

	// KeyOutFile is the output file for saving a fetched graph configuration
	KeyOutFile = "graph.out_file"
	// OutFileDefault is the default value for the output file option
	OutFileDefault = ""

	// KeyInFile is the input graph configuration file
	KeyInFile = "graph.in_file"
	// InFileDefault is the default value for the input file option
	InFileDefault = ""

	// KeyQuiet is a flag for limiting output
	KeyQuiet = "graph.quiet"
	// QuietDefault is the default value for the quiet flag
	QuietDefault = false

	// KeyForce is a flag to force overwritting files
	KeyForce = "graph.force"
	// ForceDefault is the default value for the force flag
	ForceDefault = false

	// KeyVerify is a flag to indicate list should verify graph has not been modified
	KeyVerify = "graph.verify"
	// VerifyDefault is the default value for the verify flag
	VerifyDefault = false

	// KeyLong is a flag indicating list should output long, more verbose, listings
	KeyLong = "graph.long"
	// LongDefault is the default value for the long flag
	LongDefault = false
)
