// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package regfiles

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	toml "github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	yaml "gopkg.in/yaml.v2"
)

var regTypeValidator = regexp.MustCompile(`^(check|graph|worksheet|dashboard|ruleset)$`)

// Save registration write file w/optional force overwrite or write
// formatted JSON to stdout if no file name is provided.
func Save(file string, o interface{}, force bool) error {
	if o == nil {
		return errors.New("invalid configuration (nil)")
	}

	format := filepath.Ext(file)
	if format == "" {
		format = ".json"
	}

	var data []byte
	{
		var err error

		switch format {
		case ".json":
			data, err = json.MarshalIndent(o, "", "  ")
		case ".toml":
			data, err = toml.Marshal(o)
		case ".yaml":
			data, err = yaml.Marshal(o)
		default:
			err = errors.Errorf("unknown extension/format (%s)", format)
		}
		if err != nil {
			return errors.Wrap(err, "formatting configuration")
		}
	}

	if file == "" {
		fmt.Fprintf(os.Stdout, "%s\n", string(data))
		return nil
	}

	flags := os.O_WRONLY
	if s, serr := os.Stat(file); serr != nil {
		if os.IsNotExist(serr) {
			flags |= os.O_CREATE
		} else {
			return errors.Wrapf(serr, "stat %s", file)
		}
	} else if s.IsDir() {
		return errors.Errorf("%s is a directory", file)
	} else if s.Mode().IsRegular() {
		if force {
			flags |= os.O_CREATE
		} else {
			return errors.Errorf("%s already exists, see --force", file)
		}
	}

	f, err := os.OpenFile(file, flags, 0644)
	if err != nil {
		return errors.Wrap(err, "saving configuration")
	}
	defer f.Close()

	_, err = f.Write(data)

	return err
}

// Find returns a list of registration files from the registration directory
// which match the specified registration type (e.g. check|dashboard|graph|worksheet)
func Find(regDir, regType string) (*[]string, error) {
	if regDir == "" {
		return nil, errors.Errorf("invalid registration directory (empty)")
	}
	if regType == "" {
		return nil, errors.Errorf("invalid registration type (empty)")
	}
	if !regTypeValidator.MatchString(regType) {
		return nil, errors.Errorf("invalid registration type (%s)", regType)
	}

	files, err := ioutil.ReadDir(regDir)
	if err != nil {
		return nil, errors.Wrap(err, "reading registration directory")
	}

	regFileSig := "registration-" + regType
	var regFiles []string

	for _, file := range files {
		if !file.Mode().IsRegular() {
			continue
		}

		if !strings.HasPrefix(file.Name(), regFileSig) {
			continue
		}

		regFiles = append(regFiles, file.Name())
	}

	return &regFiles, nil
}

// Load a registration file into the destination interface, returns boolean
// indicating if the file was found and any error reading/parsing the file.
func Load(regFile string, v interface{}) (bool, error) {
	if regFile == "" {
		return false, errors.Errorf("invalid registration file (empty)")
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return false, errors.Errorf("invalid interface (%v) - nil or not pointer to struct", reflect.TypeOf(v))
	}

	if _, err := os.Stat(regFile); os.IsNotExist(err) {
		log.Warn().Err(err).Str("file", regFile).Msg("not found")
		return false, nil
	}

	data, err := ioutil.ReadFile(regFile)
	if err != nil {
		log.Warn().Err(err).Str("file", regFile).Msg("not read")
		return !os.IsNotExist(err), err
	}

	if err := json.Unmarshal(data, v); err != nil {
		log.Warn().Err(err).Str("file", regFile).Msg("parsing")
		return true, errors.Wrapf(err, "parsing registration (%s)", regFile)
	}

	return true, nil
}
