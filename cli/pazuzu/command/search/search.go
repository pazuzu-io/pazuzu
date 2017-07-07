package search

import (
	"github.com/urfave/cli"
	"github.com/pazuzu-io/pazuzu/cli/pazuzu/actions"
)

var Command = cli.Command{
	Name:      "search",
	Usage:     "Search for features in registry",
	ArgsUsage: "[query] - query to be used for feature lookup for substring search in features names",
	Action:    actions.Search,
}
