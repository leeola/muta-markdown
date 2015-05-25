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

func Markdown() muta.StreamEmbedder {
	opts := Options{}
	return MarkdownOpts(opts)
}

func MarkdownOpts(opts Options) muta.StreamEmbedder {
	return muta.StreamEmbedderFunc(func(inner muta.Streamer) muta.Streamer {
		return &MarkdownStreamer{Streamer: inner, Opts: opts}
	})
}

type Options struct {
}

type MarkdownStreamer struct {
	muta.Streamer
	Opts Options
}

func (s *MarkdownStreamer) Use(embedder muta.StreamEmbedder) muta.Streamer {
	return embedder.Embed(s)
}

func (s *MarkdownStreamer) Next() (*muta.FileInfo, io.ReadCloser, error) {
	// We don't generate files, so no need to ever do anything if we don't
	// have an inner Streamer.
	if s.Streamer == nil {
		return nil, nil, nil
	}

	fi, r, err := s.Streamer.Next()
	if fi == nil || err != nil {
		return fi, r, err
	}

	// If the file isn't markdown, we don't care about it. Return it
	// unmodified.
	if filepath.Ext(fi.Name) != ".md" {
		return fi, r, err
	}

	// Rename the file to HTML
	fi.Name = fmt.Sprintf("%s.html",
		strings.TrimSuffix(fi.Name, filepath.Ext(fi.Name)))

	// Since the file is markdown, read it all so we can convert it to
	// markdown.
	markdown, err := ioutil.ReadAll(r)
	defer r.Close()
	if err != nil {
		return fi, r, err
	}

	html := blackfriday.MarkdownBasic(markdown)

	// Now return it all, with a ReadCloser for the html.
	return fi, mutil.ByteCloser(html), nil
}
