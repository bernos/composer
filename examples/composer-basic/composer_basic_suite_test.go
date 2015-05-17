package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestComposerBasic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ComposerBasic Suite")
}
