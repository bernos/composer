package composer

import (
	. "github.com/onsi/ginkgo"
	"time"
)

var _ = Describe("Testing with Ginkgo", func() {
	It("get and set", func() {

		expires := time.Now().Add(time.Minute)
		expected := "myvalue"
		key := "mykey"

		cache := NewMemoryCache()

		cache.Set(key, expected, expires)

		if actual, ok := cache.Get(key); !ok || actual != expected {
			GinkgoT().Errorf("Expected %s but got %s", expected, actual)
		}
	})
	It("expiry", func() {

		expires := time.Now().Add(time.Second)
		key := "mykey"

		cache := NewMemoryCache()

		cache.Set(key, "myvalue", expires)

		time.Sleep(time.Second * 2)

		if actual, ok := cache.Get(key); ok {
			GinkgoT().Errorf("Expected nil but got %s", actual)
		}
	})
})
