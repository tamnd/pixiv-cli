package cli

import (
	"github.com/spf13/cobra"
	"github.com/tamnd/pixiv-cli/pixiv"
)

func (a *App) rankingCmd() *cobra.Command {
	var mode, content string
	var page int

	cmd := &cobra.Command{
		Use:   "ranking",
		Short: "Fetch Pixiv ranking",
		Long:  "Fetch the Pixiv illustration or manga ranking for a given mode and content type.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			n := a.effectiveLimit(50)
			a.progressf("fetching %s %s ranking (page %d, limit %d)...", mode, content, page, n)
			illusts, err := a.client.Ranking(cmd.Context(), mode, content, page, n)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.renderOrEmpty(illusts, len(illusts))
		},
	}

	cmd.Flags().StringVar(&mode, "mode", "daily", "ranking mode: daily|weekly|monthly|rookie|original")
	cmd.Flags().StringVar(&content, "content", "illust", "content type: illust|manga|ugoira")
	cmd.Flags().IntVar(&page, "page", 1, "page number (each page has up to 50 items)")
	return cmd
}

func (a *App) modesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "modes",
		Short: "List available ranking modes and content types",
		RunE: func(cmd *cobra.Command, _ []string) error {
			modes := pixiv.Modes()
			return a.render(modes)
		},
	}
}
