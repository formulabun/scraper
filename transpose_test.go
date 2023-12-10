package main

import "testing"

func TestTranspose(t *testing.T) {
	testData := [][]int{
		{1, 4},
		{2, 5, 6},
		{},
		{},
		{3},
	}
	outChan := make(chan *int, 6)
	doneChan := make(chan struct{}, 1)
	defer close(outChan)
	defer close(doneChan)
	transpose(testData, outChan, doneChan)
	i := 1
	for i <= 6 {
		select {
		case next := <-outChan:
			if *next != i {
				t.Errorf("expected %d but got %d", i, *next)
			}
			i++
		default:
			t.Fatalf("Expected transpose to block to completion.")
		}
	}

	select {
	case <-doneChan:
		break
	default:
		t.Fatalf("Transpose did not signal it is done")
	}
}
