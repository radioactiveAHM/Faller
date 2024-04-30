package main

func MergeMap[T any](A map[string]T, B map[string]T) {
	for h, v := range B {
		A[h] = v
	}
}
