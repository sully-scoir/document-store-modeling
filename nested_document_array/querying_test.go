package nested_document_array

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var _ = Describe("Nested Document Array", func() {
	var startTime time.Time
	BeforeEach(func() {
		startTime = time.Now()
	})
	AfterEach(func() {
		fmt.Println("Runtime:", time.Since(startTime))
	})

	Context("Query by Name", func() {
		It("Client-side query given document Id", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			result, collection, err := GenerateSampleNestedArrayDocumentCollection(ctx, 100)
			Expect(err).To(BeNil())

			findOneResult := collection.FindOne(ctx, bson.M{"_id": result.InsertedID})
			Expect(findOneResult.Err()).To(BeNil())

			foundSampleDocument := &SampleDocument{}
			err = findOneResult.Decode(foundSampleDocument)
			Expect(err).To(BeNil())

			sampleNameQuery := "Test Name 1"
			clientSideFindResult := []*SampleNestedDocument{}
			for _, nestedDocument := range foundSampleDocument.NestedDocuments {
				if nestedDocument.Name == sampleNameQuery {
					clientSideFindResult = append(clientSideFindResult, nestedDocument)
				}
			}

			Expect(len(clientSideFindResult)).To(Equal(1))
			Expect(clientSideFindResult[0].Name).To(Equal(sampleNameQuery))
		})

		It("Client-side query matching array element", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			_, collection, err := GenerateSampleNestedArrayDocumentCollection(ctx, 100)
			Expect(err).To(BeNil())

			// First, look up the document that has a nested document array
			// with an element that has Name matching sampleNameQuery
			sampleNameQuery := "Test Name 1"
			findCursor, err := collection.Find(ctx, bson.M{"NestedDocuments.Name": sampleNameQuery})
			defer findCursor.Close(ctx)
			Expect(err).To(BeNil())

			findCursor.Next(ctx)
			foundSampleDocument := &SampleDocument{}
			err = findCursor.Decode(foundSampleDocument)
			Expect(err).To(BeNil())

			clientSideMatchResults := []*SampleNestedDocument{}
			for _, nestedDocument := range foundSampleDocument.NestedDocuments {
				if nestedDocument.Name == sampleNameQuery {
					clientSideMatchResults = append(clientSideMatchResults, nestedDocument)
				}
			}

			Expect(len(clientSideMatchResults)).To(Equal(1))
			Expect(clientSideMatchResults[0].Name).To(Equal(sampleNameQuery))
		})

		It("Aggregate returning single matching nested array document", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			result, collection, err := GenerateSampleNestedArrayDocumentCollection(ctx, 100)
			Expect(err).To(BeNil())

			sampleNameQuery := "Test Name 1"
			sampleNameMatch := bson.D{
				{"$match", bson.M{"NestedDocuments.Name": sampleNameQuery}},
			}

			aggregateCursor, err := collection.Aggregate(ctx, mongo.Pipeline{
				// First match returns the parent document with the array of documents to query
				bson.D{
					{"$match", bson.M{"_id": result.InsertedID}},
				},
				// Unwind lifts all the array's documents out so that they can be queried
				bson.D{
					{"$unwind", "$NestedDocuments"},
				},
				// Now we query the lifted-out array nested documents
				sampleNameMatch,
				// We want our result to be the matching nested document itself
				bson.D{
					{"$replaceRoot", bson.M{"newRoot": "$NestedDocuments"}},
				},
			})
			Expect(err).To(BeNil())

			aggregateCursor.Next(ctx)
			foundSampleNestedDocument := &SampleNestedDocument{}
			err = aggregateCursor.Decode(foundSampleNestedDocument)
			Expect(err).To(BeNil())
			Expect(aggregateCursor.Next(ctx)).To(BeFalse(), "Cursor should only find one match")

			Expect(foundSampleNestedDocument.Name).To(Equal(sampleNameQuery))
		})

		It("Aggregate returning multiple matching nested array documents", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			result, collection, err := GenerateSampleNestedArrayDocumentCollection(ctx, 100)
			Expect(err).To(BeNil())

			sampleNameQuery := "Test Name 1"
			sampleNameQuery2 := "Test Name 2"
			sampleNameMatch := bson.D{
				{"$match",
					bson.M{"$or": []bson.M{
						{"NestedDocuments.Name": sampleNameQuery},
						{"NestedDocuments.Name": sampleNameQuery2},
					}},
				},
			}

			aggregateCursor, err := collection.Aggregate(ctx, mongo.Pipeline{
				// First match returns the parent document with the array of documents to query
				bson.D{
					{"$match", bson.M{"_id": result.InsertedID}},
				},
				// Unwind lifts all the array's documents out so that they can be queried
				bson.D{
					{"$unwind", "$NestedDocuments"},
				},
				// Now we query the lifted-out array nested documents
				sampleNameMatch,
				// We want our result to be the matching nested document itself
				bson.D{
					{"$replaceRoot", bson.M{"newRoot": "$NestedDocuments"}},
				},
			})
			Expect(err).To(BeNil())

			foundSampleNestedDocuments := []*SampleNestedDocument{}
			for aggregateCursor.Next(ctx) {
				sampleNestedDocument := &SampleNestedDocument{}
				err = aggregateCursor.Decode(sampleNestedDocument)
				Expect(err).To(BeNil())
				foundSampleNestedDocuments = append(foundSampleNestedDocuments, sampleNestedDocument)
			}

			Expect(len(foundSampleNestedDocuments)).To(Equal(2))
			Expect(foundSampleNestedDocuments[0].Name).To(Equal(sampleNameQuery))
			Expect(foundSampleNestedDocuments[1].Name).To(Equal(sampleNameQuery2))
		})
	})
})
