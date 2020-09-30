package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MyRequest struct {
	Id       int
	Method   string
	FullPath string // with get params
	Headers  map[string][]string
	Body     []byte
}

var globalID = 0

func GetRequests() ([]MyRequest, error) {

	// find all
	var results []MyRequest

	// Finding multiple documents returns a cursor
	findOptions := options.Find()
	findOptions.SetLimit(500)

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return results, err
	}

	// Iterate through the cursor
	for cur.Next(context.TODO()) {
		var elem MyRequest
		err := cur.Decode(&elem)
		if err != nil {
			return results, err
		}

		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		return results, err
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	//fmt.Printf("Found multiple document: %+v\n", results)

	return results, nil

}
