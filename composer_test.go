package composer

import (
	"golang.org/x/net/html"
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

	if t2.URL != "http://example.com" {
		t.Error("Did not find correct ComposerTag url")
	}
}
