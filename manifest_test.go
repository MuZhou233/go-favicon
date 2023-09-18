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

// TestParseSize tests the extraction and parsing of image sizes.
func TestParseSize(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		path          string
		i             int
		width, height int
		square        bool
	}{
		// size read from manifest & markup
		{"kuli-0", "./testdata/kuli", 0, 512, 512, true}, // manifest
		{"kuli-1", "./testdata/kuli", 1, 400, 400, true}, // markup
		{"kuli-2", "./testdata/kuli", 2, 192, 192, true}, // manifest
		{"kuli-3", "./testdata/kuli", 3, 180, 180, true}, // markup

		// size read from manifest
		{"manifest-only-0", "./testdata/manifest-only", 0, 512, 512, true},
		{"manifest-only-1", "./testdata/manifest-only", 1, 192, 192, true},

		// size parsed from WxH in URL
		{"mozilla-0", "./testdata/mozilla", 0, 196, 196, true},
		{"mozilla-1", "./testdata/mozilla", 1, 180, 180, true},

		// size parsed from <link>
		{"multisize-0", "./testdata/multisize", 0, 48, 48, true},
		{"multisize-1", "./testdata/multisize", 1, 24, 24, true},
		{"multisize-2", "./testdata/multisize", 2, 16, 16, true},
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
			assert.Equal(t, td.width, icon.Width, "unexpected width")
			assert.Equal(t, td.height, icon.Height, "unexpected height")
			assert.Equal(t, td.square, icon.IsSquare(), "unexpected square")
		})
	}
}
