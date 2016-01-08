package primesieve

import (
	"math"
)

/// Set your CPUs L1 data cache size (in bytes) here
const SEGMENT_SIZE = 32768

const (
	RTYPE_COUNT = iota
	RTYPE_LIST_N
	RTYPE_LIST_MAX
	RTYPE_MAX
	RTYPE_N
	RTYPE_CHANNEL
)

/// Generate primes using the segmented sieve of Eratosthenes.
/// This algorithm uses O(n log log n) operations and O(sqrt(n)) space.
func SegmentedSieve(rType int, rLimit int, rChan chan int) (int, []int) {
	var rList []int
	var rInt int

	var limit int
	if rType == RTYPE_LIST_N || rType == RTYPE_N || rType == RTYPE_CHANNEL {
		limit = math.MaxInt32
	} else {
		limit = rLimit
	}

	var count int
	// anonymous fnction to handle different return types
	foundContinue := func(n int) bool {
		count++
		if rType == RTYPE_COUNT {
			rInt = count
		} else if rType == RTYPE_LIST_MAX || rType == RTYPE_LIST_N {
			rList = append(rList, n)
		} else if rType == RTYPE_N || rType == RTYPE_MAX {
			rInt = n
		} else if rType == RTYPE_CHANNEL {
			rChan <- n
		}
		if (rType == RTYPE_LIST_N || rType == RTYPE_N) && count >= rLimit {
			return false
		}
		return true
	}

	if limit < 2 || !foundContinue(2) {
		return rInt, rList
	}

	sqrt := int(math.Sqrt(float64(limit)))
	var sieve [SEGMENT_SIZE]bool

	// generate small primes <= sqrt
	var isComposite = make([]bool, sqrt+1)
	for i := 2; i*i <= sqrt; i++ {
		if !isComposite[i] {
			for j := i * i; j <= sqrt; j += i {
				isComposite[j] = true
			}
		}
	}

	s := 1
	n := 3
	var gaps2357 = []int{2, 4, 2, 4, 6, 2, 6, 4, 2, 4, 6, 6, 2, 6, 4, 2, 6, 4, 6, 8, 4, 2, 4, 2, 4, 8, 6, 4, 6, 2, 4, 6, 2, 6, 6, 4, 2, 4, 6, 2, 6, 4, 2, 4, 2, 10, 2, 10}
	var gapPos int
	var primes, next []int
	var stopIteration bool
	gaps2357Len := len(gaps2357)

	for low := 0; low <= limit; low += SEGMENT_SIZE {

		for i, _ := range sieve {
			sieve[i] = true
		}

		// current segment = interval [low, high]
		high := low + SEGMENT_SIZE
		if high > limit {
			high = limit
		}

		// store small primes needed to cross off multiples
		for ; s*s <= high; s++ {
			if !isComposite[s] {
				primes = append(primes, s)
				next = append(next, s*s-low)
			}
		}

		// sieve the current segment
		for i := 1; i < len(primes); i++ {
			j := next[i]
			for k := primes[i] * 2; j < SEGMENT_SIZE; j += k {
				sieve[j] = false
			}
			next[i] = j - SEGMENT_SIZE
		}

		for true {
			if n > high {
				break
			}
			if sieve[n-low] {
				// n is a prime
				if !foundContinue(n) {
					stopIteration = true
					break
				}
			}
			// wheel factorization optimization
			if n >= 11 {
				n += gaps2357[gapPos]
				gapPos++
				if gapPos >= gaps2357Len {
					gapPos = 0
				}
			} else if n == 3 {
				n = 5
			} else if n == 5 {
				n = 7
			} else {
				n = 11
			}
		}

		if stopIteration {
			break
		}
	}

	return rInt, rList

}

func ListN(numPrimes int) []int {
	_, l := SegmentedSieve(RTYPE_LIST_N, numPrimes, nil)
	return l
}

func ListMax(max int) []int {
	_, l := SegmentedSieve(RTYPE_LIST_MAX, max, nil)
	return l
}

func Count(max int) int {
	s, _ := SegmentedSieve(RTYPE_COUNT, max, nil)
	return s
}

func PrimeN(n int) int {
	s, _ := SegmentedSieve(RTYPE_N, n, nil)
	return s
}

func PrimeMax(max int) int {
	s, _ := SegmentedSieve(RTYPE_MAX, max, nil)
	return s
}

func Channel(c chan int) {
	go SegmentedSieve(RTYPE_CHANNEL, 0, c)
}
