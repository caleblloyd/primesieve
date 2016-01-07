package primesieve

import (
	"fmt"
	"testing"
)

func BenchmarkPrintPrimes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenPrimes(1000)
	}
}

func TestSumPrimes1000(t *testing.T) {
	const PSUM = 3682913
	const THOUSANDTH_PRIME = 7919
	var calcsum int
	primes := GenPrimes(1000)
	for _, p := range primes {
		calcsum += p
	}
	fmt.Println(primes[len(primes)-1])
	if PSUM != calcsum {
		t.Errorf("expected %d sum of GenPrimes(1000), got %d\n", PSUM, calcsum)
	}

	primes = GenPrimesMax(THOUSANDTH_PRIME)
	calcsum = 0
	for _, p := range primes {
		calcsum += p
	}
	fmt.Println(primes[len(primes)-1])
	if PSUM != calcsum {
		t.Errorf("expected %d sum of MaxGenPrimes(1000), got %d\n", PSUM, calcsum)
	}

	tp := NthPrime(1000)
	fmt.Println(tp)
	if tp != THOUSANDTH_PRIME {
		t.Errorf("expected %d for NthGetPrime(1000), got %d\n", THOUSANDTH_PRIME, tp)
	}
}