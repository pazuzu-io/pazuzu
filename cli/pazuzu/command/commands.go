package command

import (
	"github.com/pazuzu-io/pazuzu/cli/pazuzu/command/config"
	"github.com/pazuzu-io/pazuzu/cli/pazuzu/command/project"
	"github.com/pazuzu-io/pazuzu/cli/pazuzu/command/search"
)

var (
	Config  = config.Command
	Project = project.Command
	Search  = search.Command
)
