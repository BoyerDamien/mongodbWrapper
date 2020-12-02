package mongodbWrapper

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

/******************************************************************************************************
*
*					Definition
*
*******************************************************************************************************/

// Database interface
// This interface will be initialized by the init method of the wrapper
// See wrapper for more informations
type Database interface {
	// This method allows you to add one or several collections to the database structure
	// If a collection is already present, it will skip it
	AddCollections(...string)

	// This method checks if a a collection is in the database structure or not
	// This methods is here to avoid some "dangerous" memory access
	// Warning: a collection that is not in the structure could be in the database
	// Just use the AddCollections method instead
	checkCollections(string) error

	// This method allows you to insert a document into the database
	// It's wrapped by the checkCollections method (SAFE)
	InsertOne(string, interface{}) (*mongo.InsertOneResult, error)

	// This method allows you to insert several documents into the database
	// It's wrapped by the checkCollections method (SAFE)
	InsertMany(string, []interface{}) (*mongo.InsertManyResult, error)

	// This method allows you to delete a document from the database
	// It's wrapped by the checkCollections method (SAFE)
	DeleteOne(string, interface{}) (*mongo.DeleteResult, error)

	// This method allows you to delete several documents from the database
	// It's wrapped by the checkCollections method (SAFE)
	DeleteMany(string, interface{}) (*mongo.DeleteResult, error)

	// This method allows you to find one document from the database based on filter
	// It's wrapped by the checkCollections method (SAFE)
	FindOne(string, interface{}) (*mongo.SingleResult, error)

	// This method allows you to find several documents from the database based on a filter
	// It's wrapped by the checkCollections method (SAFE)
	FindMany(string, interface{}) (*mongo.Cursor, error)

	// This method allows you to update one document from the database
	// It's wrapped by the checkCollections method (SAFE)
	UpdateOne(string, interface{}, interface{}) (*mongo.UpdateResult, error)

	// This method allows you to update serveral documents from the database
	// It's wrapped by the checkCollections method (SAFE)
	UpdateMany(string, interface{}, interface{}) (*mongo.UpdateResult, error)

	// Return the number of collections added
	CollectionNumber() int

	// Return the context
	GetContext() context.Context
}

/******************************************************************************************************
*
*					Implementation
*
*******************************************************************************************************/

// Database structure
type DatabaseInfo struct {
	name        string
	collections map[string]*mongo.Collection
	ctx         context.Context
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
	return v.collections[collections].InsertOne(v.ctx, document)
}

func (v *DatabaseInfo) InsertMany(collections string, documents []interface{}) (*mongo.InsertManyResult, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].InsertMany(v.ctx, documents)
}

func (v *DatabaseInfo) DeleteOne(collections string, document interface{}) (*mongo.DeleteResult, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].DeleteOne(v.ctx, document)
}

func (v *DatabaseInfo) DeleteMany(collections string, filter interface{}) (*mongo.DeleteResult, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].DeleteMany(v.ctx, filter)
}

func (v *DatabaseInfo) FindOne(collections string, filter interface{}) (*mongo.SingleResult, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].FindOne(v.ctx, filter), nil
}

func (v *DatabaseInfo) FindMany(collections string, filter interface{}) (*mongo.Cursor, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].Find(v.ctx, filter)
}

func (v *DatabaseInfo) UpdateOne(collections string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].UpdateOne(v.ctx, filter, update)
}

func (v *DatabaseInfo) UpdateMany(collections string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	if err := v.checkCollections(collections); err != nil {
		return nil, err
	}
	return v.collections[collections].UpdateMany(v.ctx, filter, update)
}

func (v *DatabaseInfo) AddCollections(collections ...string) {
	for _, collection := range collections {
		_, ok := v.collections[collection]
		if !ok {
			v.collections[collection] = v.database.Collection(collection)
		}
	}
}

func (v *DatabaseInfo) GetContext() context.Context {
	return v.ctx
}

func (v *DatabaseInfo) CollectionNumber() int {
	return len(v.collections)
}
