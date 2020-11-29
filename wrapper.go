package mongodbWrapper

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Wrapper interface {
	Init(URI string)
	GetDatabase(name string) Database
}

type WrapperData struct {
	client    *mongo.Client
	ctx       context.Context
	databases map[string]Database
	Cancel    context.CancelFunc
}

func (v *WrapperData) Init(URI string) {
	var err error
	v.client, err = mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		log.Fatal(err)
	}
	v.ctx, v.Cancel = context.WithTimeout(context.Background(), 10*time.Second)
	if err := v.client.Connect(v.ctx); err != nil {
		log.Fatal(err)
	}
	if err := v.client.Ping(v.ctx, nil); err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully connected to %s\n", URI)
	v.databases = make(map[string]Database)
}

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
