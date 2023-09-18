// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-11-09

package favicon

import (
	"encoding/json"
	"io"
	urls "net/url"
	"path/filepath"
	"regexp"
	"strconv"
)

// Manifest is the relevant parts of a manifest.json file.
type Manifest struct {
	Icons []ManifestIcon `json:"icons"`
}

// ManifestIcon is an icon from a manifest.json file.
type ManifestIcon struct {
	URL      string `json:"src"`
	Type     string `json:"type"`
	RawSizes string `json:"sizes"`
}

type size struct {
	w, h int
}

func (p *parser) parseManifest(url string) []*Icon {
	p.find.log.Printf("loading manifest %q ...", url)
	rc, err := p.find.fetchURL(url)
	if err != nil {
		p.find.log.Printf("[ERROR] parse manifest: %v", err)
		return nil
	}
	defer rc.Close()

	return p.parseManifestReader(rc)
}

func (p *parser) parseManifestReader(r io.Reader) []*Icon {
	var (
		icons []*Icon
		man   = Manifest{}
		err   error
	)

	dec := json.NewDecoder(r)
	if err = dec.Decode(&man); err != nil {
		p.find.log.Printf("[ERROR] parse manifest: %v", err)
	}
	for _, mi := range man.Icons {
		// TODO: make URL relative to manifest, not page
		mi.URL = p.absURL(mi.URL)
		p.find.log.Printf("(manifest) %s", mi.URL)
		for _, sz := range parseSizes(mi.RawSizes) {
			icon := &Icon{
				URL:    mi.URL,
				Width:  sz.w,
				Height: sz.h,
			}
			icons = append(icons, icon)
		}
	}

	return icons
}

var (
	rxSize  = regexp.MustCompile(`(\d+)x(\d+)`)
	rxWidth = regexp.MustCompile(`-(\d+)$`)
)

func parseSizes(s string) []size {
	m := rxSize.FindAllStringSubmatch(s, -1)
	if m == nil {
		return nil
	}
	var sizes []size
	for _, l := range m {
		for i := 1; i < len(l)-1; i += 2 {
			w, _ := strconv.ParseInt(l[i], 10, 32)
			h, _ := strconv.ParseInt(l[i+1], 10, 32)
			sizes = append(sizes, size{w: int(w), h: int(h)})
		}
	}
	return sizes
}

// find dimensions in URL.
func extractSizeFromURL(url string) *size {
	// try to find WxH pattern
	v := parseSizes(url)
	if len(v) > 0 {
		return &v[0]
	}

	// look for -NNN at end of filename
	u, err := urls.Parse(url)
	if err != nil {
		return nil
	}

	var (
		name = filepath.Base(u.Path)
		ext  = filepath.Ext(name)
	)

	if m := rxWidth.FindStringSubmatch(name[:len(name)-len(ext)]); m != nil {
		n, _ := strconv.ParseInt(m[1], 10, 32)
		return &size{w: int(n), h: int(n)}
	}

	return nil
}
