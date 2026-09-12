package seg

import "sync"

type chunkJob struct {
	i int
	s string
}

type chunkResult struct {
	i     int
	sents []string
}

func mapChunks(chunks []string, workers int, fn func(string) []string) []string {
	if len(chunks) == 0 {
		return []string{}
	}
	if len(chunks) == 1 || workers < 2 {
		return fn(chunks[0])
	}
	if workers > len(chunks) {
		workers = len(chunks)
	}

	jobs := make(chan chunkJob)
	results := make(chan chunkResult, len(chunks))
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for j := range jobs {
				results <- chunkResult{i: j.i, sents: fn(j.s)}
			}
		}()
	}
	for i, c := range chunks {
		jobs <- chunkJob{i: i, s: c}
	}
	close(jobs)
	wg.Wait()
	close(results)

	parts := make([][]string, len(chunks))
	for r := range results {
		parts[r.i] = r.sents
	}
	n := 0
	for _, p := range parts {
		n += len(p)
	}
	out := make([]string, 0, n)
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}
