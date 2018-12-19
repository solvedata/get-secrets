package s3url_test

import (
	// "net/url"
	. "github.com/mexisme/get-secrets/files/s3/s3url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("The S3 URL type", func() {
	It("panics with an empty URL", func() {
		Expect(func() {
			New().WithURL("")
		}).To(Panic())
	})

	It("panics with an invalid URL", func() {
		Expect(func() {
			New().WithURL("..why::")
		}).To(Panic())
	})

	It("returns correct Bucket and Prefix from a good Base URL", func() {
		s3ish := New().WithURL("s3://a/b")

		Expect(s3ish.Bucket()).To(Equal("a"))
		Expect(s3ish.Prefix()).To(Equal("b"))
	})

	It("returns the correct Bucket when given a Bucket", func() {
		Expect(New().WithBucket("Yerp").Bucket()).To(Equal("Yerp"))
	})

	It("panics with an empty Bucket", func() {
		Expect(func() {
			New().WithBucket("")
		}).To(Panic())
	})

	It("returns the correct Prefix when given a Prefix", func() {
		s3ish := New().WithBucketPrefix("Yerp", "Derp")
		Expect(s3ish.Bucket()).To(Equal("Yerp"))
		Expect(s3ish.Prefix()).To(Equal("Derp"))
	})

	It("returns a '/' Prefix when not given a Prefix", func() {
		Expect(New().WithBucket("Yerp").Prefix()).To(Equal("/"))
	})

	It("returns a '/' Prefix when given an empty Prefix", func() {
		Expect(New().WithBucketPrefix("Yerp", "").Prefix()).To(Equal("/"))
	})

	It("trims left-most '/' from Prefix", func() {
		Expect(New().WithBucketPrefix("Yerp", "//hello").Prefix()).To(Equal("hello"))
	})

	It("trims right-most '/' from Prefix", func() {
		Expect(New().WithBucketPrefix("Yerp", "hello//").Prefix()).To(Equal("hello"))
	})

	It("trims both left-most and right-most '/' from Prefix", func() {
		Expect(New().WithBucketPrefix("Yerp", "//hello/there/you//").Prefix()).To(Equal("hello/there/you"))
	})

	It("panics with an empty Bucket and 'full' Prefix", func() {
		Expect(func() {
			New().WithBucketPrefix("", "Derp")
		}).To(Panic())
	})
})
