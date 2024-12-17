package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type ItemStr struct {
	str   string
	count int
}

func Top10(str string) []string {

	items := strings.Fields(str) //strings.Split(str, " ")

	dict := make(map[string]int)
	for _, item := range items {
		dict[item]++ // dict[item]+=1
	}

	collection := make([]ItemStr, 0, len(dict))

	for key, value := range dict {
		collection = append(collection, ItemStr{key, value})
	}

	//sort.Slice(collection, func(i, j int) bool {
	//	return false
	//})

	sort.Slice(collection, func(i, j int) bool {
		if collection[i].count == collection[j].count {
			return collection[i].str > collection[j].str
		}
		return collection[i].count < collection[j].count
	})

	result := make([]string, 0, 10)
	limit := len(collection) - 10
	if limit < 10 {
		limit = 0
	}
	for i := len(collection) - 1; i >= limit; i-- {
		result = append(result, collection[i].str)
	}

	return result
}

//func Top10(str string) []string {
//	sortedProvider := SortMap{}
//	_, _ = sortedProvider.InitByStr(str)
//
//	result, _ := sortedProvider.GetCollection(Desc, 10)
//	return result
//}

//type SortMap struct {
//	dictionary       map[string]int
//	sortedCollection []string
//}
//
//func (srtM *SortMap) AddItem(str string) (bool, error) {
//	return true, nil
//}
//
//func (srtM *SortMap) GetCollection(str string) (bool, error) {
//	return true, nil
//}
