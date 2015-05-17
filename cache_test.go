package composer

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Memory Cache", func() {

	Context("When getting a value from the cache", func() {

		var (
			cache Cache
		)

		BeforeEach(func() {
			cache = NewMemoryCache()

			cache.Set("key_one", "value_one", time.Now().Add(time.Minute))
			cache.Set("key_two", "value_two", time.Now().Add(time.Minute*-1))
		})

		It("Should return the value and true as the status if the value exists and is not expired", func() {
			value, ok := cache.Get("key_one")

			Expect(value).To(Equal("value_one"))
			Expect(ok).To(BeTrue())
		})

		It("Should return false status if the value does not exist", func() {
			value, ok := cache.Get("key_three")

			Expect(value).To(BeZero())
			Expect(ok).To(BeFalse())
		})

		It("Should return false status if the value has expired", func() {
			value, ok := cache.Get("key_two")

			Expect(value).To(BeZero())
			Expect(ok).To(BeFalse())
		})
	})

	Context("When setting a value in the cache", func() {
		var (
			cache Cache
		)

		BeforeEach(func() {
			cache = NewMemoryCache()
		})

		It("Should update the value if a new value with the same key is set", func() {
			cache.Set("key_one", "value_one", time.Now().Add(time.Minute))
			cache.Set("key_one", "value_two", time.Now().Add(time.Minute))

			value, ok := cache.Get("key_one")

			Expect(value).To(Equal("value_two"))
			Expect(ok).To(BeTrue())
		})
	})
})
