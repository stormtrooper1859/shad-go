// +build !solution

package hogwarts

func GetCourseList(prereqs map[string][]string) []string {
	result := make([]string, 0)
	result_map := make(map[string]bool, 0)

	for k := range prereqs {
		dfs(&prereqs, &result_map, &result, k)
	}

	return result
}

func dfs(prereqs *map[string][]string, result_map *map[string]bool, result *[]string, vertex string) {
	in_result, entered := (*result_map)[vertex]
	if in_result {
		return
	}
	if entered {
		panic("cyclical dependence")
	}
	(*result_map)[vertex] = false

	for _, v := range (*prereqs)[vertex] {
		dfs(prereqs, result_map, result, v)
	}

	(*result_map)[vertex] = true
	*result = append(*result, vertex)
}
