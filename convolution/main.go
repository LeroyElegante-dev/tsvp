package main

import (
	"fmt"
	"math"
	"strings"
)

const Pi = 3.1415

type Complex struct {
	Re float64
	Im float64
}

func add(a, b Complex) Complex {
	return Complex{Re: a.Re + b.Re, Im: a.Im + b.Im}
}

func sub(a, b Complex) Complex {
	return Complex{Re: a.Re - b.Re, Im: a.Im - b.Im}
}

func mul(a, b Complex) Complex {
	return Complex{
		Re: a.Re*b.Re - a.Im*b.Im,
		Im: a.Re*b.Im + a.Im*b.Re,
	}
}

func mulReal(z Complex, x float64) Complex {
	return Complex{Re: z.Re * x, Im: z.Im * x}
}

func expI(theta float64) Complex {
	return Complex{Re: math.Cos(theta), Im: math.Sin(theta)}
}

func abs(z Complex) float64 {
	return math.Hypot(z.Re, z.Im)
}

func maxAbsDiff(a, b []Complex) float64 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	m := 0.0
	for i := 0; i < n; i++ {
		d := abs(sub(a[i], b[i]))
		if d > m {
			m = d
		}
	}
	return m
}

func toComplex(xs []float64) []Complex {
	out := make([]Complex, len(xs))
	for i, v := range xs {
		out[i] = Complex{Re: v}
	}
	return out
}

func padZeros(x []Complex, m int) []Complex {
	out := make([]Complex, m)
	copy(out, x)
	return out
}

func printHeader(title string) {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("  %s\n", title)
	fmt.Println(strings.Repeat("=", 80))
}

func printVec(name string, x []Complex) {
	fmt.Printf("  %s = [\n", name)
	for i, z := range x {
		fmt.Printf("    [%2d] = %8.4f + %8.4f*i\n", i, z.Re, z.Im)
	}
	fmt.Println("  ]")
}

func printSpectrum(name string, x []Complex) {
	fmt.Printf("  %s = [\n", name)
	for i, z := range x {
		fmt.Printf("    [%2d] = %8.4f%+.4fi         |C| = %8.4f\n", i, z.Re, z.Im, abs(z))
	}
	fmt.Println("  ]")
}

func convSimple(a, b []Complex) []Complex {
	n := len(a)
	c := make([]Complex, 2*n-1)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			c[i+j] = add(c[i+j], mul(a[i], b[j]))
		}
	}
	return c
}

func dftForward(f []Complex) []Complex {
	n := len(f)
	A := make([]Complex, n)
	invN := 1.0 / float64(n)
	for k := 0; k < n; k++ {
		sum := Complex{}
		for j := 0; j < n; j++ {
			w := expI(-2 * Pi * float64(k*j) / float64(n))
			sum = add(sum, mul(w, f[j]))
		}
		A[k] = mulReal(sum, invN)
	}
	return A
}

func dftInverse(A []Complex) []Complex {
	n := len(A)
	f := make([]Complex, n)
	for k := 0; k < n; k++ {
		sum := Complex{}
		for j := 0; j < n; j++ {
			w := expI(2 * Pi * float64(k*j) / float64(n))
			sum = add(sum, mul(w, A[j]))
		}
		f[k] = sum
	}
	return f
}

func convDFT(a, b []Complex) (c, aTilde, bTilde, cTilde []Complex) {
	n := len(a)
	m := 2 * n
	ap := padZeros(a, m)
	bp := padZeros(b, m)
	aTilde = dftForward(ap)
	bTilde = dftForward(bp)
	cTilde = make([]Complex, m)
	scale := float64(m)
	for i := 0; i < m; i++ {
		cTilde[i] = mulReal(mul(aTilde[i], bTilde[i]), scale)
	}
	c = dftInverse(cTilde)
	return
}

func printTriangle(n int) {
	fmt.Println("Порядок образования слагаемых (демонстрация треугольной схемы):")
	last := 2*n - 2
	for k := 0; k <= last; k++ {
		terms := make([]string, 0, n)
		for i := 0; i < n; i++ {
			j := k - i
			if j >= 0 && j < n {
				terms = append(terms, fmt.Sprintf("a[%d]*b[%d]", i, j))
			}
		}
		fmt.Printf("  c_%d = %s (%d слагаемых)\n", k, strings.Join(terms, " + "), len(terms))
	}
}

func main() {
	a := toComplex([]float64{1, 2, 3, 4, 3, 2, 1, 0})
	b := toComplex([]float64{2, 1, 0, 1, 2, 3, 2, 1})
	n := len(a)

	fmt.Printf("Входные векторы (n = %d):\n", n)
	printVec("a", a)
	printVec("b", b)
	fmt.Println()

	printHeader("ЛАБОРАТОРНАЯ РАБОТА 7: ПРОСТАЯ СВЕРТКА (ПРЯМОЙ АЛГОРИТМ) - O(N^2)")
	cSimple := convSimple(a, b)
	fmt.Printf("Результат свертки c = a * b (длина 2n - 1 = %d):\n", 2*n-1)
	printVec("c_simple", cSimple)
	fmt.Printf("Количество операций умножения: %d (теоретически n^2 = %d)\n\n", n*n, n*n)
	printTriangle(n)
	fmt.Println()

	printHeader("ЛАБОРАТОРНАЯ РАБОТА 8: СВЕРТКА ЧЕРЕЗ ПРЕОБРАЗОВАНИЯ ФУРЬЕ (ДПФ) - O(N^2)")
	m := 2 * n
	tDFT := 3*m*m + m
	fmt.Printf("Массивы дополнены n нулями до длины 2n = %d.\n", m)
	fmt.Printf("Теоретическая трудоемкость: 3 * (2n)^2 + 2n = 3 * %d^2 + %d = %d операций.\n\n", m, m, tDFT)

	cDFT, aTilde, bTilde, cTilde := convDFT(a, b)
	fmt.Println("Спектр A_tilde (ДПФ вектора a, дополненного нулями):")
	printSpectrum("A_tilde", aTilde)
	fmt.Println("Спектр B_tilde (ДПФ вектора b, дополненного нулями):")
	printSpectrum("B_tilde", bTilde)
	fmt.Println("Спектр произведения C_tilde = 2n * A_tilde * B_tilde:")
	printSpectrum("C_tilde", cTilde)
	fmt.Println("Восстановленный вектор свертки c_dft = InverseDft(C_tilde):")
	printVec("c_dft", cDFT)
	fmt.Printf("Максимальная погрешность |c_прост - c_ДПФ|: %E\n", maxAbsDiff(cSimple, cDFT))
}
