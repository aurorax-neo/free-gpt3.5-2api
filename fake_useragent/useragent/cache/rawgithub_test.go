package cache

import (
	"fmt"
	"github.com/EDDYCJY/fake-useragent/setting"
	"net/http"
	"testing"
)

var r = NewRawCache(setting.CACHE_URL, fmt.Sprintf(setting.TEMP_FILE_NAME, setting.CACHE_VERSION))
var rawResp *http.Response

func TestRaw_Get(t *testing.T) {
	resp, exist, err := r.Get()
	if err != nil {
		t.Errorf("r.Get err: %v", err)
	}
	if exist == false {
		t.Errorf("r.Get not exist")
	}

	rawResp = resp
}

func TestRaw_IsExist(t *testing.T) {
	exist := r.IsExist(rawResp)
	if exist == false {
		t.Errorf("r.IsExist not exist")
	}
}

func TestRaw_Read(t *testing.T) {
	defer rawResp.Body.Close()
	body, err := r.Read(rawResp.Body)
	if err != nil {
		t.Errorf("r.Get err: %v", err)
	}
	if len(body) == 0 {
		t.Errorf("r.Read len is zero")
	}
}
