package nested_document_array

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GenerateNestedDocsByCount(count int) {
	It("When size is"+strconv.Itoa(count)+" nested documents", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_, err := GenerateNestedDocumentArray(ctx, count)
		Expect(err).To(BeNil())

		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
		Expect(err).To(BeNil())
		database := client.Database(DocumentStoreModelingDatabaseName)
		result := database.RunCommand(context.Background(), bson.M{"collStats": NestedDocumentArrayCollectionName})
		var document bson.M
		err = result.Decode(&document)
		Expect(err).To(BeNil())

		fmt.Printf("Nested Doc Count: %d, Avg Object Size: %d bytes\n", count, document["avgObjSize"])
	})
}

var _ = Describe("Nested Document Array", func() {
	Context("Size of nested array documents", func() {
		nestedDocumentCounts := []int{10, 100, 1000, 10000}

		for _, count := range nestedDocumentCounts {
			GenerateNestedDocsByCount(count)
		}
	})
})
