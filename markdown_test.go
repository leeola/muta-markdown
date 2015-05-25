package markdown

import (
	"io/ioutil"
	"testing"

	"github.com/leeola/muta"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMarkdownStreamerNext(t *testing.T) {
	Convey("Should immediately return nil fi", t, func() {
		s := &MarkdownStreamer{Streamer: &muta.MockStreamer{}}
		fi, r, err := s.Next()
		So(fi, ShouldBeNil)
		So(r, ShouldBeNil)
		So(err, ShouldBeNil)
	})

	Convey("Should not create markdown on error", t, func() {
		s := &MarkdownStreamer{Streamer: &muta.MockStreamer{
			Files:    []string{"error.md"},
			Contents: []string{"error: **foo**"},
		}}
		_, r, err := s.Next()
		So(err, ShouldNotBeNil)
		b, err := ioutil.ReadAll(r)
		So(string(b), ShouldEqual, "error: **foo**")
	})

	Convey("Should not modify non-markdown files", t, func() {
		s := &MarkdownStreamer{Streamer: &muta.MockStreamer{
			Files:    []string{"file"},
			Contents: []string{"This **isn't** a markdown file!"},
		}}
		_, r, err := s.Next()
		So(err, ShouldBeNil)
		b, err := ioutil.ReadAll(r)
		So(string(b), ShouldEqual, "This **isn't** a markdown file!")
	})

	Convey("Should compile markdown files to html", t, func() {
		s := &MarkdownStreamer{Streamer: &muta.MockStreamer{
			Files:    []string{"file.md"},
			Contents: []string{"This **is** a markdown file!"},
		}}
		_, r, err := s.Next()
		So(err, ShouldBeNil)
		b, err := ioutil.ReadAll(r)
		So(string(b), ShouldEqual,
			"<p>This <strong>is</strong> a markdown file!</p>\n")
	})

	Convey("Should rename markdown files to html", t, func() {
		s := &MarkdownStreamer{Streamer: &muta.MockStreamer{
			Files: []string{"file.md"},
		}}
		fi, _, err := s.Next()
		So(err, ShouldBeNil)
		So(fi.Name, ShouldEqual, "file.html")
	})
}
