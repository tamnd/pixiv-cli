// Package cli assembles the pixiv command tree from the pixiv
// domain on top of the any-cli/kit framework.
package cli

import (
	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/pixiv-cli/pixiv"
)

// Build metadata, set via -ldflags at release time.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// NewApp assembles the kit application from the pixiv domain. The
// domain's Register installs the client factory and every operation, so the
// binary and a host (ant, which blank-imports the package) share one source of
// truth. kit.Run turns the App into the CLI, plus the serve and mcp surfaces and
// the typed-error-to-exit-code mapping.
func NewApp() *kit.App {
	id := pixiv.Domain{}.Info().Identity
	id.Version = Version

	app := kit.New(id)
	(pixiv.Domain{}).Register(app)
	app.AddCommand(newVersionCmd())
	return app
}
