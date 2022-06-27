package nested_document_array

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

var _ = Describe("Nested Document Array", func() {
	FContext("Query by Name", func() {
		It("Client-side query", func() {
			ctx := context.TODO()
			client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
			Expect(err).To(BeNil())

			database := client.Database(DocumentStoreModelingDatabaseName)
			collection := database.Collection(NestedDocumentArrayCollectionName)
			err = collection.Drop(ctx)
			Expect(err).To(BeNil())

			sampleDocument := &SampleDocument{}
			count := 100
			for i := 1; i <= count; i++ {
				sampleNestedDocument := &SampleNestedDocument{
					Name:    "Test Name " + strconv.Itoa(i),
					Phone:   strconv.Itoa(rangeIn(1000000000, 9999999999)),
					Address: strconv.Itoa(rangeIn(1000, 9999)) + " Street",
				}
				sampleDocument.NestedDocuments = append(sampleDocument.NestedDocuments, sampleNestedDocument)
			}

			result, err := collection.InsertOne(ctx, sampleDocument)
			Expect(err).To(BeNil())

			// Query by Name
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
