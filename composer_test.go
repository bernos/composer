package composer

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetAttributes(t *testing.T) {
	z := html.NewTokenizer(strings.NewReader("<div name=\"brendan\" title=\"captain\"></div>"))
	z.Next()

	attributes := getAttributes(z, make([]*html.Attribute, 0))

	expected := 2

	if len(attributes) != expected {
		t.Errorf("Expected %d attributes. Found %d", expected, len(attributes))
	}

	if attributes[0].Key != "name" {
		t.Errorf("Expected attribute name brendan, but found %s", attributes[0].Key)
	}
}

func TestGetComposerTag(t *testing.T) {
	z1 := html.NewTokenizer(strings.NewReader("<div name=\"brendan\" title=\"captain\"></div>"))
	z2 := html.NewTokenizer(strings.NewReader("<div composer-url=\"http://example.com\" title=\"captain\"></div>"))

	z1.Next()
	z2.Next()

	t1 := getComposerTag(z1)
	t2 := getComposerTag(z2)

	if t1 != nil {
		t.Error("Expected no ComposerTag but found one")
	}

	if t2 == nil {
		t.Error("Did not find ComposerTag when one was expected")
	}

	if t2.Url != "http://example.com" {
		t.Error("Did not find correct ComposerTag url")
	}
}

func TestCompose(t *testing.T) {
	serverUrl := ""

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/a" {
			fmt.Fprintf(w, "Hello, client from <div composer-url=\"%s/b\"></div>", serverUrl)
		} else if r.URL.Path == "/b" {
			fmt.Fprint(w, "Hello, client from b")
		} else {
			fmt.Fprint(w, "Hello, "+r.URL.Path)
		}
	}))

	serverUrl = ts.URL

	defer ts.Close()

	unparsed := strings.NewReader(fmt.Sprintf("<div><p composer-url=\"%s/a\" title=\"captain\">other stuff</p></div>", ts.URL))
	parsed := Compose(unparsed)

	buf := new(bytes.Buffer)
	buf.ReadFrom(parsed)

	result := buf.String()

	if result != "blah" {
		t.Errorf("Expected blah but got %s", result)
	}
}

func TestPipeline(t *testing.T) {
	tags := []*ComposerTag{
		&ComposerTag{Url: "http://example/com/a"},
		&ComposerTag{Url: "http://example/com/b"},
		&ComposerTag{Url: "http://example/com/c"},
	}

	loader := func(url string) io.Reader {
		return strings.NewReader(url)
	}

	pipeline := BuildTagPipeline(tags, loader)

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
