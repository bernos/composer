package composer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestComposer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Composer Suite")
}
