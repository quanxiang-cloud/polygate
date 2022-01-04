package header

// DeepCopy copy the first none-empty header value
func DeepCopy(src []string) string {
	for _, elem := range src {
		if elem != "" {
			return elem
		}
	}
	return ""
}
