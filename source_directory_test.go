package migration

import (
	"time"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDirectorySource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Source Directory")
}

var _ = Describe("Source Directory", func() {
	Describe("File pattern", func() {
		It("It should match a file", func() {
			ms := DirectorySourcePattern.FindStringSubmatch("20010203040506_Description1")
			Expect(ms).To(HaveLen(3))
			Expect(ms[1]).To(Equal("20010203040506"))
			Expect(ms[2]).To(Equal("Description1"))
		})
	})

	Describe("List", func() {
		It("It should list the files of a directory", func() {
			d := DirectorySource{
				Directory: "test/migrations1",
				Extension: "sql",
			}
			ms, err := d.List()
			Expect(err).To(BeNil())
			Expect(ms).To(HaveLen(2))
			Expect(ms[0].GetId()).To(Equal(time.Date(2017, 10, 25, 19, 17, 47, 0, time.UTC)))
			Expect(ms[0].GetDescription()).To(Equal("description1"))

			Expect(ms[1].GetId()).To(Equal(time.Date(2017, 10, 25, 21, 33, 03, 0, time.UTC)))
			Expect(ms[1].GetDescription()).To(Equal("description2"))
		})
	})
})
