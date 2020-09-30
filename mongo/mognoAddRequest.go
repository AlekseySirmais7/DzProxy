package mongo

import (
	"context"
	"sync"
)

var idMutex = sync.Mutex{}

func AddRequest(request MyRequest) error {

	// почему curl на localhost формируется иначе?
	if request.FullPath == "localhost:5001http://localhost:5001/testparam" {
		request.FullPath = "http://localhost:5001/testparam"
	}

	idMutex.Lock()
	request.Id = globalID
	globalID++
	idMutex.Unlock()

	_, err := collection.InsertOne(context.TODO(), request)
	if err != nil {
		return err
	}

	return nil
}
