package s3_test

import (
	// TODO: Add this back when some tests are written:
	// . "github.com/mexisme/get-secrets/files/s3"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("The S3 type", func() {
	It("", func() {
		Skip("Placeholder")
		Expect(true).To(BeTrue())
	})

	It("panics with a bad Session", func() {
		Skip("Because the AWS SDK is bloody terrible to unit test.")
	})
})
