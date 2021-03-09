package mongodbwrapper

import (
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	wrapper     Wrapper            = &WrapperData{}
	credentials options.Credential = options.Credential{
		Username: "test",
		Password: "password",
	}
)

func Test_Init_Good_URI(t *testing.T) {
	got := wrapper.Init("mongodb://localhost:27017", credentials)
	defer wrapper.Close()
	if got != nil {
		t.Errorf("Error: wrapper.Init(URI) -> uri = mongodb://localhost:27017")
	}
}

func Test_GetDataBase(t *testing.T) {
	if err := wrapper.Init("mongodb://localhost:27017", credentials); err != nil {
		t.Errorf("Error: wrapper.Init(URI) -> uri = mongodb://localhost:27017")
	}
	defer wrapper.Close()

	for i := 0; i < 10; i++ {
		wrapper.GetDatabase("test0")
	}
	if wrapper.DatabaseNumber() > 1 {
		t.Errorf("Error: wrapper.GetDatabase() -> wrong number of database : ok = 1 ; ko = %d", wrapper.DatabaseNumber())
	}
	for i := 1; i < 5; i++ {
		wrapper.GetDatabase(fmt.Sprintf("test%d", i))
	}
	if wrapper.DatabaseNumber() != 5 {
		t.Errorf("Error: wrapper.GetDatabase() -> wrong number of database : ok = 5 ; ko = %d", wrapper.DatabaseNumber())
	}
}
