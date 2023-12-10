package main

func transpose[T any](data [][]T, outChan chan *T, done chan struct{}) {
	indices := make([]int, len(data))
	// how many rows to go through
	remaining := len(data)

	i := 0
	for remaining != 0 {
		row := data[i]
		if indices[i] == len(row) {
			remaining--
			indices[i]++
		} else if indices[i] < len(row) {
			outChan <- &row[indices[i]]
			indices[i]++
		}
		i = (i + 1) % len(data)
	}
	done <- struct{}{}
}
