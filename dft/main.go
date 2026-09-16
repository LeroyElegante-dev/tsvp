package main

import (
	"fmt"
	"math"
	"math/cmplx"
)

func formatC(z complex128) string {
	return fmt.Sprintf("%.4f%+.4fi", real(z), imag(z))
}

// Прямое ДПФ (2.1):
// A_k = (1/N) * Σ_{j=0}^{N-1} exp(-2πi * k*j / N) * f_j
func dftForward(f []complex128) []complex128 {
	n := len(f)
	A := make([]complex128, n)
	ops := 0

	for k := 0; k < n; k++ {
		var sum complex128
		for j := 0; j < n; j++ {
			theta := -2 * math.Pi * float64(k*j) / float64(n)
			w := cmplx.Exp(complex(0, theta))
			sum += w * f[j]
			ops += 5
		}
		A[k] = sum / complex(float64(n), 0)
		fmt.Printf("A[%d] = %s\n", k, formatC(A[k]))
	}

	fmt.Printf("Трудоёмкость прямого ДПФ: C*N² = 5*%d² = %d\n", n, ops)
	return A
}

// Обратное ДПФ (2.2):
// f_k = Σ_{j=0}^{N-1} exp(2πi * k*j / N) * A_j
func dftInverse(A []complex128) []complex128 {
	n := len(A)
	f := make([]complex128, n)
	ops := 0

	for k := 0; k < n; k++ {
		var sum complex128
		for j := 0; j < n; j++ {
			theta := 2 * math.Pi * float64(k*j) / float64(n)
			w := cmplx.Exp(complex(0, theta))
			sum += w * A[j]
			ops += 5
		}
		f[k] = sum
		fmt.Printf("f[%d] = %s\n", k, formatC(f[k]))
	}

	fmt.Printf("Трудоёмкость обратного ДПФ: C*N² = 5*%d² = %d\n", n, ops)
	return f
}

func main() {
	src := []float64{1, 2, 3, 4}
	f := make([]complex128, len(src))
	for i, v := range src {
		f[i] = complex(v, 0)
	}

	fmt.Println("Исходный массив f:")
	for i, v := range f {
		fmt.Printf("f[%d] = %s\n", i, formatC(v))
	}

	fmt.Println("\nПрямое ДПФ (формула 2.1):")
	A := dftForward(f)

	fmt.Println("\nОбратное ДПФ (формула 2.2):")
	restored := dftInverse(A)

	fmt.Println("\nВосстановленный массив:")
	for i, v := range restored {
		fmt.Printf("f[%d] = %s\n", i, formatC(v))
	}
}
