package mongodbWrapper

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

// This interface will be initialized by the init method of the wrapper
// See wrapper for more informations
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

// Database structure
type DatabaseInfo struct {
	name        string
	collections map[string]*mongo.Collection
	Ctx         context.Context
	database    *mongo.Database
}

// This method checks if a a collection is in the database structure or not
// This methods is here to avoid some "dangerous" memory access
// Warning: a collection that is not in the structure could be in the database
// Just use the AddCollections method instead
func (v *DatabaseInfo) checkCollections(collections string) error {
	_, ok := v.collections[collections]
	if !ok {
		return fmt.Errorf("Collection %s does not exist", collections)
	}
	return nil
}

// This method allows you to insert a document into the database
// It's wrapped by the checkCollections method (SAFE)
func (v *DatabaseInfo) InsertOne(collections string, document interface{}) (*mongo.InsertOneResult, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].InsertOne(v.Ctx, document)
}

// This method allows you to insert several documents into the database
// It's wrapped by the checkCollections method (SAFE)
func (v *DatabaseInfo) InsertMany(collections string, documents []interface{}) (*mongo.InsertManyResult, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].InsertMany(v.Ctx, documents)
}

// This method allows you to delete a document from the database
// It's wrapped by the checkCollections method (SAFE)
func (v *DatabaseInfo) DeleteOne(collections string, document interface{}) (*mongo.DeleteResult, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].DeleteOne(v.Ctx, document)
}

// This method allows you to delete several documents from the database
// It's wrapped by the checkCollections method (SAFE)
func (v *DatabaseInfo) DeleteMany(collections string, documents []interface{}) (*mongo.DeleteResult, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].DeleteMany(v.Ctx, documents)
}

// This method allows you to find one document from the database based on filter
// It's wrapped by the checkCollections method (SAFE)
func (v *DatabaseInfo) FindOne(collections string, filter interface{}) (*mongo.SingleResult, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].FindOne(v.Ctx, filter), nil
}

// This method allows you to find several documents from the database based on a filter
// It's wrapped by the checkCollections method (SAFE)
func (v *DatabaseInfo) FindMany(collections string, filter interface{}) (*mongo.Cursor, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].Find(v.Ctx, filter)
}

// This method allows you to update one document from the database
// It's wrapped by the checkCollections method (SAFE)
func (v *DatabaseInfo) UpdateOne(collections string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].UpdateOne(v.Ctx, filter, update)
}

// This method allows you to update serveral documents from the database
// It's wrapped by the checkCollections method (SAFE)
func (v *DatabaseInfo) UpdateMany(collections string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].UpdateMany(v.Ctx, filter, update)
}

// This method allows you to add one or several collections to the database structure
// If a collection is already present, it will skip it
func (v *DatabaseInfo) AddCollections(collections ...string) {
	for _, collection := range collections {
		_, ok := v.collections[collection]
		if !ok {
			v.collections[collection] = v.database.Collection(collection)
		}
	}
}
