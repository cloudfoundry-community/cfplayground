package tutorials_test

import (
	"io/ioutil"
	"strings"

	. "github.com/simonleung8/cfplayground/tutorials"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tutorials", func() {
	tutorial := New("./courses/")

	Describe("New()", func() {
		It("reads each JSON file in directory courses/ and create a new course in map", func() {
			numFiles := 0

			list, _ := ioutil.ReadDir("./courses")
			for _, info := range list {
				if strings.HasSuffix(info.Name(), ".json") {
					numFiles++
				}
			}
			Ω(tutorial.TotalCourses()).To(Equal(numFiles))
		})
	})

	Describe("StartCourse()", func() {
		var (
			instruct string
			step     string
		)

		Context("If no tutorial is in progress", func() {
			It("should report tutorial not in progress", func() {
				Ω(tutorial.InProgress()).To(BeFalse())
			})
		})

		Context("Starting a tutorial", func() {
			It("should report tutorial is in progress", func() {
				instruct, step = tutorial.StartCourse("Listing applications")
				Ω(tutorial.InProgress()).To(BeTrue())
			})

			It("returns the a instruction set", func() {
				Ω(instruct).ToNot(Equal(""))
			})

			It("returns step index", func() {
				Ω(step).To(ContainSubstring("1/"))
			})
		})

		PContext("Advancing to the next step in the tutorial", func() {

		})
	})

})
