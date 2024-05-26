package algorithms

func SliceIntersection(a, b []string) []string {
	set := make(map[string]bool)
	for _, v := range a {
		set[v] = true
	}
	var intersection []string
	for _, v := range b {
		if set[v] {
			intersection = append(intersection, v)
		}
	}
	return intersection
}

func SliceUnion(a, b []string) []string {
	set := make(map[string]bool)
	for _, v := range a {
		set[v] = true
	}
	var intersection []string
	for _, v := range b {
		if !set[v] {
			intersection = append(intersection, v)
		}
	}
	return intersection
}

func SliceExist(list []string, str string) bool {
	set := make(map[string]bool)
	for _, v := range list {
		set[v] = true
	}
	return set[str]
}
