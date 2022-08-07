package test_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Init Tests", Ordered, func() {

	It("init", func() {
		Expect(app.Run([]string{"barge", "init"})).To(Succeed())
	})
})
