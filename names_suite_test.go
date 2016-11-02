package names_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNames(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Names Suite")
}
