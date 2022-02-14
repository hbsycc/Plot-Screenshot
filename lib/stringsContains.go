package lib

import "sort"

func StringsContains(a []string, x string) bool {
	sort.Strings(a)
	index := sort.SearchStrings(a, x)
	//index的取值：0 ~ (len(str_array)-1)
	return index < len(a) && a[index] == x
}
