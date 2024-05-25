package utils

func Filter(slice []string, f func(string) bool) []string {
	var result []string
	for _, value := range slice {
		if f(value) {
			result = append(result, value)
		}
	}
	return result
}

func NotEmpty(s string) bool {
	return s != ""
}
