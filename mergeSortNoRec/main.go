package main

import "fmt"

func merge(a []int, l, m, r int, buf []int) {
	i, j, k := l, m, l
	for i < m && j < r {
		if a[i] <= a[j] {
			buf[k] = a[i]
			i++
		} else {
			buf[k] = a[j]
			j++
		}
		k++
	}
	for i < m {
		buf[k] = a[i]
		i++
		k++
	}
	for j < r {
		buf[k] = a[j]
		j++
		k++
	}
	copy(a[l:r], buf[l:r])
}

func mergeSort(a []int) {
	n := len(a)
	buf := make([]int, n)
	iter := 0

	for width := 1; width < n; width *= 2 {
		for l := 0; l < n; l += 2 * width {
			m := l + width
			if m > n {
				m = n
			}
			r := l + 2*width
			if r > n {
				r = n
			}
			if m < r {
				merge(a, l, m, r, buf)
				iter++
				fmt.Println("Итерация", iter, ":", a)
			}
		}
	}
}

func main() {
	arr := []int{64, 25, 12, 22, 11}

	fmt.Println("Исходный массив:", arr)

	mergeSort(arr)

	fmt.Println("Отсортированный массив:", arr)
}
