package search

import (
	"github.com/urfave/cli"
	"github.com/zalando-incubator/pazuzu/cli/pazuzu/actions"
)

var Command = cli.Command{
	Name:      "search",
	Usage:     "Search for features in registry",
	ArgsUsage: "[regexp] - Regexp to be used for feature lookup",
	Action:    actions.Search,
}
