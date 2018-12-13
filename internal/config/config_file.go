// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	toml "github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// LoadConfigFile will attempt to load json|toml|yaml configuration files.
// `base` is the full path and base name of the configuration file to load.
// `target` is an interface in to which the data will be loaded. Checks for
// '<base>.json', '<base>.toml', and '<base>.yaml'.
func LoadConfigFile(base string, target interface{}) error {

	if base == "" {
		return errors.Errorf("invalid config file (empty)")
	}

	rt := reflect.ValueOf(target)
	if rt.Kind() != reflect.Ptr || rt.IsNil() {
		return errors.Errorf("invalid target (%s) not a pointer", reflect.TypeOf(target).String())
	}

	limitExt := ""
	extensions := []string{".yaml", ".json", ".toml"}
	loaded := false

	// check if base already has a known extension
	// remove it and limit the list of extensions the
	// indicated one
	for _, ext := range extensions {
		if strings.HasSuffix(base, ext) {
			base = strings.Replace(base, ext, "", -1)
			limitExt = ext
			break
		}
	}
	if limitExt != "" {
		extensions = []string{limitExt}
	}

	for _, ext := range extensions {
		cfg := base + ext
		data, err := ioutil.ReadFile(cfg)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return errors.Wrap(err, "reading config file")
		}
		parseErrMsg := fmt.Sprintf("parsing config file (%s)", cfg)
		switch ext {
		case ".json":
			if err := json.Unmarshal(data, target); err != nil {
				return errors.Wrap(err, parseErrMsg)
			}
			loaded = true
		case ".toml":
			if err := toml.Unmarshal(data, target); err != nil {
				return errors.Wrap(err, parseErrMsg)
			}
			loaded = true
		case ".yaml":
			if err := yaml.Unmarshal(data, target); err != nil {
				return errors.Wrap(err, parseErrMsg)
			}
			loaded = true
		}
	}

	if !loaded {
		return errors.Errorf("no config found matching (%s%s)", base, strings.Join(extensions, "|"))
	}

	return nil
}

// LoadConfigFile will attempt to load json|toml|yaml configuration files.
// `base` is the full path and base name of the configuration file to load.
// `target` is an interface in to which the data will be loaded. Checks for
// '<base>.json', '<base>.toml', and '<base>.yaml'.
func SaveConfigFile(cfgFile string, target interface{}, force bool) error {
	if cfgFile == "" {
		return errors.Errorf("invalid config file (empty)")
	}
	if target == nil {
		return errors.Errorf("invalid target (nil)")
	}

	// rt := reflect.ValueOf(target)
	// if rt.Kind() != reflect.Ptr || rt.IsNil() {
	// 	return errors.Errorf("invalid target (%s) not a pointer", reflect.TypeOf(target).String())
	// }

	ext := filepath.Ext(cfgFile)

	var err error
	var data []byte
	switch ext {
	case ".json":
		data, err = json.MarshalIndent(target, "", "  ")
	case ".toml":
		data, err = toml.Marshal(target)
	case ".yaml":
		data, err = yaml.Marshal(target)
	default:
		return errors.Errorf("unknown format requested (%s)", ext)
	}

	if err != nil {
		return errors.Wrap(err, "formatting configuration")
	}

	return ioutil.WriteFile(cfgFile, data, 0644)
}
