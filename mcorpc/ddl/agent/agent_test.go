package agent

import (
	"path"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFileContent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "McoRPC/DDL/Agent")
}

var _ = Describe("McoRPC/DDL/Agent", func() {
	var pkg *DDL
	var err error

	BeforeEach(func() {
		pkg, err = New(path.Join("testdata", "package.json"))
		Expect(err).ToNot(HaveOccurred())
	})

	var _ = Describe("New", func() {
		It("Should fail for missing json files", func() {
			d, err := New(path.Join("testdata", "missing.json"))
			Expect(err.Error()).To(MatchRegexp("could not load DDL data: open.+missing.json"))
			Expect(d).To(BeNil())
		})

		It("Should fail for invalid json files", func() {
			d, err := New(path.Join("testdata", "invalid.json"))
			Expect(err).To(MatchError("could not parse JSON data in testdata/invalid.json: unexpected end of JSON input"))
			Expect(d).To(BeNil())
		})

		It("Should correctly load valid DDL files", func() {
			Expect(pkg.Metadata.Author).To(Equal("R.I.Pienaar <rip@devco.net>"))
			Expect(pkg.Metadata.Description).To(Equal("Manage Operating System Packages"))
			Expect(pkg.Metadata.License).To(Equal("Apache-2.0"))
			Expect(pkg.Metadata.Name).To(Equal("package"))
			Expect(pkg.Metadata.Timeout).To(Equal(180))
			Expect(pkg.Metadata.URL).To(Equal("https://github.com/choria-plugins/package-agent"))
			Expect(pkg.Metadata.Version).To(Equal("5.0.0"))
		})
	})

	var _ = Describe("ActionList", func() {
		It("Should return the correct list", func() {
			Expect(pkg.ActionNames()).To(Equal([]string{"apt_checkupdates", "apt_update", "checkupdates", "count", "install", "md5", "purge", "status", "uninstall", "update", "yum_checkupdates", "yum_clean"}))
		})
	})

	var _ = Describe("Timeout", func() {
		It("Should handle 0 second timeouts as 10 seconds", func() {
			pkg.Metadata.Timeout = 0

			Expect(pkg.Timeout()).To(Equal(time.Duration(10 * time.Second)))
		})

		It("Should handle timeouts correctly", func() {
			Expect(pkg.Timeout()).To(Equal(time.Duration(180 * time.Second)))
		})
	})

	var _ = Describe("ActionInterface", func() {
		It("Should retrieve the correct interface", func() {
			act, err := pkg.ActionInterface("install")
			Expect(err).ToNot(HaveOccurred())

			Expect(act.Name).To(Equal("install"))
			Expect(act.Description).To(Equal("Install a package"))
			Expect(act.Display).To(Equal("failed"))
			Expect(act.Output).To(HaveLen(8))
		})

		It("Should handle unknown interfaces", func() {
			act, err := pkg.ActionInterface("unknown")
			Expect(err).To(HaveOccurred())
			Expect(act).To(BeNil())
		})
	})
})
