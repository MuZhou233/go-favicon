// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-11-09

// Command favicon finds favicons for a URL and prints their format, size and URL.
// It can search HTML tags, manifest files and common locations on the server.
// Results are printed in a pretty table by default, but can be dumped as
// JSON, CSV or TSV.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	stdlog "log"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/muzhou233/go-favicon"
)

var (
	// set by Makefile.
	version   = "undefined"
	buildDate = "undefined"
)

var (
	fs          = flag.NewFlagSet("favicon", flag.ExitOnError)
	flagHelp    = fs.Bool("h", false, "show this message and exit")
	flagJSON    = fs.Bool("json", false, "output favicon list as JSON")
	flagCSV     = fs.Bool("csv", false, "output favicon list as CSV")
	flagTSV     = fs.Bool("tsv", false, "output favicon list as TSV")
	flagSquare  = fs.Bool("square", false, "only show square icons")
	flagVerbose = fs.Bool("v", false, "show informational messages")
	flagVersion = fs.Bool("version", false, "show version number and exit")

	log  *stdlog.Logger
	opts []favicon.Option
)

func usage() {
	fmt.Fprint(fs.Output(),
		`usage: favicon <url>

retrieve, favicons for URL

`)
	fs.PrintDefaults()
}

func init() {
	fs.Usage = usage
	log = stdlog.New(os.Stderr, "", 0)
}

func main() {
	checkErr(fs.Parse(os.Args[1:]))

	if *flagVersion {
		fmt.Printf("favicon %s\n", version)
		fmt.Printf("built: %s\n", buildDate)
		return
	}

	if *flagHelp || fs.NArg() == 0 {
		usage()
		return
	}

	u := fs.Arg(0)
	s := strings.ToLower(u)
	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		log.Fatalf("invalid URL: %q", s)
	}

	if *flagVerbose {
		opts = append(opts, favicon.WithLogger(log))
	}

	f := favicon.New(opts...)
	icons, err := f.Find(u)
	checkErr(err)
	// log.Printf("%d icon(s) found for %q", len(icons), u)

	if *flagSquare {
		var clean []*favicon.Icon
		for _, icon := range icons {
			if icon.IsSquare() {
				clean = append(clean, icon)
			}
		}
		icons = clean
	}

	if *flagJSON {
		data, err1 := json.MarshalIndent(icons, "", "  ")
		checkErr(err1)
		fmt.Println(string(data))
		return
	}

	if *flagTSV {
		fmt.Println("#\tformat\twidth\theight\tURL")
		for i, icon := range icons {
			fmt.Printf("%d\t%s\t%d\t%d\t%s\n", i+1, icon.MimeType, icon.Width, icon.Height, icon.URL)
		}
		return
	}

	t := table.NewWriter()
	t.AppendHeader(table.Row{
		"#", "format", "width", "height", "URL",
	})

	for i, icon := range icons {
		t.AppendRow(table.Row{
			i + 1, icon.MimeType, icon.Width, icon.Height, icon.URL,
		})
	}
	if *flagCSV {
		fmt.Println(t.RenderCSV())
		return
	}

	fmt.Println(t.Render())
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
