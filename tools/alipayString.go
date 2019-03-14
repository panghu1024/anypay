package tools

import (
	"sort"
	"strings"
	"net/url"
)

//参数排序
func SortData(data url.Values) (str string,err error){

	var list = make([]string, 0, 0)
	for key := range data {
		var value = strings.TrimSpace(data.Get(key))
		if len(value) > 0 {
			list = append(list, key+"="+value)
		}
	}
	sort.Strings(list)
	return strings.Join(list, "&"),nil
}
