package hw03frequencyanalysis

func Top10(str string) []string {
	sortedProvider := SortMap{}
	_, _ = sortedProvider.InitByStr(str)

	result, _ := sortedProvider.GetCollection(Desc, 10)
	return result
}

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
