package seg

const parallelMinBytes = 32 * 1024

// splitChunks divides s into around n pieces at blank lines so each worker
// gets a contiguous slice. The original bytes are not rewritten; chunks
// concatenate back to s.
func splitChunks(s string, n int) []string {
	if n < 2 || len(s) == 0 {
		return []string{s}
	}
	breaks := paragraphEnds(s)
	if len(breaks) < 2 {
		return []string{s}
	}
	target := (len(s) + n - 1) / n
	if target < 1 {
		target = 1
	}
	var chunks []string
	start := 0
	for _, end := range breaks {
		if end-start >= target && end > start {
			chunks = append(chunks, s[start:end])
			start = end
		}
	}
	if start < len(s) {
		chunks = append(chunks, s[start:])
	}
	if len(chunks) == 0 {
		return []string{s}
	}
	return chunks
}

// paragraphEnds returns byte indices after each blank-line run, plus len(s).
func paragraphEnds(s string) []int {
	var ends []int
	for i := 0; i < len(s); {
		if s[i] != '\n' {
			i++
			continue
		}
		j := i + 1
		for j < len(s) && (s[j] == ' ' || s[j] == '\t' || s[j] == '\r') {
			j++
		}
		if j < len(s) && s[j] == '\n' {
			j++
			for j < len(s) && (s[j] == '\n' || s[j] == ' ' || s[j] == '\t' || s[j] == '\r') {
				j++
			}
			if j > i {
				ends = append(ends, j)
			}
			i = j
			continue
		}
		i++
	}
	if len(ends) == 0 || ends[len(ends)-1] != len(s) {
		ends = append(ends, len(s))
	}
	return ends
}
