package composer_test

import (
	"bytes"
	"io"
	"strings"

	. "github.com/bernos/composer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Composerpipeline", func() {
	var (
		tags     []*ComposerTag
		loader   func(string) io.Reader
		pipeline <-chan *ComposerTag
	)

	BeforeEach(func() {
		tags = []*ComposerTag{
			&ComposerTag{Url: "http://example/com/a"},
			&ComposerTag{Url: "http://example/com/b"},
			&ComposerTag{Url: "http://example/com/c"},
		}

		loader = func(url string) io.Reader {
			return strings.NewReader(url)
		}

		pipeline = BuildTagPipeline(tags, loader)
	})

	Describe("Pipeline", func() {
		It("Should emit loaded content", func() {
			for tag := range pipeline {
				buf := new(bytes.Buffer)
				buf.ReadFrom(tag.Content)
				Expect(tag.Url).To(Equal(buf.String()))
			}
		})
	})
})

/*

func TestPipeline(t *testing.T) {
	tags := []*ComposerTag{
		&ComposerTag{Url: "http://example/com/a"},
		&ComposerTag{Url: "http://example/com/b"},
		&ComposerTag{Url: "http://example/com/c"},
	}

	loader := func(url string) io.Reader {
		return strings.NewReader(url)
	}

	pipeline := buildTagPipeline(tags, loader)

	for tag := range pipeline {
		url := tag.Url
		buf := new(bytes.Buffer)
		buf.ReadFrom(tag.Content)
		content := buf.String()

		if url != content {
			t.Errorf("Expected %s but got %s", url, content)
		}
	}
}
*/
