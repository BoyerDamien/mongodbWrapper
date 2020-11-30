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
	FindOne(string, interface{}, interface{}) error
	FindMany(string, interface{}, interface{}, interface{}) error
}

type DatabaseInfo struct {
	name        string
	collections map[string]*mongo.Collection
	ctx         context.Context
	database    *mongo.Database
}

func (v *DatabaseInfo) InsertOne(collections string, document interface{}) (*mongo.InsertOneResult, error) {
	return v.collections[collections].InsertOne(v.ctx, document)
}

func (v *DatabaseInfo) InsertMany(collections string, documents []interface{}) (*mongo.InsertManyResult, error) {
	return v.collections[collections].InsertMany(v.ctx, documents)
}

func (v *DatabaseInfo) DeleteOne(collections string, document interface{}) (*mongo.DeleteResult, error) {
	return v.collections[collections].DeleteOne(v.ctx, document)
}

func (v *DatabaseInfo) DeleteMany(collections string, documents []interface{}) (*mongo.DeleteResult, error) {
	return v.collections[collections].DeleteMany(v.ctx, documents)
}

func (v *DatabaseInfo) FindOne(collections string, filter interface{}, model interface{}) error {
	result := v.collections[collections].FindOne(v.ctx, filter)
	if result.Err() != nil {
		return result.Err()
	}
	if err := result.Decode(model); err != nil {
		return err
	}
	return nil
}

func (v *DatabaseInfo) FindMany(collections string, filter interface{}, model interface{}, results []interface{}) error {
	cur, err := v.collections[collections].Find(v.ctx, filter)
	if err != nil {
		return err
	}
	for cur.Next(v.ctx) {
		err = cur.Decode(model)
		if err != nil {
			return err
		}
		results = append(results, model)
	}
	return nil
}

func (v *DatabaseInfo) AddCollections(collections ...string) {
	for _, collection := range collections {
		_, ok := v.collections[collection]
		if !ok {
			v.collections[collection] = v.database.Collection(collection)
		}
	}
}
