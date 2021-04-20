// +build !solution

package hogwarts

func GetCourseList(prereqs map[string][]string) []string {
	result := make([]string, 0)
	resultMap := make(map[string]bool, 0)

	for k := range prereqs {
		dfs(prereqs, resultMap, &result, k)
	}

	return result
}

func dfs(prereqs map[string][]string, resultMap map[string]bool, result *[]string, vertex string) {
	inResult, entered := resultMap[vertex]
	if inResult {
		return
	}
	if entered {
		panic("cyclical dependence")
	}
	resultMap[vertex] = false

	for _, v := range prereqs[vertex] {
		dfs(prereqs, resultMap, result, v)
	}

	resultMap[vertex] = true
	*result = append(*result, vertex)
}
