package engine

func unique(arr []string) *[]string {
	occurred := map[string]bool{}
	var result []string

	for e := range arr {
		if !occurred[arr[e]] {
			occurred[arr[e]] = true
			result = append(result, arr[e])
		}
	}
	return &result
}
