// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package templates

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	agentapi "github.com/circonus-labs/circonus-agent/api"
	cosiapi "github.com/circonus-labs/cosi-server/api"
	"github.com/circonus-labs/cosi-tool/internal/config"
	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Templates defines the template object
type Templates struct {
	client CosiAPI
}

// FetchAllResult defines the result from a fetching attempt when
// fetching all templates.
type FetchAllResult struct {
	IDX      int               // input templateList index
	Template *cosiapi.Template // a fetched template or nil
	Err      error             // an error or nil
}

const (
	// KeyID is the template id
	KeyID = "template.id"
	// IDDefault is the default value for the id fetch option
	IDDefault = ""

	// KeyIDList is the template list for fetch all
	KeyIDList = "template.id_list"

	// KeyShow displays the list of templates 'fetch all' would attempt to fetch
	KeyShow = "template.show"
	// ShowDefault is the default value for the show flag
	ShowDefault = false

	// KeyOutFile is the output file for saving a fetched template configuration
	KeyOutFile = "template.out_file"
	// OutFileDefault is the default value for the output file option
	OutFileDefault = ""

	// KeyQuiet is a flag for limiting output
	KeyQuiet = "template.quiet"
	// QuietDefault is the default value for the quiet flag
	QuietDefault = false

	// KeyForce is a flag to force overwritting files
	KeyForce = "template.force"
	// ForceDefault is the default value for the force flag
	ForceDefault = false
)

var (
	// IDListDefault is the base list of default templates.
	// Based on the metrics returned by the running agent, additional
	// graph templates are appended.
	IDListDefault = []string{"check-system", "dashboard-system", "worksheet-system"}
)

// New returns a new templates instance
func New(client CosiAPI) (*Templates, error) {
	if client == nil {
		osType := viper.GetString(config.KeySystemOSType)
		osDist := viper.GetString(config.KeySystemOSDistro)
		osVers := viper.GetString(config.KeySystemOSVersion)
		sysArch := viper.GetString(config.KeySystemArch)
		cosiURL := viper.GetString(config.KeyCosiURL)

		if osType == "" {
			return nil, errors.Errorf("invalid os type (empty)")
		}
		if osDist == "" {
			return nil, errors.Errorf("invalid os distro (empty)")
		}
		if osVers == "" {
			return nil, errors.Errorf("invalid os version (empty)")
		}
		if sysArch == "" {
			return nil, errors.Errorf("invalid system arch (empty)")
		}
		if cosiURL == "" {
			return nil, errors.Errorf("invalid cosi url (empty)")
		}

		cli, err := cosiapi.New(&cosiapi.Config{
			OSType:    osType,
			OSDistro:  osDist,
			OSVersion: osVers,
			SysArch:   sysArch,
			CosiURL:   cosiURL,
		})
		if err != nil {
			return nil, errors.Wrap(err, "creating cosi-server client")
		}
		client = cli
	}

	t := &Templates{
		client: client,
	}

	return t, nil
}

// Fetch retrieves a template using the cosi-server API
func (t *Templates) Fetch(id string) (*cosiapi.Template, error) {
	if id == "" {
		return nil, errors.Errorf("invalid id (empty)")
	}

	template, err := t.client.FetchTemplate(id)
	if err != nil {
		return nil, err
	}

	return template, nil
}

// FetchAll retrieves all templates specified using the cosi-server API
func (t *Templates) FetchAll(templateList []string) (*[]FetchAllResult, error) {
	if len(templateList) == 0 {
		return nil, errors.Errorf("invalid template list (empty)")
	}

	ret := make([]FetchAllResult, 0, len(templateList))
	for idx, id := range templateList {
		template, err := t.Fetch(id)
		ret = append(ret, FetchAllResult{idx, template, err})
	}

	return &ret, nil
}

// Load returns a cosi template. It will load the template specified by <id>
// if found in the dir. If not, it will fetch the template from the cosi api
// and save it into the dir.
func (t *Templates) Load(dir string, id string) (*cosiapi.Template, bool, error) {
	if dir == "" {
		return nil, false, errors.Errorf("invalid directory (empty)")
	}
	if id == "" {
		return nil, false, errors.New("invalid id (empty)")
	}

	fn := path.Join(dir, "template-"+id+cosiapi.TemplateFileExtension)

	data, err := ioutil.ReadFile(fn)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, true, errors.Wrap(err, "reading template")
		}
		// not found, retrieve from cosi api
		tmpl, ferr := t.Fetch(id)
		if ferr != nil {
			return nil, !strings.Contains(ferr.Error(), "404 Not Found"), ferr
		}
		if err := regfiles.Save(fn, *tmpl, true); err != nil { // NOTE: deref ptr, toml.Marshal can't take &struct{}
			return nil, false, err
		}
		return tmpl, true, nil
	}

	tv := cosiapi.Template{}
	if err := toml.Unmarshal(data, &tv); err != nil {
		return nil, true, errors.Wrap(err, "parsing template")
	}

	return &tv, true, nil
}

// DefaultTemplateList builds a list of default templates based on system and
// active plugins from the agent
func DefaultTemplateList(metrics *agentapi.Metrics) (*[]string, error) {
	if metrics == nil {
		return nil, errors.New("invalid metric list (nil)")
	}
	if len(*metrics) == 0 {
		return nil, errors.New("invalid metric list (empty)")
	}
	list := IDListDefault
	groups := make(map[string]bool)
	for mname := range *metrics {
		nameParts := strings.Split(mname, "`")
		metricGroup := mname
		if len(nameParts) > 1 {
			metricGroup = nameParts[0]
		}
		groups["graph-"+metricGroup] = true
	}
	for mg := range groups {
		list = append(list, mg)
	}
	return &list, nil
}
