package pixiv

import (
	"fmt"
	"strings"
)

// Illust is the record emitted for a ranked illustration.
type Illust struct {
	Rank     int    `json:"rank"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Tags     string `json:"tags"`
	Type     string `json:"type"`
	Pages    int    `json:"pages"`
	Date     string `json:"date"`
	IllustID string `json:"illust_id"`
	URL      string `json:"url"`
}

// ModeInfo describes one ranking mode / content pair.
type ModeInfo struct {
	Mode        string `json:"mode"`
	Content     string `json:"content"`
	Description string `json:"description"`
}

// Modes returns every supported ranking combination.
func Modes() []ModeInfo {
	return []ModeInfo{
		{Mode: "daily", Content: "illust", Description: "Top illustrations today"},
		{Mode: "daily", Content: "manga", Description: "Top manga today"},
		{Mode: "daily", Content: "ugoira", Description: "Top animations today"},
		{Mode: "weekly", Content: "illust", Description: "Top illustrations this week"},
		{Mode: "weekly", Content: "manga", Description: "Top manga this week"},
		{Mode: "weekly", Content: "ugoira", Description: "Top animations this week"},
		{Mode: "monthly", Content: "illust", Description: "Top illustrations this month"},
		{Mode: "monthly", Content: "manga", Description: "Top manga this month"},
		{Mode: "rookie", Content: "illust", Description: "Top rookie illustrations"},
		{Mode: "rookie", Content: "manga", Description: "Top rookie manga"},
		{Mode: "original", Content: "illust", Description: "Top original illustrations"},
	}
}

// ─── wire types ──────────────────────────────────────────────────────────────

type wireResponse struct {
	Contents  []wireIllust `json:"contents"`
	Date      string       `json:"date"`
	NextDate  string       `json:"next_date"`
	PrevDate  string       `json:"prev_date"`
	RankTotal int          `json:"rank_total"`
}

type wireIllust struct {
	Rank            int      `json:"rank"`
	Title           string   `json:"title"`
	Date            string   `json:"date"`
	URL             string   `json:"url"`
	UserName        string   `json:"user_name"`
	UserID          string   `json:"user_id"`
	IllustID        int      `json:"illust_id"`
	IllustType      string   `json:"illust_type"`
	IllustPageCount any      `json:"illust_page_count"` // string or number
	Tags            []string `json:"tags"`
	Width           int      `json:"width"`
	Height          int      `json:"height"`
}

func wireToIllust(w wireIllust) Illust {
	pages := 1
	switch v := w.IllustPageCount.(type) {
	case float64:
		pages = int(v)
	case string:
		if v != "" {
			var n int
			_, _ = fmt.Sscanf(v, "%d", &n)
			if n > 0 {
				pages = n
			}
		}
	}

	illustType := w.IllustType
	if illustType == "" {
		illustType = "illust"
	}

	artworkURL := fmt.Sprintf("https://www.pixiv.net/artworks/%d", w.IllustID)

	return Illust{
		Rank:     w.Rank,
		Title:    w.Title,
		Artist:   w.UserName,
		Tags:     strings.Join(w.Tags, ", "),
		Type:     illustType,
		Pages:    pages,
		Date:     w.Date,
		IllustID: fmt.Sprintf("%d", w.IllustID),
		URL:      artworkURL,
	}
}
