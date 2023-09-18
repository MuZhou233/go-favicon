// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-11-09

package favicon_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/muzhou233/go-favicon"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type debugLogger struct {
	t *testing.T
}

func (l debugLogger) Printf(format string, v ...interface{}) {
	l.t.Logf(format, v...)
}

// TestBaseURL verifies absolute links.
func TestBaseURL(t *testing.T) {
	t.Parallel()
	file, err := os.Open("testdata/kuli/index.html")
	require.Nil(t, err, "unexpected error")
	defer file.Close()

	f := favicon.New(
		favicon.WithLogger(debugLogger{t}),
		favicon.OnlyICO,
		favicon.IgnoreWellKnown,
		favicon.IgnoreManifest,
	)
	require.Nil(t, err, "unexpected error")

	var (
		baseURL = "https://www.kulturliste-duesseldorf.de"
		x       = "https://www.kulturliste-duesseldorf.de/favicon-rot.ico"
		icons   []*favicon.Icon
	)
	icons, err = f.FindReader(file, baseURL)
	require.Nil(t, err, "unexpected error")
	// for _, i := range icons {
	// 	fmt.Println(i)
	// }
	require.Equal(t, 1, len(icons), "unexpected favicon count")
	assert.Equal(t, x, icons[0].URL, "unexpected favicon URL")
}

// TestFindHTML parses HTML only.
func TestFindHTML(t *testing.T) {
	t.Parallel()
	file, err := os.Open("testdata/github/index.html")
	require.Nil(t, err, "unexpected error")
	defer file.Close()

	f := favicon.New(favicon.WithLogger(debugLogger{t}))
	require.Nil(t, err, "unexpected error")

	var icons []*favicon.Icon
	icons, err = f.FindReader(file)
	require.Nil(t, err, "unexpected error")
	assert.Equal(t, 6, len(icons), "unexpected favicon count")
}

// TestHTTP tests fetching via HTTP.
func TestHTTP(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, path string
		xcount     int
	}{
		{"github", "./testdata/github", 17},
		{"kuli", "./testdata/kuli", 7},
		{"mozilla", "./testdata/mozilla", 4},
		{"no-markup", "./testdata/no-markup", 3},
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
			assert.Equal(t, td.xcount, len(icons), "unexpected favicon count")
		})
	}
}

// TestIgnore verifies Ignore* Options.
func TestIgnore(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, path      string
		ignoreWellKnown bool
		ignoreManifest  bool
		xcount          int
	}{
		// ignore well-known
		{"github-ignore-well-known", "./testdata/github", true, false, 17},
		{"kuli-ignore-well-known", "./testdata/kuli", true, false, 7},
		{"mozilla-ignore-well-known", "./testdata/mozilla", true, false, 4},
		{"no-markup-ignore-well-known", "./testdata/no-markup", true, false, 2},
		{"manifest-only-ignore-well-known", "./testdata/manifest-only", true, false, 2},

		// ignore manifest
		{"no-markup-ignore-manifest", "./testdata/no-markup", false, true, 1},
		{"manifest-only-ignore-manifest", "./testdata/manifest-only", false, true, 0},

		// ignore well-known & manifest
		{"github-ignore-both", "./testdata/github", true, true, 6},
		{"kuli-ignore-both", "./testdata/kuli", true, true, 5},
		{"mozilla-ignore-both", "./testdata/mozilla", true, true, 4},
		{"no-markup-ignore-both", "./testdata/no-markup", true, true, 0},
		{"manifest-only-both", "./testdata/manifest-only", true, true, 0},
	}

	for _, td := range tests {
		td := td
		t.Run(td.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(http.FileServer(http.Dir(td.path)))
			defer ts.Close()

			opts := []favicon.Option{
				favicon.WithClient(ts.Client()),
				favicon.WithLogger(debugLogger{t}),
			}

			if td.ignoreWellKnown {
				opts = append(opts, favicon.IgnoreWellKnown)
			}
			if td.ignoreManifest {
				opts = append(opts, favicon.IgnoreManifest)
			}

			f := favicon.New(opts...)
			icons, err := f.Find(ts.URL + "/index.html")
			require.Nil(t, err, "unexpected error")
			assert.Equal(t, td.xcount, len(icons), "unexpected favicon count")
		})
	}
}

// TestFilter verifies filtering Options.
func TestFilter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, path string
		opts       []favicon.Option
		xcount     int
	}{
		{"no-options", "./testdata/multiformat", []favicon.Option{}, 9},
		{"only-square", "./testdata/multiformat", []favicon.Option{favicon.OnlySquare}, 6},
		{"ignore-nosize", "./testdata/multiformat", []favicon.Option{favicon.IgnoreNoSize}, 8},
		{"only-ico", "./testdata/multiformat", []favicon.Option{favicon.OnlyICO}, 1},
		{"only-png", "./testdata/multiformat", []favicon.Option{favicon.OnlyPNG}, 7},
		{"only-square-png", "./testdata/multiformat", []favicon.Option{favicon.OnlyPNG, favicon.OnlySquare}, 4},
		{"only-jpeg", "./testdata/multiformat", []favicon.Option{favicon.OnlyMimeType("image/jpeg")}, 1},
		{"only-square-sized", "./testdata/multiformat", []favicon.Option{favicon.OnlySquare, favicon.IgnoreNoSize}, 5},
		{"only-400", "./testdata/multiformat", []favicon.Option{favicon.MinWidth(400), favicon.MaxWidth(400)}, 1},
		{"width-100-and-200", "./testdata/multiformat", []favicon.Option{favicon.MinWidth(100), favicon.MaxWidth(200)}, 4},
		{"width+height-100-and-200", "./testdata/multiformat", []favicon.Option{favicon.MinWidth(100),
			favicon.MaxWidth(200), favicon.MinHeight(100), favicon.MaxHeight(200)}, 2},
	}

	for _, td := range tests {
		td := td
		t.Run(td.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(http.FileServer(http.Dir(td.path)))
			defer ts.Close()

			opts := []favicon.Option{
				favicon.WithClient(ts.Client()),
				favicon.WithLogger(debugLogger{t}),
			}
			opts = append(opts, td.opts...)
			f := favicon.New(opts...)
			icons, err := f.Find(ts.URL + "/index.html")
			require.Nil(t, err, "unexpected error")
			for _, i := range icons {
				t.Log(i)
			}
			assert.Equal(t, td.xcount, len(icons), "unexpected favicon count")
		})
	}
}
