package command

import (
	"github.com/zalando-incubator/pazuzu/cli/pazuzu/command/clean"
	"github.com/zalando-incubator/pazuzu/cli/pazuzu/command/config"
	"github.com/zalando-incubator/pazuzu/cli/pazuzu/command/search"
	"github.com/zalando-incubator/pazuzu/cli/pazuzu/command/compose"
	"github.com/zalando-incubator/pazuzu/cli/pazuzu/command/build"
)

var (
	Clean   = clean.Command
	Config  = config.Command
	Search  = search.Command
	Build   = build.Command
	Compose = compose.Command
)
