package aws_test

import (
	// TODO: Add this back when some tests are written:
	// . "github.com/mexisme/get-secrets/config/aws"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("The config/aws type", func() {
	It("", func() {
		Skip("Placeholder")
		Expect(true).To(BeTrue())
	})

	It("panics with a bad Session", func() {
		Skip("Because the AWS SDK is bloody terrible to unit test.")
	})
})
