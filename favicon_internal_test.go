package favicon

import (
	urls "net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type debugLogger struct {
	t *testing.T
}

func (l debugLogger) Printf(format string, v ...interface{}) {
	l.t.Logf(format, v...)
}

// TestFindManifest finds favicons in manifest.
func TestFindManifest(t *testing.T) {
	t.Parallel()
	file, err := os.Open("testdata/github/manifest.json")
	require.Nil(t, err, "unexpected error")
	defer file.Close()

	f := New(WithLogger(debugLogger{t}))
	require.Nil(t, err, "unexpected error")
	p := f.newParser()
	p.baseURL = mustURL("https://github.com")

	icons := p.parseManifestReader(file)
	assert.Equal(t, 11, len(icons), "unexpected favicon count")
}

// TestParserAbsURL tests resolution of URLs.
func TestParserAbsURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		in, x string
		base  *urls.URL
	}{
		{"empty", "", "", nil},
		{"onlyBaseURL", "", "", mustURL("https://github.com")},
		{"noBaseURL", "/root", "/root", nil},
		{"baseURL", "/root", "https://github.com/root", mustURL("https://github.com")},
		{"absURL", "https://github.com/root", "https://github.com/root", mustURL("https://github.com")},
		// absolute URLs returned as-is
		{"absURLDifferentBase", "https://github.com/root", "https://github.com/root", mustURL("https://google.com")},
		{"absURLNoBase", "https://github.com/root", "https://github.com/root", nil},
	}

	for _, td := range tests {
		td := td
		t.Run(td.name, func(t *testing.T) {
			t.Parallel()
			p := parser{baseURL: td.base}
			v := p.absURL(td.in)
			assert.Equal(t, td.x, v, "unexpected URL")
		})
	}
}

func mustURL(s string) *urls.URL {
	u, err := urls.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}
