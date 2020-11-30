package mongodbWrapper

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/******************************************************************************************************
*
*					Definition
*
*******************************************************************************************************/

type Wrapper interface {

	// This method will:
	//	- Create a mongodb client
	//	- Create a context and its cancel function -> don't forget to use "defer (*Wrapper).Close()"
	//	- Test the connection with the database
	// If something went wrong, an error will be returned
	Init(URI string) error

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

// The wrapper structure
type WrapperData struct {
	client    *mongo.Client
	Ctx       context.Context
	databases map[string]Database
	cancel    context.CancelFunc
}

func (v *WrapperData) Init(URI string) error {
	var err error
	v.client, err = mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		return err
	}
	v.Ctx, v.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	if err := v.client.Connect(v.Ctx); err != nil {
		return err
	}
	if err := v.client.Ping(v.Ctx, nil); err != nil {
		return err
	}
	v.databases = make(map[string]Database)
	return nil
}

func (v *WrapperData) GetDatabase(name string) Database {
	_, ok := v.databases[name]
	if !ok {
		v.databases[name] = &DatabaseInfo{
			name,
			make(map[string]*mongo.Collection),
			v.Ctx,
			v.client.Database(name),
		}
	}
	return v.databases[name]
}

func (v *WrapperData) DatabaseNumber() int {
	return len(v.databases)
}

func (v *WrapperData) Close() {
	v.cancel()
}
