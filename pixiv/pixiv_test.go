package pixiv_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tamnd/pixiv-cli/pixiv"
)

func TestRanking(t *testing.T) {
	payload := map[string]any{
		"contents": []map[string]any{
			{
				"rank":              1,
				"title":             "Test Artwork",
				"date":              "2026-06-14 00:00:00",
				"url":               "https://i.pximg.net/img-master/img/test.jpg",
				"user_name":         "TestArtist",
				"user_id":           "12345",
				"illust_id":         99999,
				"illust_type":       "0",
				"illust_page_count": "1",
				"tags":              []string{"original", "art"},
				"width":             1200,
				"height":            800,
			},
			{
				"rank":              2,
				"title":             "Second Artwork",
				"date":              "2026-06-14 00:00:00",
				"url":               "https://i.pximg.net/img-master/img/test2.jpg",
				"user_name":         "AnotherArtist",
				"user_id":           "67890",
				"illust_id":         88888,
				"illust_type":       "0",
				"illust_page_count": float64(3),
				"tags":              []string{"fanart"},
				"width":             800,
				"height":            1200,
			},
		},
		"date":       "20260614",
		"next_date":  "20260615",
		"prev_date":  "20260613",
		"rank_total": 500,
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the Referer header is set
		if r.Header.Get("Referer") == "" {
			http.Error(w, "missing referer", http.StatusForbidden)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
	defer ts.Close()

	cfg := pixiv.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.Rate = 0

	client := pixiv.NewClient(cfg)
	illusts, err := client.Ranking(context.Background(), "daily", "illust", 1, 0)
	if err != nil {
		t.Fatalf("Ranking: %v", err)
	}

	if len(illusts) != 2 {
		t.Fatalf("expected 2 illusts, got %d", len(illusts))
	}

	first := illusts[0]
	if first.Rank != 1 {
		t.Errorf("rank: got %d, want 1", first.Rank)
	}
	if first.Title != "Test Artwork" {
		t.Errorf("title: got %q, want %q", first.Title, "Test Artwork")
	}
	if first.Artist != "TestArtist" {
		t.Errorf("artist: got %q, want %q", first.Artist, "TestArtist")
	}
	if first.Tags != "original, art" {
		t.Errorf("tags: got %q, want %q", first.Tags, "original, art")
	}
	if first.IllustID != "99999" {
		t.Errorf("illust_id: got %q, want %q", first.IllustID, "99999")
	}
	if first.URL != "https://www.pixiv.net/artworks/99999" {
		t.Errorf("url: got %q", first.URL)
	}

	second := illusts[1]
	if second.Pages != 3 {
		t.Errorf("pages: got %d, want 3", second.Pages)
	}
}

func TestRankingLimit(t *testing.T) {
	contents := make([]map[string]any, 10)
	for i := range contents {
		contents[i] = map[string]any{
			"rank":              i + 1,
			"title":             "Art",
			"user_name":         "Artist",
			"illust_id":         1000 + i,
			"illust_page_count": "1",
			"tags":              []string{},
		}
	}
	payload := map[string]any{
		"contents":   contents,
		"rank_total": 500,
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
	defer ts.Close()

	cfg := pixiv.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.Rate = 0

	client := pixiv.NewClient(cfg)
	illusts, err := client.Ranking(context.Background(), "daily", "illust", 1, 5)
	if err != nil {
		t.Fatalf("Ranking: %v", err)
	}
	if len(illusts) != 5 {
		t.Fatalf("expected 5 illusts with limit, got %d", len(illusts))
	}
}

func TestModes(t *testing.T) {
	modes := pixiv.Modes()
	if len(modes) == 0 {
		t.Fatal("Modes() returned empty list")
	}
	for _, m := range modes {
		if m.Mode == "" {
			t.Errorf("mode entry has empty Mode: %+v", m)
		}
		if m.Content == "" {
			t.Errorf("mode entry has empty Content: %+v", m)
		}
		if m.Description == "" {
			t.Errorf("mode entry has empty Description: %+v", m)
		}
	}
}
