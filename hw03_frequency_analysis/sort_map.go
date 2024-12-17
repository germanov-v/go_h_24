package hw03frequencyanalysis

import (
	"fmt"
	"strings"
)

type TypeSort int

const (
	Asc  TypeSort = iota // ascending. first in top
	Desc                 // descending. last in top
)

type sortMapItem struct {
	str      string
	position int
	count    int
}

type SortMap struct {
	dictionary       map[string]sortMapItem
	sortedCollection []string
}

func (strMap *SortMap) InitByStr(str string) (bool, error) {

	var items = strings.Split(str, " ")
	strMap.dictionary = make(map[string]sortMapItem)
	strMap.sortedCollection = make([]string, 0, len(items))
	for _, item := range items {
		if item != "" {
			_, err := strMap.AddItem(item)
			if err != nil {
				return false, fmt.Errorf("error was occurred by %s. Internal error is %s", item, err)
			}
		}

	}

	if len(strMap.sortedCollection) != len(strMap.dictionary) {
		return false, fmt.Errorf("not eqaul length by sortedList (len=%d) and map (len=%d)", len(strMap.sortedCollection),
			len(strMap.dictionary))
	}

	return true, nil
}

func (strMap *SortMap) AddItem(str string) (bool, error) {

	item, exits := strMap.dictionary[str]

	if !exits {
		item = sortMapItem{
			str:      str,
			position: -1,
			count:    0,
		}
	}

	item.count++

	_, err := strMap.addToStrToSortedList(&item)
	strMap.dictionary[str] = item
	if err != nil {
		// rollback
		delete(strMap.dictionary, str)
		return false, fmt.Errorf("error occuкred by string %s", str)
	}

	return true, err
}

func (strMap *SortMap) addToStrToSortedList(item *sortMapItem) (int, error) {

	startPosition := 0
	needOffset := false
	if item.position == -1 {
		strMap.sortedCollection = append([]string{item.str}, strMap.sortedCollection...)
		item.position = 0
		needOffset = true
	} else {
		startPosition = item.position
	}

	stopOffsetIndex := 0
	for i := startPosition + 1; i < len(strMap.sortedCollection); i++ {
		if strMap.dictionary[strMap.sortedCollection[i-1]].count > strMap.dictionary[strMap.sortedCollection[i]].count {
			tempMore := strMap.dictionary[strMap.sortedCollection[i-1]]
			tempLess := strMap.dictionary[strMap.sortedCollection[i]]
			tempMore.position += 1
			tempLess.position -= 1
			strMap.dictionary[strMap.sortedCollection[i-1]] = tempMore
			strMap.dictionary[strMap.sortedCollection[i]] = tempLess
			strMap.sortedCollection[i-1], strMap.sortedCollection[i] = strMap.sortedCollection[i], strMap.sortedCollection[i-1]
			stopOffsetIndex = i
		} else if stopOffsetIndex > 0 {
			// тут мы понимаем, что дальше двигаться нет смысла - остальные элементы массива
			break
		} else if i > startPosition+1 && stopOffsetIndex == 0 {
			// тут мы понимаем, что дальше двигаться так же не имеет смысла - весь массив остается упорядоченным
			break
		}
	}

	if needOffset {
		for i := stopOffsetIndex + 1; i < len(strMap.sortedCollection); i++ {
			temp := strMap.dictionary[strMap.sortedCollection[i]]
			temp.position++
			strMap.dictionary[strMap.sortedCollection[i]] = temp
		}
	}

	return item.position, nil
}

func (strMap *SortMap) GetCollection(asc TypeSort, count int) ([]string, error) {

	if count < 1 || count > len(strMap.sortedCollection) {
		return nil, fmt.Errorf("Invalid count %d", count)
	}

	var result []string

	switch asc {
	case Asc:
		for i := 0; i < count; i++ {
			result = append(result, strMap.sortedCollection[i])
		}

	case Desc:
		for i := len(strMap.sortedCollection) - 1; i >= len(strMap.sortedCollection)-count; i-- {
			result = append(result, strMap.sortedCollection[i])
		}

	default:
		return nil, fmt.Errorf("Invalid typesort")
	}
	return result, nil
}
