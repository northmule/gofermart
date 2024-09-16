package storage

import (
	"sync"
	"testing"
)

func TestSessionStorage_Add(t *testing.T) {
	storage := NewSessionStorage()

	storage.Add("key1", "value1")
	if val, ok := storage.Values["key1"]; !ok || val != "value1" {
		t.Errorf("Expected key1 with value1, but got %v", val)
	}
}

func TestSessionStorage_Get(t *testing.T) {
	storage := NewSessionStorage()

	storage.Add("key1", "value1")
	if val, ok := storage.Get("key1"); !ok || val != "value1" {
		t.Errorf("Expected key1 with value1, but got %v", val)
	}
}

func TestSessionStorage_GetAll(t *testing.T) {
	storage := NewSessionStorage()

	storage.Add("key1", "value1")
	storage.Add("key2", "value2")

	allValues := storage.GetAll()
	if len(allValues) != 2 {
		t.Errorf("Expected 2 values, but got %d", len(allValues))
	}
	if allValues["key1"] != "value1" || allValues["key2"] != "value2" {
		t.Errorf("Expected values to be key1:value1, key2:value2, but got %v", allValues)
	}
}

func TestSessionStorage_ThreadSafety(t *testing.T) {
	storage := NewSessionStorage()
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			storage.Add("key", "value")
		}(i)
	}

	wg.Wait()

	if len(storage.Values) != 1 {
		t.Errorf("Expected 1 unique key, but got %d", len(storage.Values))
	}
}
