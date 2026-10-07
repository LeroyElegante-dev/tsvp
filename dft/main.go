package main

import (
	"fmt"
	"math"
	"os"
)

const Pi = 3.1415

type Complex struct {
	Re float64
	Im float64
}

func formatC(z Complex) string {
	return fmt.Sprintf("%.7f%+.7fi", z.Re, z.Im)
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
func dftForward(f []Complex, file *os.File) []Complex {
	n := len(f)
	A := make([]Complex, n)
	ops := 0
	countMul := 0

	for k := 0; k < n; k++ {
		sum := Complex{Re: 0, Im: 0}
		// fmt.Printf("A[%d]:\n", k)
		fmt.Fprintf(file, "A[%d]: ", k)

		for j := 0; j < n; j++ {
			theta := -2 * Pi * float64(k*j) / float64(n)
			w := expI(theta)
			term := mul(w, f[j]) // w*f_j
			countMul += 1
			sum = add(sum, term) // суммы произведений w_j*fj
			ops += 5

			// fmt.Printf("  w[%d] = %s\n", j, formatC(w))
			// fmt.Fprintf(file, "  w[%d] = %s\n", j, formatC(w))

		}

		A[k] = divReal(sum, float64(n))
		// fmt.Printf("  = %s\n\n", formatC(A[k]))
		fmt.Fprintf(file, "%s\n\n", formatC(A[k]))
	}

	fmt.Fprintf(file, "Трудоёмкость прямого ДПФ: C*N² = 5*%d² = %d\n", n, ops)
	fmt.Fprintf(file, "CountMul: %d\n", countMul)

	return A
}

// Обратное ДПФ (2.2):
// f_k = Σ_{j=0}^{N-1} exp(2πi * k*j / N) * A_j
func dftInverse(A []Complex, file *os.File) []Complex {
	n := len(A)
	f := make([]Complex, n)
	ops := 0
	countMul := 0

	for k := 0; k < n; k++ {
		sum := Complex{Re: 0, Im: 0}
		fmt.Fprintf(file, "f[%d]: ", k)

		for j := 0; j < n; j++ {
			theta := 2 * Pi * float64(k*j) / float64(n)
			w := expI(theta)
			term := mul(w, A[j])
			countMul += 1
			sum = add(sum, term)
			ops += 5

			// fmt.Fprintf(file, "  w[%d] = %s\n", j, formatC(w))

		}

		f[k] = sum
		fmt.Fprintf(file, "%s\n\n", formatC(f[k]))
	}

	fmt.Fprintf(file, "Трудоёмкость обратного ДПФ: C*N² = 5*%d² = %d\n", n, ops)
	fmt.Fprintf(file, "CountMul: %d\n", countMul)

	return f
}

func runDft(size int) {
	filename := fmt.Sprintf("dft_output_%d.txt", size)
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Ошибка создания файла: %v\n", err)
		return
	}
	defer file.Close()

	src := make([]float64, size)
	for i := 0; i < size; i++ {
		src[i] = float64(i)
	}

	f := make([]Complex, len(src))
	for i, v := range src {
		f[i] = Complex{Re: v, Im: 0}
	}

	fmt.Fprintln(file, "Исходный массив f:")
	for i, v := range f {
		fmt.Fprintf(file, "f[%d] = %s\n", i, formatC(v))
	}
	fmt.Fprintln(file, "")

	fmt.Fprintln(file, "\nПрямое ДПФ:")
	A := dftForward(f, file)

	fmt.Fprintln(file, "\nОбратное ДПФ:")
	_ = dftInverse(A, file)

}

func main() {
	runDft(10)
	runDft(100)
	runDft(400)

	// src := []float64{1, 2, 3, 4}
	// f := make([]Complex, len(src))
	// for i, v := range src {
	// 	f[i] = Complex{Re: v, Im: 0}
	// }

	// fmt.Println("Исходный массив f:")
	// for i, v := range f {
	// 	fmt.Printf("f[%d] = %s\n", i, formatC(v))
	// }

	// fmt.Println("\nПрямое ДПФ:")
	// A := dftForward(f)

	// fmt.Println("\nОбратное ДПФ:")
	// restored := dftInverse(A)

	// fmt.Println("\nВосстановленный массив:")
	// for i, v := range restored {
	// 	fmt.Printf("f[%d] = %s\n", i, formatC(v))
	// }
}