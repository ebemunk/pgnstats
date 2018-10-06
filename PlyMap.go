package main

import (
	"encoding/json"
	"sort"
	"sync"
)

type PlyMap struct {
	sync.Map
}

func (f PlyMap) MarshalJSON() ([]byte, error) {
	tmpMap := make(map[int]float64)
	keys := make([]int, 0)

	f.Range(func(k, v interface{}) bool {
		tmpMap[k.(int)] = v.(float64)
		keys = append(keys, k.(int))
		return true
	})

	sort.Ints(keys)
	sortedValues := make([]float64, keys[len(keys)-1]+1)

	// log.Println("yo", max)

	for _, k := range keys {
		sortedValues[k] = tmpMap[k]
	}

	return json.Marshal(sortedValues)
}
