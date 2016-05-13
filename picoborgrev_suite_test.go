package picoborgrev_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPicoborgrev(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Picoborgrev Suite")
}
