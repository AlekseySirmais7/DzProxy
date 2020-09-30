package mongo

import (
	"context"
	"gopkg.in/mgo.v2/bson"
	"log"
)

func GetOneRequest(id int) (MyRequest, error) {
	var request MyRequest
	log.Println("GetOneRequest id: ", id)
	err := collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&request)
	if err != nil {
		panic(err)
	}
	return request, err
}
