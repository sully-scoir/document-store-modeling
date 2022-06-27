package nested_document_array_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNestedDocumentArray(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NestedDocumentArray Suite")
}
