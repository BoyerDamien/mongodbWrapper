package mongodbwrapper

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/******************************************************************************************************
*
*					Definition
*
*******************************************************************************************************/

// Wrapper interface
type Wrapper interface {

	// This method will:
	//	- Create a mongodb client
	//	- Create a context and its cancel function -> don't forget to use "defer (*Wrapper).Close()"
	//	- Test the connection with the database
	// If something went wrong, an error will be returned
	Init(URI string, credentials options.Credential) error

	// This method will create a database and return an interface between it and your program
	// If the database already exists, no new database will be created
	// The interface will be always created
	GetDatabase(name string) Database

	// Return the number of databases that you added to the wrapper
	DatabaseNumber() int

	// Close the connection with the database
	Close()
}

/******************************************************************************************************
*
*					Implementation
*
*******************************************************************************************************/

// WrapperData structure
type WrapperData struct {
	client    *mongo.Client
	ctx       context.Context
	databases map[string]Database
	cancel    context.CancelFunc
}

// Init methods init db connection
func (v *WrapperData) Init(URI string, credentials options.Credential) error {
	var err error
	v.client, err = mongo.NewClient(options.Client().ApplyURI(URI).SetAuth(credentials))
	if err != nil {
		return err
	}
	v.ctx, v.cancel = context.WithCancel(context.Background())
	if err := v.client.Connect(v.ctx); err != nil {
		return err
	}
	if err := v.client.Ping(v.ctx, nil); err != nil {
		return err
	}
	v.databases = make(map[string]Database)
	return nil
}

// GetDatabase returns or create a db form a name if not already exists
func (v *WrapperData) GetDatabase(name string) Database {
	_, ok := v.databases[name]
	if !ok {
		v.databases[name] = &DatabaseInfo{
			name,
			make(map[string]*mongo.Collection),
			v.ctx,
			v.client.Database(name),
		}
	}
	return v.databases[name]
}

// DatabaseNumber return the number of db stored
func (v *WrapperData) DatabaseNumber() int {
	return len(v.databases)
}

// Close close db connection
func (v *WrapperData) Close() {
	v.cancel()
}
