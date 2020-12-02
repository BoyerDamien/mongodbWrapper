package mongodbWrapper

import (
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

type Model struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func Test_AddCollection(t *testing.T) {
	if err := wrapper.Init("mongodb://localhost:27017"); err != nil {
		t.Errorf("Error: wrapper.Init(URI) -> uri = mongodb://localhost:27017: %s", err)
	}
	defer wrapper.Close()
	db := wrapper.GetDatabase("test")

	for i := 0; i < 10; i++ {
		db.AddCollections("coll1")
	}
	if db.CollectionNumber() > 1 {
		t.Errorf("Error: db.AddCollections() -> wrong number of collection : ok = 1 ; ko = %d", db.CollectionNumber())
	}
	for i := 1; i < 5; i++ {
		db.AddCollections(fmt.Sprintf("test%d", i))
	}
	if db.CollectionNumber() != 5 {
		t.Errorf("Error: db.AddCollections() -> wrong number of collection : ok = 5 ; ko = %d", db.CollectionNumber())
	}
}

func Test_InsertOne(t *testing.T) {
	if err := wrapper.Init("mongodb://localhost:27017"); err != nil {
		t.Errorf("Error: wrapper.Init(URI) -> uri = mongodb://localhost:27017: %s", err)
	}
	defer wrapper.Close()
	db := wrapper.GetDatabase("test")
	db.AddCollections("coll1")

	testModel := Model{"test", 30}
	_, err := db.InsertOne("coll1", testModel)
	if err != nil {
		t.Errorf("Error: db.InsertOne() -> %s", err)
	}

	var decoded Model
	res, err := db.FindOne("coll1", testModel)
	if err != nil {
		t.Errorf("Error: db.FindOne() -> %s", err)
	}
	res.Decode(&decoded)
	if decoded.Age != testModel.Age || decoded.Name != testModel.Name {
		t.Errorf("Error: db.FindOne() -> wrong result %v", decoded)
	}
	_, err = db.UpdateOne("coll1", decoded, bson.M{"$set": bson.M{"name": "changed"}})
	if err != nil {
		t.Errorf("Error: db.UpdateOne() -> %s", err)
	}
	_, err = db.DeleteOne("coll1", bson.M{"name": "changed"})
	if err != nil {
		t.Errorf("Error: db.DeleteOne() -> %s", err)
	}
}

func Test_InserMany(t *testing.T) {
	if err := wrapper.Init("mongodb://localhost:27017"); err != nil {
		t.Errorf("Error: wrapper.Init(URI) -> uri = mongodb://localhost:27017: %s", err)
	}
	defer wrapper.Close()
	db := wrapper.GetDatabase("test2")
	db.AddCollections("coll2")

	var toInsert []interface{} = []interface{}{Model{"test", 30}, Model{"test", 30}}
	_, err := db.InsertMany("coll2", toInsert)
	if err != nil {
		t.Errorf("Error: db.InsertMany() -> %s", err)
	}

	var decoded []Model
	cur, err := db.FindMany("coll2", bson.M{"name": "test"})
	if err != nil {
		t.Errorf("Error: db.FindMany() -> %s", err)
	}
	_, err = db.UpdateMany("coll2", bson.M{"name": "test"}, bson.M{"$set": bson.M{"name": "changed"}})
	if err != nil {
		t.Errorf("Error: db.UpdateMany() -> %s", err)
	}
	if err := cur.All(db.GetContext(), &decoded); err != nil {
		t.Errorf("Error: cur.All() -> %s", err)
	}
	if len(decoded) != 2 {
		t.Errorf("Error: cur.All() -> wrong number of model finded: ok = %d -- ko: %d", 2, len(decoded))
	}

	_, err = db.DeleteMany("coll2", bson.M{"name": "changed"})
	if err != nil {
		t.Errorf("Error: db.DeleteMany() -> %s", err)
	}
}
