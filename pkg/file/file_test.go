package file

import (
	"encoding/json"
	"os"
	"testing"
)

func TestFile(t *testing.T) {

	const testFile = "./testFile.txt"

	expectedId := 10

	type testModel struct {
		Id int `json:"id"`
	}

	if Exists(testFile) {
		t.Errorf("file should not yet exist")
	}

	err := WriteStringAsLine(testFile, func(yield func([]byte) bool) {
		data, _ := json.Marshal(&testModel{Id: expectedId})
		yield(data)
		return
	})

	if err != nil {
		t.Errorf("file creation failed")
	}

	if !Exists(testFile) {
		t.Errorf("file should have been created")
	}

	readData := testModel{}
	for range ReadJsonLineAsStruct(testFile, &readData) {
	}

	if readData.Id != 10 {
		t.Errorf("expected %v, got %v", expectedId, readData.Id)
	}

	err = os.Remove(testFile)
	if err != nil {
		t.Errorf("could not delete %v", testFile)
	}
}
