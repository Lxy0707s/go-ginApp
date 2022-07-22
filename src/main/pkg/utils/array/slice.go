package array

import "reflect"

// StrSliceUnique is to delete duplicate entries in the string array
func StrSliceUnique(ss []string) []string {
	m := make(map[string]bool)
	for _, s := range ss {
		m[s] = true
	}
	var r []string
	for s := range m {
		r = append(r, s)
	}
	return r
}

// StrSliceUnique is to delete duplicate entries in the string array
// 会保证切片的顺序
func StrSliceUnique2(dataSlice []string) []string {
	result := make([]string, 0)
	temp := map[string]struct{}{}
	for _, item := range dataSlice {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// IntSliceUnique is to delete duplicate entries in the int array
func IntSliceUnique(ss []int) []int {
	m := make(map[int]bool)
	for _, s := range ss {
		m[s] = true
	}
	var r []int
	for s := range m {
		r = append(r, s)
	}
	return r
}

// StrFind is check if subStr in ss
func StrFind(ss []string, subStr string) bool {
	for _, s := range ss {
		if s == subStr {
			return true
		}
	}
	return false
}

// IntFind is check if sub in ss
func IntFind(ss []int, sub int) bool {
	for _, s := range ss {
		if s == sub {
			return true
		}
	}
	return false
}

// StrIndex is return the index of subStr in ss
func StrIndex(ss []string, subStr string) int {
	for i, s := range ss {
		if s == subStr {
			return i
		}
	}
	return -1
}

// Replace replace the old item in ss with new item
func Replace(ss []string, old, new string) []string {
	var r []string
	for _, v := range ss {
		if v == old {
			r = append(r, new)
			continue
		}
		r = append(r, v)
	}
	return r
}

// Delete is to delete ele in ss
func Delete(ss []string, ele string) []string {
	var r []string
	for _, v := range ss {
		if v == ele {
			continue
		}
		r = append(r, v)
	}
	return r
}

// HaveIntersection check is slice1 and slice2 have intersection
func HaveIntersection(slice1, slice2 []string) bool {
	for _, ele1 := range slice1 {
		for _, ele2 := range slice2 {
			if ele1 == ele2 {
				return true
			}
		}
	}
	return false
}

// Include check the items of slice2 is all in slice1
func Include(slice1, slice2 []string) bool {
	for _, ele := range slice2 {
		if !StrFind(slice1, ele) {
			return false
		}
	}
	return true
}

// 判断obj是否在target中，target支持的类型array,slice,map
/**
   example1:
    c := "a"
   	d := [4]string{"b", "c", "d", "a"}
   	fmt.Println(Contain(c, d))

   example2:
	 a := 1
     b := [3]int{1, 2, 3}

     fmt.Println(Contain(a, b))
*/
func Contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

// Computes the difference of string slice
func StrSliceDiff(slice1, slice2 []string) (diffslice []string) {
	for _, v := range slice1 {
		if !StrFind(slice2, v) {
			diffslice = append(diffslice, v)
		}
	}
	return
}

// Computes the difference of int slice
func IntSliceDiff(slice1, slice2 []int) (diffslice []int) {
	for _, v := range slice1 {
		if !IntFind(slice2, v) {
			diffslice = append(diffslice, v)
		}
	}
	return
}

//两map交集
func TwoMapIntersection(map1 map[string][]string, map2 map[string][]string) map[string][]string {
	results := map[string][]string{}
	for k, v := range map1 {
		if _, ok := map2[k]; ok {
			StrListIntersection := TwoStrListIntersection(v, map2[k])
			results[k] = StrListIntersection
		}
	}
	return results
}

//两数组交集
func TwoStrListIntersection(ss1 []string, ss2 []string) []string {
	var results []string
	res := map[string]int{}
	for _, item := range ss1 {
		res[item] += 1
	}
	for _, item := range ss2 {
		res[item] += 1
	}
	for k, v := range res {
		if v >= 2 {
			results = append(results, k)
		}
	}
	return results
}
