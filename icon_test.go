// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-11-10

package favicon_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/muzhou233/go-favicon"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFormat tests the extraction and parsing of file extensions.
// NOTE: MIME types aren't tested because Go uses the system MIME
// database, so results differ between machines.
func TestFormat(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		path string
		i    int
		ext  string
	}{
		// read from manifest & markup
		{"kuli-0", "./testdata/kuli", 0, "png"}, // manifest
		{"kuli-1", "./testdata/kuli", 1, "png"}, // markup
		{"kuli-2", "./testdata/kuli", 2, "png"}, // manifest
		{"kuli-3", "./testdata/kuli", 3, "png"}, // markup
		{"kuli-6", "./testdata/kuli", 6, "ico"}, // /favicon.ico

		// read from manifest
		{"manifest-only-0", "./testdata/manifest-only", 0, "png"},
		{"manifest-only-1", "./testdata/manifest-only", 1, "png"},

		// size parsed from WxH in URL
		{"mozilla-0", "./testdata/mozilla", 0, "png"},
		{"mozilla-1", "./testdata/mozilla", 1, "png"},

		// size parsed from <link>
		{"multisize-0", "./testdata/multisize", 0, "ico"},
		{"multisize-1", "./testdata/multisize", 1, "ico"},
		{"multisize-2", "./testdata/multisize", 2, "ico"},
	}

	for _, td := range tests {
		td := td
		t.Run(td.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(http.FileServer(http.Dir(td.path)))
			defer ts.Close()

			f := favicon.New(favicon.WithLogger(debugLogger{t}))
			icons, err := f.Find(ts.URL + "/index.html")
			require.Nil(t, err, "unexpected error")
			require.Greater(t, len(icons), td.i, "too few icons found")
			icon := icons[td.i]
			assert.Equal(t, td.ext, icon.FileExt, "unexpected extension")
		})
	}
}

// TestIconCopy verifies that icon copies are the same as the original.
func TestIconCopy(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, path string
	}{
		{"github", "./testdata/github"},
		{"kuli", "./testdata/kuli"},
		{"mozilla", "./testdata/mozilla"},
		{"no-markup", "./testdata/no-markup"},
	}

	for _, td := range tests {
		td := td
		t.Run(td.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(http.FileServer(http.Dir(td.path)))
			defer ts.Close()

			f := favicon.New(favicon.WithClient(ts.Client()), favicon.WithLogger(debugLogger{t}))
			icons, err := f.Find(ts.URL + "/index.html")
			require.Nil(t, err, "unexpected error")

			for _, icon := range icons {
				i := icon.Copy()
				assert.Equal(t, icon, i, "unequal copy")
			}
		})
	}
}
