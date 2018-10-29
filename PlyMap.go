package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"sort"
// 	"sync"
// )

// //PlyMap is an marshal-able sync.Map
// type PlyMap struct {
// 	sync.Map
// }

// //MarshalJSON marshals sync.Map to json, where data is sorted by key
// func (plyMap *PlyMap) MarshalJSON() ([]byte, error) {
// 	tmpMap := make(map[int]float64)
// 	keys := make([]int, 0)

// 	plyMap.Range(func(k, v interface{}) bool {
// 		tmpMap[k.(int)] = v.(float64)
// 		keys = append(keys, k.(int))
// 		return true
// 	})

// 	if len(keys) < 1 {
// 		return json.Marshal(keys)
// 	}

// 	sort.Ints(keys)
// 	sortedValues := make([]float64, keys[len(keys)-1]+1)

// 	for _, k := range keys {
// 		if k == 0 {
// 			fmt.Println("wa", tmpMap[k], tmpMap[k]/4616)
// 		}
// 		sortedValues[k] = tmpMap[k] / 4616
// 	}

// 	return json.Marshal(sortedValues)
// }

// //StoreOrAdd stores the value in the plymap, or adds value to stored value
// func (plyMap *PlyMap) StoreOrAdd(key, value interface{}, totalGames uint32) {
// 	val, loaded := plyMap.LoadOrStore(key, value)
// 	if loaded {
// 		plyMap.Store(key, val.(float64)+value.(float64))
// 		// plyMap.Store(key, ((val.(float64)*float64(totalGames))+value.(float64))/(float64(totalGames)+1))
// 	}
// }
