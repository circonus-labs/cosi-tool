// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

// Package reset handles resetting all cosi created assets
package reset

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	circapi "github.com/circonus-labs/go-apiclient"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

const (
	// KeyForce skips prompt and assumes 'yes'
	KeyForce = "reset.force"
	// DefaultForce is the default value for the force option
	DefaultForce = false
)

// Reset provides an interactive prompt for resetting all assets created (checks, visuals, etc.)
func Reset(client CircAPI, regDir string, force bool) error {
	if client == nil {
		return errors.New("invalid client (nil)")
	}
	if regDir == "" {
		return errors.New("invalid regdir (empty)")
	}

	proceed := force
	if !proceed {
		for {
			color.Yellow("Reset will remove the check and all visuals, continue? (type 'yes' or 'no')")
			var response string
			_, err := fmt.Scanln(&response)
			if err != nil {
				color.Red(err.Error())
				os.Exit(1)
			}
			if response == "yes" {
				proceed = true
				break
			} else if response == "no" {
				proceed = false
				break
			}
			color.Red("type 'yes' or 'no'")
		}
	}

	if proceed {
		for _, fn := range []func(CircAPI, string) error{
			deleteWorksheets,
			deleteDashboards,
			deleteGraphs,
			deleteRulesets,
			deleteChecks,
			deleteTemplates,
		} {
			err := fn(client, regDir)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func removeRegistration(regFile string) error {
	if regFile == "" {
		return errors.New("invalid regfile (empty)")
	}
	color.Cyan("\tDeleting %s\n", regFile)
	return os.Remove(regFile)
}

func deleteWorksheets(client CircAPI, regDir string) error {
	if client == nil {
		return errors.New("invalid client (nil)")
	}
	if regDir == "" {
		return errors.New("invalid regdir (empty)")
	}
	assetType := "worksheet"
	assets, err := regfiles.Find(regDir, assetType)
	if err != nil {
		return errors.Wrapf(err, "loading '%s' registrations", assetType)
	}
	if len(*assets) == 0 {
		return nil
	}
	color.HiWhite("Processing %s(s)\n", assetType)
	for _, asset := range *assets {
		regFile := filepath.Join(regDir, asset)
		var v circapi.Worksheet
		ok, err := regfiles.Load(regFile, &v)
		if err != nil {
			return err
		}
		if ok {
			color.Cyan("\tRemoving %s - %s\n", assetType, v.CID)
			if _, err := client.DeleteWorksheetByCID(circapi.CIDType(&v.CID)); err != nil {
				return err
			}
			if err := removeRegistration(regFile); err != nil {
				return err
			}
		}
	}
	return nil
}

func deleteDashboards(client CircAPI, regDir string) error {
	if client == nil {
		return errors.New("invalid client (nil)")
	}
	if regDir == "" {
		return errors.New("invalid regdir (empty)")
	}
	assetType := "dashboard"
	assets, err := regfiles.Find(regDir, assetType)
	if err != nil {
		return errors.Wrapf(err, "loading '%s' registrations", assetType)
	}
	if len(*assets) == 0 {
		return nil
	}
	color.HiWhite("Processing %s(s)\n", assetType)
	for _, asset := range *assets {
		regFile := filepath.Join(regDir, asset)
		var v circapi.Dashboard
		ok, err := regfiles.Load(regFile, &v)
		if err != nil {
			return err
		}
		if ok {
			color.Cyan("\tRemoving %s - %s\n", assetType, v.CID)
			if _, err := client.DeleteDashboardByCID(circapi.CIDType(&v.CID)); err != nil {
				return err
			}
			if err := removeRegistration(regFile); err != nil {
				return err
			}
		}
	}
	return nil
}

func deleteGraphs(client CircAPI, regDir string) error {
	if client == nil {
		return errors.New("invalid client (nil)")
	}
	if regDir == "" {
		return errors.New("invalid regdir (empty)")
	}
	assetType := "graph"
	assets, err := regfiles.Find(regDir, assetType)
	if err != nil {
		return errors.Wrapf(err, "loading '%s' registrations", assetType)
	}
	if len(*assets) == 0 {
		return nil
	}
	color.HiWhite("Processing %s(s)\n", assetType)
	for _, asset := range *assets {
		regFile := filepath.Join(regDir, asset)
		var v circapi.Graph
		ok, err := regfiles.Load(regFile, &v)
		if err != nil {
			return err
		}
		if ok {
			color.Cyan("\tRemoving %s - %s\n", assetType, v.CID)
			if _, err := client.DeleteGraphByCID(circapi.CIDType(&v.CID)); err != nil {
				return err
			}
			if err := removeRegistration(regFile); err != nil {
				return err
			}
		}
	}
	return nil
}

func deleteRulesets(client CircAPI, regDir string) error {
	if client == nil {
		return errors.New("invalid client (nil)")
	}
	if regDir == "" {
		return errors.New("invalid regdir (empty)")
	}
	assetType := "ruleset"
	assets, err := regfiles.Find(regDir, assetType)
	if err != nil {
		return errors.Wrapf(err, "loading '%s' registrations", assetType)
	}
	if len(*assets) == 0 {
		return nil
	}
	color.HiWhite("Processing %s(s)\n", assetType)
	for _, asset := range *assets {
		regFile := filepath.Join(regDir, asset)
		var v circapi.RuleSet
		ok, err := regfiles.Load(regFile, &v)
		if err != nil {
			return err
		}
		if ok {
			color.Cyan("\tRemoving %s - %s\n", assetType, v.CID)
			if _, err := client.DeleteRuleSetByCID(circapi.CIDType(&v.CID)); err != nil {
				return err
			}
			if err := removeRegistration(regFile); err != nil {
				return err
			}
		}
	}
	return nil
}

func deleteChecks(client CircAPI, regDir string) error {
	if client == nil {
		return errors.New("invalid client (nil)")
	}
	if regDir == "" {
		return errors.New("invalid regdir (empty)")
	}
	assetType := "check"
	assets, err := regfiles.Find(regDir, assetType)
	if err != nil {
		return errors.Wrapf(err, "loading '%s' registrations", assetType)
	}
	if len(*assets) == 0 {
		return nil
	}
	color.HiWhite("Processing %s(s)\n", assetType)
	for _, asset := range *assets {
		regFile := filepath.Join(regDir, asset)
		var v circapi.CheckBundle
		ok, err := regfiles.Load(regFile, &v)
		if err != nil {
			return err
		}
		if ok {
			color.Cyan("\tRemoving %s - %s\n", assetType, v.CID)
			if _, err := client.DeleteCheckBundleByCID(circapi.CIDType(&v.CID)); err != nil {
				return err
			}
			if err := removeRegistration(regFile); err != nil {
				return err
			}
		}
	}
	return nil
}

func deleteTemplates(client CircAPI, regDir string) error {
	if client == nil {
		return errors.New("invalid client (nil)")
	}
	if regDir == "" {
		return errors.New("invalid regdir (empty)")
	}

	files, err := ioutil.ReadDir(regDir)
	if err != nil {
		return errors.Wrap(err, "reading registration directory")
	}
	if len(files) == 0 {
		return nil
	}
	color.HiWhite("Processing templates\n")
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "template-") {
			tfile := filepath.Join(regDir, file.Name())
			color.Cyan("\tDeleting %s\n", tfile)
			if err := os.Remove(tfile); err != nil {
				return err
			}
		}
	}
	return nil
}
