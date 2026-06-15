package pixiv

import (
	"context"
	"fmt"

	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/any-cli/kit/errs"
)

func init() { kit.Register(Domain{}) }

// Domain is the pixiv kit driver.
type Domain struct{}

// Info describes the scheme, hostnames, and binary identity.
func (Domain) Info() kit.DomainInfo {
	return kit.DomainInfo{
		Scheme: "pixiv",
		Hosts:  []string{Host},
		Identity: kit.Identity{
			Binary: "pixiv",
			Short:  "A command line for Pixiv artwork rankings.",
			Long: `A command line for Pixiv artwork rankings.

pixiv reads public Pixiv data over plain HTTPS and prints output that pipes
into the rest of your tools. No API key or account required.`,
			Site: Host,
			Repo: "https://github.com/tamnd/pixiv-cli",
		},
	}
}

// Register installs the client factory and every operation onto app.
func (Domain) Register(app *kit.App) {
	app.SetClient(newClient)

	kit.Handle(app, kit.OpMeta{Name: "ranking", Group: "read", List: true,
		Summary: "Fetch the Pixiv illustration ranking",
		Args: []kit.Arg{
			{Name: "mode", Help: "ranking mode: daily, weekly, monthly, rookie, original (default: daily)", Optional: true},
		},
	}, listRanking)

	kit.Handle(app, kit.OpMeta{Name: "modes", Group: "read", List: true,
		Summary: "List available ranking modes and content types",
	}, listModes)
}

// newClient builds the client from kit config.
func newClient(_ context.Context, cfg kit.Config) (any, error) {
	c := DefaultConfig()
	if cfg.UserAgent != "" {
		c.UserAgent = cfg.UserAgent
	}
	if cfg.Rate > 0 {
		c.Rate = cfg.Rate
	}
	if cfg.Retries > 0 {
		c.Retries = cfg.Retries
	}
	if cfg.Timeout > 0 {
		c.Timeout = cfg.Timeout
	}
	return NewClient(c), nil
}

// --- inputs ---

type rankingInput struct {
	Mode    string  `kit:"arg" help:"ranking mode: daily, weekly, monthly, rookie, original"`
	Content string  `kit:"flag" help:"content type: illust, manga, ugoira (default: illust)"`
	Page    int     `kit:"flag" help:"page number (default: 1)"`
	Limit   int     `kit:"flag,inherit" help:"max results"`
	Client  *Client `kit:"inject"`
}

type modesInput struct{}

// --- handlers ---

func listRanking(ctx context.Context, in rankingInput, emit func(*Illust) error) error {
	mode := in.Mode
	if mode == "" {
		mode = "daily"
	}
	content := in.Content
	if content == "" {
		content = "illust"
	}
	page := in.Page
	if page <= 0 {
		page = 1
	}

	illusts, err := in.Client.Ranking(ctx, mode, content, page, in.Limit)
	if err != nil {
		return mapErr(err)
	}
	for i := range illusts {
		if err := emit(&illusts[i]); err != nil {
			return err
		}
	}
	return nil
}

func listModes(_ context.Context, _ modesInput, emit func(*ModeInfo) error) error {
	for _, m := range Modes() {
		mc := m
		if err := emit(&mc); err != nil {
			return err
		}
	}
	return nil
}

// Classify turns any accepted input into the canonical (type, id).
func (Domain) Classify(input string) (uriType, id string, err error) {
	return "", "", errs.Usage("pass an artwork URL like https://www.pixiv.net/en/artworks/ID")
}

// Locate is the inverse: the live https URL for a (type, id).
func (Domain) Locate(uriType, id string) (string, error) {
	if uriType != "artwork" {
		return "", errs.Usage("pixiv has no resource type %q", uriType)
	}
	return fmt.Sprintf("%s/en/artworks/%s", BaseURL, id), nil
}

// mapErr converts a library error into the kit error kind.
func mapErr(err error) error {
	return err
}
