package composer

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type ComposerTag struct {
	Content io.Reader
	Offset  int
	Url     string
}

func Compose(r io.Reader) io.Reader {
	buf, tags := extractComposerTags(r)
	output := new(bytes.Buffer)

	if len(tags) > 0 {
		offset := 0

		for tag := range buildTagPipeline(tags, fetchFromUrl) {
			output.Write(buf.Next(tag.Offset - offset))

			if tag.Content != nil {
				output.ReadFrom(Compose(tag.Content))
			} else {
				output.WriteString("<!-- composer failed to load content from " + tag.Url + " -->")
			}

			offset += tag.Offset
		}
	}

	output.Write(buf.Next(buf.Len()))

	return output
}

/**
 * extractComposerTags will read html content from an io.Reader and extract any
 * composer tags, creating a ComposerTag for each. It will also append all html
 * content outside of the composer tags to to a buffer, ready to merged with
 * the tags
 */
func extractComposerTags(r io.Reader) (*bytes.Buffer, []*ComposerTag) {
	z := html.NewTokenizer(r)
	buf := new(bytes.Buffer)
	depth := 0
	tags := make([]*ComposerTag, 0)

	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			return buf, tags
		}

		if depth == 0 {
			if tag := getComposerTag(z); tag != nil {
				tag.Offset = buf.Len()
				tags = append(tags, tag)
				if tt == html.StartTagToken {
					depth++
				}
			} else {
				buf.Write(z.Raw())
			}
		} else if tt == html.EndTagToken {
			depth--
		}
	}
}

func getComposerTag(z *html.Tokenizer) *ComposerTag {

	attributes := getAttributes(z, make([]*html.Attribute, 0))

	composerTag := new(ComposerTag)

	for i := range attributes {
		if attributes[i].Key == "composer-url" {
			composerTag.Url = attributes[i].Val
		}
	}

	if len(composerTag.Url) > 0 {
		return composerTag
	}

	return nil
}

func getAttributes(z *html.Tokenizer, a []*html.Attribute) []*html.Attribute {
	key, val, more := z.TagAttr()

	a = append(a, &html.Attribute{Key: string(key), Val: string(val)})

	if more {
		return getAttributes(z, a)
	}

	return a
}

/**
 * fetchFromUrl will fetch content for the url and return the resulting io.Reader
 */
func fetchFromUrl(url string) io.Reader {
	ch := make(chan io.Reader, 1)

	go func(url string) {

		res, err := http.Get(url)

		if err != nil {
			ch <- strings.NewReader(fmt.Sprintf("<-- composer encountered an error while loading content from %s - %v", url, err))
		} else {
			greeting, err := ioutil.ReadAll(res.Body)
			res.Body.Close()

			if err != nil {
				ch <- strings.NewReader(fmt.Sprintf("<-- composer encountered an error while loading content from %s - %v", url, err))
			} else {
				ch <- bytes.NewReader(greeting)
			}
		}
	}(url)

	select {
	case r := <-ch:
		return r
	case <-time.After(time.Second * 1):
		return nil
	}
}
