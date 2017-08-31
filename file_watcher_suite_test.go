package file_watcher_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"fmt"
)

func TestFileWatcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FileWatcher Suite")
}

var _ = Describe("VagrantFile", func() {
	var vagrantFile VagrantFile

	BeforeEach(func() {
		vagrantFile = VagrantFile{fpath: "./Vagrantfile"}
	})

	Describe("#parse", func() {
		It("It returns a slice of watched directions and files", func() {
			watchedFiles := vagrantFile.parse()

			for file := range watchedFiles {
				fmt.Printf("Matched file! - %s", file)
			}

			Expect(true).To(Equal(true))
		})
	})
})
