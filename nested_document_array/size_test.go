package nested_document_array

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SampleNestedDocument struct {
	Name    string
	Phone   string
	Address string
}

type SampleDocument struct {
	Id              primitive.ObjectID `bson:"_id"`
	NestedDocuments []*SampleNestedDocument
}

func rangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}

var (
	DocumentStoreModelingDatabaseName = "document_store_modeling"
	NestedDocumentArrayCollectionName = "nested_document_array"
)

func GenerateNestedDocsByCount(count int) {
	It("When size is"+strconv.Itoa(count)+" nested documents", func() {
		ctx := context.TODO()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
		Expect(err).To(BeNil())

		database := client.Database(DocumentStoreModelingDatabaseName)
		collection := database.Collection(NestedDocumentArrayCollectionName)
		err = collection.Drop(ctx)
		Expect(err).To(BeNil())

		sampleDocument := &SampleDocument{}
		for i := 1; i <= count; i++ {
			sampleNestedDocument := &SampleNestedDocument{
				Name:    "Test Name " + strconv.Itoa(i),
				Phone:   strconv.Itoa(rangeIn(1000000000, 9999999999)),
				Address: strconv.Itoa(rangeIn(1000, 9999)) + " Street",
			}
			sampleDocument.NestedDocuments = append(sampleDocument.NestedDocuments, sampleNestedDocument)
		}

		_, err = collection.InsertOne(ctx, sampleDocument)

		Expect(err).To(BeNil())

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
