package markdown

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/leeola/muta"
	"github.com/leeola/muta/mutil"
	"github.com/russross/blackfriday"
)

const pluginName string = "muta-markdown"

func Markdown() muta.Streamer {
	opts := Options{}
	return MarkdownOpts(opts)
}

func MarkdownOpts(opts Options) muta.Streamer {
	return &MarkdownStreamer{Opts: opts}
}

type Options struct {
}

type MarkdownStreamer struct {
	Opts Options
}

func (s *MarkdownStreamer) Next(fi muta.FileInfo, rc io.ReadCloser) (
	muta.FileInfo, io.ReadCloser, error) {

	// MarkdownStreamer does not create any files, so if no files are
	// given to it, just return.
	if fi == nil {
		return fi, rc, nil
	}

	// If the file isn't markdown, we don't care about it. Return it
	// unmodified.
	if filepath.Ext(fi.Name()) != ".md" {
		return fi, rc, nil
	}

	// Rename the file to HTML
	fi.SetName(fmt.Sprintf("%s.html",
		strings.TrimSuffix(fi.Name(), filepath.Ext(fi.Name())),
	))

	// Since the file is markdown, read it all so we can convert it to
	// markdown.
	markdown, err := ioutil.ReadAll(rc)
	defer rc.Close()
	if err != nil {
		return fi, rc, err
	}

	html := blackfriday.MarkdownBasic(markdown)

	// Now return it all, with a ReadCloser for the html.
	return fi, mutil.ByteCloser(html), nil
}
