// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-11-10

package favicon_test

import (
	"fmt"

	"github.com/muzhou233/go-favicon"
)

// Find favicons using default options.
func ExampleNew() {
	// Find icons defined in HTML, the manifest file and at default locations
	icons, err := favicon.New().Find("https://www.deanishe.net")
	if err != nil {
		panic(err)
	}
	// icons are sorted widest first
	for _, i := range icons {
		fmt.Printf("%dx%d\t%s\n", i.Width, i.Height, i.FileExt)
	}
	// Output:
	// 256x256	png
	// 192x192	png
	// 180x180	png
	// 32x32	png
	// 16x16	png
	// 0x0	png
	// 0x0	ico
}

// Find favicons using custom options. Passing IgnoreManifest and IgnoreWellKnown
// causes the Finder to only retrieve the initial URL (HTML page).
func ExampleNew_withOptions() {
	f := favicon.New(
		// Don't look for or parse a manifest.json file
		favicon.IgnoreManifest,
		// Don't request files like /favicon.ico to see if they exist
		favicon.IgnoreWellKnown,
	)
	// Find icons defined in HTML, the manifest file and at default locations
	icons, err := f.Find("https://www.deanishe.net")
	if err != nil {
		panic(err)
	}
	// icons are sorted widest first
	for _, i := range icons {
		fmt.Printf("%dx%d\t%s\n", i.Width, i.Height, i.MimeType)
	}
	// Output:
	// 180x180	image/png
	// 32x32	image/png
	// 16x16	image/png
}
