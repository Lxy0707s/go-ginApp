package mapsort

import (
	"sort"
)

type valueType string

var (
	ArrayType  valueType = "array"
	StringType valueType = "string"
	MapType    valueType = "map"
	MapMapType valueType = "mapmap"
	BoolType   valueType = "bool"
)

// SortStringKeys get keys after sort
func SortStringKeys(m interface{}, vType valueType) []string {
	if m == nil {
		return nil
	}
	if vType == ArrayType {
		return sortKeyOfArray(m.(map[string][]string))
	}
	if vType == StringType {
		return sortKeyOfString(m.(map[string]string))
	}
	if vType == MapType {
		return sortKeyOfMap(m.(map[string]map[string]string))
	}
	if vType == BoolType {
		return sortKeyOfBool(m.(map[string]bool))
	}
	if vType == MapMapType {
		return sortKeyOfMapMap(m.(map[string]map[string]map[string]string))
	}
	return nil
}

// SortMapKey sort map keys
func SortMapKey(m map[string]interface{}) []string {
	result := make([]string, 0)
	for k := range m {
		result = append(result, k)
	}
	sort.Strings(result)
	return result
}

func sortKeyOfArray(m map[string][]string) []string {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortKeyOfString(m map[string]string) []string {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortKeyOfBool(m map[string]bool) []string {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortKeyOfMap(m map[string]map[string]string) []string {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortKeyOfMapMap(m map[string]map[string]map[string]string) []string {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
