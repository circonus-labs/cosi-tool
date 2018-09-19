// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package graphs

import (
	"regexp"

	"github.com/rs/zerolog/log"
)

type globalFilters struct {
	include []*regexp.Regexp
	exclude []*regexp.Regexp
}

func compileFilters(filterList []string) []*regexp.Regexp {
	filters := make([]*regexp.Regexp, 0)
	if len(filterList) == 0 {
		return filters
	}

	for _, filterRx := range filterList {
		rx, err := regexp.Compile(filterRx)
		if err != nil {
			log.Warn().Err(err).Str("filter", filterRx).Msg("bad filter")
			continue
		}
		filters = append(filters, rx)
	}

	return filters
}
