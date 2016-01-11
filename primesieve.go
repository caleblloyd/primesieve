package primesieve

import (
	"math"
)

// Set your CPUs L1 data cache size (in bytes) here
const SEGMENT_SIZE = 32768

// wheel factorization optimization length
const WHEEL_LEN = 48

const (
	RTYPE_COUNT = iota
	RTYPE_LIST_N
	RTYPE_LIST_MAX
	RTYPE_MAX
	RTYPE_N
	RTYPE_CHANNEL
)

// minimum number of composites to dynamically allocate
const ALLOCATE_MIN = 1024

// Generate primes using the segmented sieve of Eratosthenes.
// This algorithm uses O(n log log n) operations and O(sqrt(n)) space.
func SegmentedSieve(rType uint8, rLimit uint64, rChan chan uint64) (uint64, []uint64) {
	var rList []uint64
	var rInt uint64
	var isComposite []bool

	var limit uint64
	if rType == RTYPE_LIST_N || rType == RTYPE_N || rType == RTYPE_CHANNEL {
		limit = math.MaxUint64
	} else {
		limit = rLimit
	}

	var count uint64
	// anonymous function to handle different return types
	foundContinue := func(n uint64) bool {
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

	var sieve [SEGMENT_SIZE]bool

	s := uint64(2)
	n := uint64(3)
	var wheel2357 = []uint8{2, 4, 2, 4, 6, 2, 6, 4, 2, 4, 6, 6, 2, 6, 4, 2, 6, 4, 6, 8, 4, 2, 4, 2, 4, 8, 6, 4, 6, 2, 4, 6, 2, 6, 6, 4, 2, 4, 6, 2, 6, 4, 2, 4, 2, 10, 2, 10}
	var wheelPos uint8
	var primes, next []uint64
	var stopIteration bool

	for low := uint64(0); low <= limit; low += SEGMENT_SIZE {

		for i, _ := range sieve {
			sieve[i] = true
		}

		// current segment = interval [low, high]
		high := low + SEGMENT_SIZE
		if high > limit {
			high = limit
		}

		sqrt := uint64(math.Sqrt(float64(high))) + 1
		if s < sqrt {
			// dynamically allocate isComposite
			oldLen := uint64(len(isComposite))
			allocate := int64(sqrt - oldLen)
			if allocate > 0 {
				if allocate < ALLOCATE_MIN {
					allocate = ALLOCATE_MIN
				}

				// backfill new allocation
				isComposite = append(isComposite, make([]bool, allocate)...)
				sqrt = uint64(len(isComposite))

				for i := uint64(0); i < uint64(len(primes)); i++ {
					p := primes[i] * primes[i]
					if p < oldLen {
						p = oldLen - (oldLen % primes[i]) + primes[i]
					}
					for ; p < sqrt; p += primes[i] {
						isComposite[p] = true
					}
				}
			}

			// find new small primes
			for ; s < sqrt; s++ {
				if !isComposite[s] {
					primes = append(primes, s)
					next = append(next, s*s-low)
					for j := s * s; j < sqrt; j += s {
						isComposite[j] = true
					}
				}
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
				n += uint64(wheel2357[wheelPos])
				wheelPos++
				if wheelPos >= WHEEL_LEN {
					wheelPos = 0
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

func ListN(numPrimes uint64) []uint64 {
	_, l := SegmentedSieve(RTYPE_LIST_N, numPrimes, nil)
	return l
}

func ListMax(max uint64) []uint64 {
	_, l := SegmentedSieve(RTYPE_LIST_MAX, max, nil)
	return l
}

func Count(max uint64) uint64 {
	s, _ := SegmentedSieve(RTYPE_COUNT, max, nil)
	return s
}

func PrimeN(n uint64) uint64 {
	s, _ := SegmentedSieve(RTYPE_N, n, nil)
	return s
}

func PrimeMax(max uint64) uint64 {
	s, _ := SegmentedSieve(RTYPE_MAX, max, nil)
	return s
}

func Channel(c chan uint64) {
	go SegmentedSieve(RTYPE_CHANNEL, 0, c)
}
