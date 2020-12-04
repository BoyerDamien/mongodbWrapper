package mongodbWrapper

import (
	"fmt"
	"testing"
)

var wrapper Wrapper = &WrapperData{}

func Test_Init_Good_URI(t *testing.T) {
	got := wrapper.Init("mongodb://localhost:27017")
	defer wrapper.Close()
	if got != nil {
		t.Errorf("Error: wrapper.Init(URI) -> uri = mongodb://localhost:27017")
	}
}

func Test_Init_Wrong_Uri(t *testing.T) {
	got := wrapper.Init("wrong")
	defer wrapper.Close()
	if got != nil {
		return
	}
	t.Errorf("Error: wrapper.Init(WRONG URI) -> uri = wrong")
}

func Test_GetDataBase(t *testing.T) {
	if err := wrapper.Init("mongodb://localhost:27017"); err != nil {
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
