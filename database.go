package mongodbWrapper

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type Database interface {
	AddCollections(...string)
	checkCollections(string) error
	InsertOne(string, interface{}) (*mongo.InsertOneResult, error)
	InsertMany(string, []interface{}) (*mongo.InsertManyResult, error)
	DeleteOne(string, interface{}) (*mongo.DeleteResult, error)
	DeleteMany(string, []interface{}) (*mongo.DeleteResult, error)
	FindOne(string, interface{}) (*mongo.SingleResult, error)
	FindMany(string, interface{}) (*mongo.Cursor, error)
	UpdateOne(string, interface{}, interface{}) (*mongo.UpdateResult, error)
	UpdateMany(string, interface{}, interface{}) (*mongo.UpdateResult, error)
}

type DatabaseInfo struct {
	name        string
	collections map[string]*mongo.Collection
	Ctx         context.Context
	database    *mongo.Database
}

func (v *DatabaseInfo) checkCollections(collections string) error {
	_, ok := v.collections[collections]
	if !ok {
		return fmt.Errorf("Collection %s does not exist", collections)
	}
	return nil
}

func (v *DatabaseInfo) InsertOne(collections string, document interface{}) (*mongo.InsertOneResult, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].InsertOne(v.Ctx, document)
}

func (v *DatabaseInfo) InsertMany(collections string, documents []interface{}) (*mongo.InsertManyResult, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].InsertMany(v.Ctx, documents)
}

func (v *DatabaseInfo) DeleteOne(collections string, document interface{}) (*mongo.DeleteResult, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].DeleteOne(v.Ctx, document)
}

func (v *DatabaseInfo) DeleteMany(collections string, documents []interface{}) (*mongo.DeleteResult, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].DeleteMany(v.Ctx, documents)
}

func (v *DatabaseInfo) FindOne(collections string, filter interface{}) (*mongo.SingleResult, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].FindOne(v.Ctx, filter), nil
}

func (v *DatabaseInfo) FindMany(collections string, filter interface{}) (*mongo.Cursor, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].Find(v.Ctx, filter)
}

func (v *DatabaseInfo) UpdateOne(collections string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].UpdateOne(v.Ctx, filter, update)
}

func (v *DatabaseInfo) UpdateMany(collections string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].UpdateMany(v.Ctx, filter, update)
}

func (v *DatabaseInfo) AddCollections(collections ...string) {
	for _, collection := range collections {
		_, ok := v.collections[collection]
		if !ok {
			v.collections[collection] = v.database.Collection(collection)
		}
	}
}
