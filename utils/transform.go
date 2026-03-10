package utils

func MapStrings(input []string, f func(string) string) []string {
	result := make([]string, len(input))

	for i, s := range input {
		result[i] = f(s)
	}

	return result
}

func MapStringsE(input []string, f func(string) (string, error)) ([]string, error) {
	result := make([]string, len(input))

	for i, s := range input {
		r, err := f(s)
		if err != nil {
			return nil, err
		}
		result[i] = r
	}

	return result, nil
}
