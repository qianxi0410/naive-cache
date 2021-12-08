package naive_cache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestSprintN(t *testing.T) {
	str := fmt.Sprintf("%sgroup1/key1", defaultBasePath)

	s := strings.SplitN(str[len(defaultBasePath):], "/", 2)

	if len(s) != 2 {
		t.Logf("error lens")
	}

	if s[0] != "group1" {
		t.Fail()
	}

	if s[1] != "key1" {
		t.Fail()
	}
}

func TestServeHttp(t *testing.T) {
	var db = map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}

	NewGroup("test", GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}), 2<<10)

	addr := "localhost:9999"
	peers := NewHttpPool(addr)
	log.Println("naive-cache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
