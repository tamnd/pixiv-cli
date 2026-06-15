package pixiv

import (
	"testing"

	"github.com/tamnd/any-cli/kit"
)

func TestDomainInfo(t *testing.T) {
	info := Domain{}.Info()
	if info.Scheme != "pixiv" {
		t.Errorf("Scheme = %q, want pixiv", info.Scheme)
	}
	if len(info.Hosts) == 0 || info.Hosts[0] != Host {
		t.Errorf("Hosts = %v, want [%s]", info.Hosts, Host)
	}
	if info.Identity.Binary != "pixiv" {
		t.Errorf("Identity.Binary = %q, want pixiv", info.Identity.Binary)
	}
}

func TestLocate(t *testing.T) {
	got, err := Domain{}.Locate("artwork", "12345678")
	want := "https://www.pixiv.net/en/artworks/12345678"
	if err != nil || got != want {
		t.Errorf("Locate = (%q, %v), want (%q, nil)", got, err, want)
	}
}

func TestHostWiring(t *testing.T) {
	h, err := kit.Open()
	if err != nil {
		t.Fatal(err)
	}

	info, ok := h.Domain("pixiv")
	if !ok {
		t.Fatal("pixiv domain not registered")
	}
	if info.Identity.Binary != "pixiv" {
		t.Errorf("Binary = %q, want pixiv", info.Identity.Binary)
	}
}
