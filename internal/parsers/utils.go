package parsers

func calculateIndentSpacesCnt(indent string) int {
	cnt := 0
	for _, r := range indent {
		if r == '\t' {
			cnt += 4 // TODO: move tab size to config
		} else {
			cnt++
		}
	}
	return cnt
}
