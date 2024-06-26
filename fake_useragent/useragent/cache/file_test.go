package cache

import (
	"fmt"
	"testing"

	"github.com/EDDYCJY/fake-useragent/setting"
)

var f = NewFileCache(GetTempDir(), fmt.Sprintf(setting.TEMP_FILE_TEST_NAME, setting.VERSION))

func TestFile_Write(t *testing.T) {
	err := f.Write([]byte("test"))
	if err != nil {
		t.Errorf("f.Write err: %v", err)
	}
}

func TestFile_Read(t *testing.T) {
	value, err := f.Read()
	if err != nil {
		t.Errorf("f.Read err: %v", err)
	}

	str := string(value)
	if str != "test" {
		t.Errorf("Expected 'test', got %s", str)
	}
}

func TestFile_Remove(t *testing.T) {
	err := f.Remove()
	if err != nil {
		t.Errorf("f.Remove err: %v", err)
	}
}

func TestFile_IsExist(t *testing.T) {
	exist, err := f.IsExist()
	if exist == true {
		t.Errorf("Expected false, got true")
	}
	if err != nil {
		t.Errorf("f.IsExist err: %v", err)
	}
}
