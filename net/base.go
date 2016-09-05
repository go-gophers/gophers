package net

func errorLabelValue(err error) string {
	if err == nil {
		return "0"
	}
	return "1"
}
