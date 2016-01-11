package primesieve

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func BenchmarkPrintPrimes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ListN(1000)
	}
}

func TestSumPrimes1000(t *testing.T) {
	const PSUM = 3682913
	const THOUSANDTH_PRIME = 7919
	var calcsum uint64
	primes := ListN(1000)
	for _, p := range primes {
		calcsum += p
	}
	fmt.Println(primes[len(primes)-1])
	if PSUM != calcsum {
		t.Errorf("expected %d sum of ListN(1000), got %d\n", PSUM, calcsum)
	}

	primes = ListMax(THOUSANDTH_PRIME)
	calcsum = 0
	for _, p := range primes {
		calcsum += p
	}
	fmt.Println(primes[len(primes)-1])
	if PSUM != calcsum {
		t.Errorf("expected %d sum of ListMax(1000), got %d\n", PSUM, calcsum)
	}

	tp := PrimeN(1000)
	fmt.Println(tp)
	if tp != THOUSANDTH_PRIME {
		t.Errorf("expected %d for PrimeN(1000), got %d\n", THOUSANDTH_PRIME, tp)
	}
}

func TestPerformance(t *testing.T) {
	bm := func(max uint64, s string, expected uint64) {
		start := time.Now()
		c := Count(max)
		elapsed := time.Since(start)
		fmt.Printf("%s count: %d time: %s", s, c, elapsed)
		fmt.Println()
		if c != expected {
			t.Errorf("expected %d for Count(%d), got %d\n", expected, max, c)
		}
	}

	bm(uint64(math.Pow(10, 7)), "10^7", 664579)
	bm(uint64(math.Pow(10, 8)), "10^8", 5761455)
	bm(uint64(math.Pow(10, 9)), "10^9", 50847534)
	bm(uint64(math.Pow(2, 31))-1, "2^31-1", 105097565)
	bm(uint64(math.Pow(2, 32)), "2^32", 203280221)
}
