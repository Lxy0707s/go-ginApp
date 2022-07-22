package array

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

// HaveDifference check is slice1 and slice2 have difference
func HaveDifference(slice1, slice2 []string) bool {
	for _, ele1 := range slice1 {
		for _, ele2 := range slice2 {
			if ele1 != ele2 {
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
