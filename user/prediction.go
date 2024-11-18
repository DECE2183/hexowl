package user

func Predict(word string) string {
	for k := range variables {
		if len(k) < len(word) {
			continue
		}
		if k[:len(word)] == word {
			return k
		}
	}

	for k := range functions {
		if len(k) < len(word) {
			continue
		}
		if k[:len(word)] == word {
			return k + "()"
		}
	}

	return ""
}
