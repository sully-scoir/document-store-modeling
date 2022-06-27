package nested_document_array

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var _ = Describe("Nested Document Array", func() {
	Context("Query by Name", func() {
		It("Client-side query", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			result, err := GenerateNestedDocumentArray(ctx, 100)
			Expect(err).To(BeNil())

			client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
			Expect(err).To(BeNil())
			collection := client.
				Database(DocumentStoreModelingDatabaseName).
				Collection(NestedDocumentArrayCollectionName)

			findOneResult := collection.FindOne(ctx, bson.M{"_id": result.InsertedID})
			Expect(findOneResult.Err()).To(BeNil())

			foundSampleDocument := &SampleDocument{}
			err = findOneResult.Decode(foundSampleDocument)
			Expect(err).To(BeNil())

			sampleNameQuery := "Test Name 1"
			matchResults := []*SampleNestedDocument{}
			for _, nestedDocument := range foundSampleDocument.NestedDocuments {
				if nestedDocument.Name == sampleNameQuery {
					matchResults = append(matchResults, nestedDocument)
				}
			}

			Expect(len(matchResults)).To(Equal(1))
			Expect(matchResults[0].Name).To(Equal(sampleNameQuery))
		})
	})
})
