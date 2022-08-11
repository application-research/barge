package tests_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Login Tests", Ordered, func() {

	It("login", func() {
		Expect(app.Run([]string{"barge", "login", os.Getenv("ESTUARY_API_KEY")})).To(Succeed())
	})
})
