package composer

import (
	"flag"
	"io/ioutil"
	"net/http"

	"github.com/elazarl/goproxy"
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
		content := Compose(resp.Body)
		resp.Body.Close()
		resp.Body = ioutil.NopCloser(content)

		return resp
	})

	return &c
}
