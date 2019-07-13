package utils

func ContainsString(list []string, match string) bool {

	for _, l := range list {
		if l == match {
			return true
		}
	}

	return false
}
