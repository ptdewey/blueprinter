package utils

func MergeMaps(maps ...map[string]any) map[string]any {
	merged := make(map[string]any)

	for _, m := range maps {
		for key, value := range m {
			merged[key] = value
		}
	}

	return merged
}
