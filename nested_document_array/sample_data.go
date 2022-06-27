package nested_document_array

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"strconv"
)

var (
	DocumentStoreModelingDatabaseName = "document_store_modeling"
	NestedDocumentArrayCollectionName = "nested_document_array"
)

type SampleNestedDocument struct {
	Name    string `bson:"Name"`
	Phone   string `bson:"Phone"`
	Address string `bson:"Address"`
}

type SampleDocument struct {
	Id              primitive.ObjectID      `bson:"_id"`
	NestedDocuments []*SampleNestedDocument `bson:"NestedDocuments"`
}

func rangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}

func GenerateSampleNestedArrayDocumentCollection(ctx context.Context, nestedDocumentCount int) (result *mongo.InsertOneResult, collection *mongo.Collection, err error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	database := client.Database(DocumentStoreModelingDatabaseName)
	collection = database.Collection(NestedDocumentArrayCollectionName)
	err = collection.Drop(ctx)
	if err != nil {
		return result, collection, err
	}

	sampleDocument := &SampleDocument{}
	for i := 1; i <= nestedDocumentCount; i++ {
		sampleNestedDocument := &SampleNestedDocument{
			Name:    "Test Name " + strconv.Itoa(i),
			Phone:   strconv.Itoa(rangeIn(1000000000, 9999999999)),
			Address: strconv.Itoa(rangeIn(1000, 9999)) + " Street",
		}
		sampleDocument.NestedDocuments = append(sampleDocument.NestedDocuments, sampleNestedDocument)
	}

	result, err = collection.InsertOne(ctx, sampleDocument)
	return result, collection, err
}
