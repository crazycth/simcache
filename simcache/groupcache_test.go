package simcache

import (
	"fmt"
	"log"
	"testing"
)

var mockDB = map[string]string{
	"test_A": "answer_A",
	"test_B": "answer_B",
	"test_C": "answer_C",
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(mockDB))
	cache := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := mockDB[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range mockDB {
		if view, err := cache.Get(k); err != nil || view.String() != v {
			t.Fatalf("failed to get value")
		}
		if _, err := cache.Get(k); err != nil || loadCounts[k] != 1 {
			t.Fatalf("cache %s miss", k)
		}
	}

	if view, err := cache.Get("unknown"); err == nil {
		t.Fatalf("the value of unknown should be empty , but %s get", view)
	}
}

func TestGetGroup(t *testing.T) {
	groupName := "scores"
	NewGroup(groupName, 2<<10, GetterFunc(
		func(key string) (bytes []byte, err error) { return }))
	if group := GetGroup(groupName); group == nil || group.name != groupName {
		t.Fatalf("tsetGroup fail")
	}
}
