package command

import (
	"github.com/zalando-incubator/pazuzu/cli/pazuzu/command/config"
	"github.com/zalando-incubator/pazuzu/cli/pazuzu/command/project"
	"github.com/zalando-incubator/pazuzu/cli/pazuzu/command/search"
)

var (
	Config  = config.Command
	Project = project.Command
	Search  = search.Command
)
