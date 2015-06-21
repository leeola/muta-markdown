package markdown

import (
	"io/ioutil"
	"testing"

	"github.com/leeola/muta"
	"github.com/leeola/muta/mutil"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMarkdownStreamerNext(t *testing.T) {
	Convey("Should not attempt to modify nil fi", t, func() {
		s := &MarkdownStreamer{}
		fi, rc, err := s.Next(nil, nil)
		So(fi, ShouldBeNil)
		So(rc, ShouldBeNil)
		So(err, ShouldBeNil)
	})

	Convey("Should not modify non-markdown files", t, func() {
		s := &MarkdownStreamer{}
		_, rc, err := s.Next(
			muta.NewFileInfo("file.txt"),
			mutil.StringCloser("This **isn't** a markdown file!"),
		)
		So(err, ShouldBeNil)
		b, err := ioutil.ReadAll(rc)
		So(string(b), ShouldEqual, "This **isn't** a markdown file!")
	})

	Convey("Should compile markdown files to html", t, func() {
		s := &MarkdownStreamer{}
		_, rc, err := s.Next(
			muta.NewFileInfo("file.md"),
			mutil.StringCloser("This **is** a markdown file!"),
		)
		So(err, ShouldBeNil)
		b, err := ioutil.ReadAll(rc)
		So(string(b), ShouldEqual,
			"<p>This <strong>is</strong> a markdown file!</p>\n")
	})

	Convey("Should rename markdown files to html", t, func() {
		s := &MarkdownStreamer{}
		fi, _, err := s.Next(
			muta.NewFileInfo("file.md"),
			mutil.StringCloser("foo"),
		)
		So(err, ShouldBeNil)
		So(fi.Name(), ShouldEqual, "file.html")
	})
}
