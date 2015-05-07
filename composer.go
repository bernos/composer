package composer

import (
	"bytes"
	"flag"
	"github.com/elazarl/goproxy"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	//"strings"
)

var (
	listen = flag.String("listen", "localhost:1080", "listen on address")
)

type ComposerHttpServer struct {
	proxy *goproxy.ProxyHttpServer
}

func (c *ComposerHttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Scheme = "http"
	r.Host = "localhost:8080"
	r.URL.Host = "localhost:8080"

	c.proxy.ServeHTTP(w, r)
}

func NewComposerHttpServer() *ComposerHttpServer {
	c := ComposerHttpServer{
		proxy: goproxy.NewProxyHttpServer(),
	}

	c.proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		content := parseHtml(resp.Body)
		resp.Body.Close()
		resp.Body = ioutil.NopCloser(content)

		return resp
	})

	return &c
}

type ComposerTag struct {
	URL string
}

func parseHtml(r io.Reader) io.Reader {
	z := html.NewTokenizer(r)
	buf := new(bytes.Buffer)

	for {
		tt := z.Next()
		log.Printf("tt: %v", tt)

		if tt == html.ErrorToken {
			return buf
		}

		if _, hasAttr := z.TagName(); hasAttr {
			attributes := getAttributes(z, make([]*html.Attribute, 0))

			for i := range attributes {
				log.Printf("Attr %s=%s", attributes[i].Key, attributes[i].Val)
			}
		}

		buf.Write(z.Raw())
		log.Print(string(z.Raw()))
	}

}

func getComposerTag(z *html.Tokenizer) *ComposerTag {
	attributes := getAttributes(z, make([]*html.Attribute, 0))

	composerTag := new(ComposerTag)

	for i := range attributes {
		if attributes[i].Key == "composer-url" {
			composerTag.URL = attributes[i].Val
		}
	}

	if len(composerTag.URL) > 0 {
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
