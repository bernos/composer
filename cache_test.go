package composer

import (
	"testing"
	"time"
)

func TestGetAndSet(t *testing.T) {
	expires := time.Now().Add(time.Minute)
	expected := "myvalue"
	key := "mykey"

	cache := NewMemoryCache()

	cache.Set(key, expected, expires)

	if actual, ok := cache.Get(key); !ok || actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestExpiry(t *testing.T) {
	expires := time.Now().Add(time.Second)
	key := "mykey"

	cache := NewMemoryCache()

	cache.Set(key, "myvalue", expires)

	time.Sleep(time.Second * 2)

	if actual, ok := cache.Get(key); ok {
		t.Errorf("Expected nil but got %s", actual)
	}
}
