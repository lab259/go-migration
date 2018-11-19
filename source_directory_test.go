package migration_test

import (
	"github.com/lab259/go-migration"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Source Directory", func() {
	Describe("List", func() {
		It("should list the files of a directory", func() {
			d := migration.DirectorySource{
				Directory: "test/migrations1",
				Extension: "sql",
			}
			ms, err := d.List()
			Expect(err).To(BeNil())
			Expect(ms).To(HaveLen(2))
			Expect(ms[0].GetID()).To(Equal(time.Date(2017, 10, 25, 19, 17, 47, 0, time.UTC)))
			Expect(ms[0].GetDescription()).To(Equal("description1"))

			Expect(ms[1].GetID()).To(Equal(time.Date(2017, 10, 25, 21, 33, 03, 0, time.UTC)))
			Expect(ms[1].GetDescription()).To(Equal("description2"))
		})
	})
})
