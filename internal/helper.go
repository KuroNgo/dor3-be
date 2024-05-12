package internal

import (
	"go.mongodb.org/mongo-driver/bson"
	"sync"
)

func ToDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}

var (
	Wg sync.WaitGroup
)
