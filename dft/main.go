package main

import (
	"fmt"
	"math"
)

type Complex struct {
	Re float64
	Im float64
}

func formatC(z Complex) string {
	return fmt.Sprintf("%.4f%+.4fi", z.Re, z.Im)
}

// (a+bi) + (c+di) = (a+c) + (b+d)i
func add(a, b Complex) Complex {
	return Complex{
		Re: a.Re + b.Re,
		Im: a.Im + b.Im,
	}
}

// (a+bi)(c+di) = (ac - bd) + (ad + bc)i
func mul(a, b Complex) Complex {
	return Complex{
		Re: a.Re*b.Re - a.Im*b.Im,
		Im: a.Re*b.Im + a.Im*b.Re,
	}
}

// (a+bi) / n = (a/n) + (b/n)i
func divReal(z Complex, n float64) Complex {
	return Complex{
		Re: z.Re / n,
		Im: z.Im / n,
	}
}

// exp(iθ) = cos(θ) + i sin(θ)
func expI(theta float64) Complex {
	return Complex{
		Re: math.Cos(theta),
		Im: math.Sin(theta),
	}
}

// Прямое ДПФ (2.1):
// A_k = (1/N) * Σ_{j=0}^{N-1} exp(-2πi * k*j / N) * f_j
func dftForward(f []Complex) []Complex {
	n := len(f)
	A := make([]Complex, n)
	ops := 0

	for k := 0; k < n; k++ {
		sum := Complex{Re: 0, Im: 0}
		fmt.Printf("A[%d]:\n", k)

		for j := 0; j < n; j++ {
			theta := -2 * math.Pi * float64(k*j) / float64(n)
			w := expI(theta)
			term := mul(w, f[j]) // w*f_j
			sum = add(sum, term) // суммы произведений w_j*fj
			ops += 5

			fmt.Printf("  w[%d] = %s\n", j, formatC(w))


		}

		A[k] = divReal(sum, float64(n))
		fmt.Printf("  = %s\n\n", formatC(A[k]))
	}

	fmt.Printf("Трудоёмкость прямого ДПФ: C*N² = 5*%d² = %d\n", n, ops)
	return A
}

// Обратное ДПФ (2.2):
// f_k = Σ_{j=0}^{N-1} exp(2πi * k*j / N) * A_j
func dftInverse(A []Complex) []Complex {
	n := len(A)
	f := make([]Complex, n)
	ops := 0

	for k := 0; k < n; k++ {
		sum := Complex{Re: 0, Im: 0}
		fmt.Printf("f[%d]:\n", k)

		for j := 0; j < n; j++ {
			theta := 2 * math.Pi * float64(k*j) / float64(n)
			w := expI(theta)
			term := mul(w, A[j])
			sum = add(sum, term)
			ops += 5

			fmt.Printf("  w[%d] = %s\n", j, formatC(w))

		}

		f[k] = sum
		fmt.Printf("  = %s\n\n", formatC(f[k]))
	}

	fmt.Printf("Трудоёмкость обратного ДПФ: C*N² = 5*%d² = %d\n", n, ops)
	return f
}

func main() {
	src := []float64{1, 2, 3, 4}
	f := make([]Complex, len(src))
	for i, v := range src {
		f[i] = Complex{Re: v, Im: 0}
	}

	fmt.Println("Исходный массив f:")
	for i, v := range f {
		fmt.Printf("f[%d] = %s\n", i, formatC(v))
	}

	fmt.Println("\nПрямое ДПФ:")
	A := dftForward(f)

	fmt.Println("\nОбратное ДПФ:")
	restored := dftInverse(A)

	fmt.Println("\nВосстановленный массив:")
	for i, v := range restored {
		fmt.Printf("f[%d] = %s\n", i, formatC(v))
	}
}
