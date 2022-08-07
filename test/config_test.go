package test_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config Tests", Ordered, func() {
	It("config", func() {
		Expect(app.Run([]string{"barge", "config"})).To(Succeed())
	})
})
