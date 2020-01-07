package godash

type StringMapper func(item string, i int) string
type StringChecker func(item string, i int) bool

func StringMap(list []string, mapper StringMapper) []string {
	newList := []string{}
	for i, item := range list {
		newList = append(newList, mapper(item, i))
	}
	return newList
}

func StringPartition(list []string, checker StringChecker) ([]string, []string) {
	newLeftList := []string{}
	newRightList := []string{}
	for i, item := range list {
		if checker(item, i) {
			newLeftList = append(newLeftList, item)
		} else {
			newRightList = append(newRightList, item)
		}
	}
	return newLeftList, newRightList
}

func StringFilter(list []string, checker StringChecker) []string {
	newList, _ := StringPartition(list, checker)
	return newList
}

func StringUniq(list []string) (result []string) {
	strMap := make(map[string]bool)
	for _, str := range list {
		strMap[str] = true
	}

	for i := range list {
		str := list[i]
		if _, ok := strMap[str]; ok {
			delete(strMap, str)
			result = append(result, str)
		}
	}
	return
}

func StringReject(list []string, checker StringChecker) []string {
	_, newList := StringPartition(list, checker)
	return newList
}

func StringEvery(list []string, checker StringChecker) bool {
	var checked = true
	for i, item := range list {
		checked = checked && checker(item, i)
	}
	return checked
}

func StringSome(list []string, checker StringChecker) bool {
	for i, item := range list {
		if checker(item, i) {
			return true
		}
	}
	return false
}

func StringIncludes(list []string, target string) bool {
	return StringSome(list, func(item string, i int) bool {
		return item == target
	})
}

func StringFind(list []string, checker StringChecker) string {
	for i, item := range list {
		if checker(item, i) {
			return item
		}
	}
	return ""
}
