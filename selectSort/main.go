package main

import "fmt"

func selectionSort(arr []int) {
	n := len(arr)
	iter := 0

	for i := 0; i < n-1; i++ {
		minIdx := i

		for j := i + 1; j < n; j++ {
			if arr[j] < arr[minIdx] {
				minIdx = j
			}
		}

		if minIdx != i {
			arr[i], arr[minIdx] = arr[minIdx], arr[i]
		}
		iter++
		fmt.Println("Итерация", iter, ":", arr)
	}
}

func main() {
	arr := []int{64, 25, 12, 22, 11}

	fmt.Println("Исходный массив:", arr)

	selectionSort(arr)

	fmt.Println("Отсортированный массив:", arr)
}
