package mongodbWrapper

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Database interface {
	AddCollections(...string)
	InsertOne(string, interface{}) (*mongo.InsertOneResult, error)
	InsertMany(string, []interface{}) (*mongo.InsertManyResult, error)
	DeleteOne(string, interface{}) (*mongo.DeleteResult, error)
	DeleteMany(string, []interface{}) (*mongo.DeleteResult, error)
	FindOne(string, interface{}) *mongo.SingleResult
	FindMany(string, interface{}) (*mongo.Cursor, error)
}

type DatabaseInfo struct {
	name        string
	collections map[string]*mongo.Collection
	Ctx         context.Context
	database    *mongo.Database
}

func (v *DatabaseInfo) InsertOne(collections string, document interface{}) (*mongo.InsertOneResult, error) {
	return v.collections[collections].InsertOne(v.Ctx, document)
}

func (v *DatabaseInfo) InsertMany(collections string, documents []interface{}) (*mongo.InsertManyResult, error) {
	return v.collections[collections].InsertMany(v.Ctx, documents)
}

func (v *DatabaseInfo) DeleteOne(collections string, document interface{}) (*mongo.DeleteResult, error) {
	return v.collections[collections].DeleteOne(v.Ctx, document)
}

func (v *DatabaseInfo) DeleteMany(collections string, documents []interface{}) (*mongo.DeleteResult, error) {
	return v.collections[collections].DeleteMany(v.Ctx, documents)
}

func (v *DatabaseInfo) FindOne(collections string, filter interface{}) *mongo.SingleResult {
	return v.collections[collections].FindOne(v.Ctx, filter)
}

func (v *DatabaseInfo) FindMany(collections string, filter interface{}) (*mongo.Cursor, error) {
	return v.collections[collections].Find(v.Ctx, filter)
}

func (v *DatabaseInfo) AddCollections(collections ...string) {
	for _, collection := range collections {
		_, ok := v.collections[collection]
		if !ok {
			v.collections[collection] = v.database.Collection(collection)
		}
	}
}
