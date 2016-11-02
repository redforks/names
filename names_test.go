package names_test

import (
	. "github.com/redforks/names"
	"github.com/redforks/testing/matcher"
	"github.com/redforks/testing/reset"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func wrapNextFunc(k Kind) func() (string, error) {
	return func() (string, error) {
		return Next(k)
	}
}

var _ = Describe("Names", func() {

	BeforeEach(func() {
		reset.Enable()
	})

	AfterEach(func() {
		reset.Disable()
	})

	DescribeTable("Tests", func(f func() (string, error)) {
		var names []string
		empties := 0
		for i := 0; i < 1002; i++ {
			var name string
			Ω(f()).Should(matcher.Save(&name))
			names = append(names, name)
			if name == "" {
				empties++
			}
		}
		Ω(empties).Should(BeNumerically("<", 50))
	},
		Entry("Person kind", wrapNextFunc(Person)),
		Entry("Product kind", wrapNextFunc(Product)),
		Entry("Address kind", wrapNextFunc(Address)),
		Entry("Firm kind", wrapNextFunc(Firm)),
		Entry("Person", NextPerson),
		Entry("Product", NextProduct),
		Entry("Address", NextAddress),
		Entry("Firm", NextFirm),
	)
})
