package primesieve

import (
	"container/list"
	"container/ring"
	"runtime"
)

const PRIME_GROUP_SIZE = 100000
const TRY_PER_THREAD_SIZE = 100000

func wheel2357() *ring.Ring {
	var gaps2357 = []int{2, 4, 2, 4, 6, 2, 6, 4, 2, 4, 6, 6, 2, 6, 4, 2, 6, 4, 6, 8, 4, 2, 4, 2, 4, 8, 6, 4, 6, 2, 4, 6, 2, 6, 6, 4, 2, 4, 6, 2, 6, 4, 2, 4, 2, 10, 2, 10}
	r := ring.New(len(gaps2357))
	for _, i := range gaps2357 {
		r.Value = i
		r = r.Next()
	}
	return r
}

type PrimeGroup struct {
	Primes    []int
	PrimesLen int
	Capped    bool
}

func NewPrimeGroup() *PrimeGroup {
	return &PrimeGroup{
		Primes: make([]int, PRIME_GROUP_SIZE),
	}
}

func (pg *PrimeGroup) Add(prime int) bool {
	if !pg.Capped {
		pg.Primes[pg.PrimesLen] = prime
		pg.PrimesLen++
		if pg.PrimesLen >= PRIME_GROUP_SIZE {
			pg.Capped = true
		}
		return true
	}
	return false
}

func (pg *PrimeGroup) Compare(tg *TryGroup) {
	tryLen := tg.TryLen
	tg.Reset()
	for t := 0; t < tryLen; t++ {
		try := tg.Try[t]
		pass := true
		var lastPrime int
		for i := 0; i < pg.PrimesLen; i++ {
			prime := pg.Primes[i]
			if try%prime == 0 {
				pass = false
				break
			}
			if lastPrime*prime > try {
				break
			}
			lastPrime = prime
		}
		if pass {
			tg.Add(try)
		}
	}
}

type TryGroup struct {
	Try    []int
	TryLen int
}

func NewTryGroup() *TryGroup {
	return &TryGroup{
		Try: make([]int, TRY_PER_THREAD_SIZE),
	}
}

func (tg *TryGroup) Add(try int) bool {
	if tg.TryLen < TRY_PER_THREAD_SIZE {
		tg.Try[tg.TryLen] = try
		tg.TryLen++
		return true
	}
	return false
}

func (tg *TryGroup) Reset() {
	tg.TryLen = 0
}

type PrimeGroupList struct {
	End         int
	generate    int
	outCh       chan int
	primeGroups *list.List
	threads     int
}

func NewPrimeGroupList(outCh chan int) *PrimeGroupList {
	primeGroups := list.New()
	primeGroups.PushBack(NewPrimeGroup())
	threads := runtime.GOMAXPROCS(0)
	return &PrimeGroupList{
		outCh:       outCh,
		primeGroups: primeGroups,
		threads:     threads,
	}
}

func (pgl *PrimeGroupList) Add(prime int) {
	pg := pgl.primeGroups.Back().Value.(*PrimeGroup)
	if !pg.Add(prime) {
		pg := NewPrimeGroup()
		pg.Add(prime)
		pgl.primeGroups.PushBack(pg)
	}
	pgl.End = prime
	pgl.outCh <- prime
}

func (pgl *PrimeGroupList) Generate() {
	ch := make(chan struct{}, pgl.threads)
	tg := make([]*TryGroup, pgl.threads)
	for i := 0; i < pgl.threads; i++ {
		tg[i] = NewTryGroup()
	}
	curGap := wheel2357()
	// 2 is prime, but we won't generate factors of 2 so no need to add to pgl for comparison
	pgl.outCh <- 2
	// everything else gets added to pgl for comparison
	pgl.Add(3)
	pgl.Add(5)
	pgl.Add(7)
	pgl.Add(11)
	gapTotal := 11
	for true {
		// generate potential primes
		for i := 0; i < pgl.threads; i++ {
			tg[i].Reset()
		}
		max := pgl.End * 3
		tgi := 0
		for true {
			gapVal := curGap.Value.(int)
			gapTotal += gapVal
			if gapTotal > max {
				gapTotal -= gapVal
				break
			}
			if !tg[tgi].Add(gapTotal) {
				tgi++
				if tgi >= pgl.threads {
					gapTotal -= gapVal
					break
				}
				tg[tgi].Add(gapTotal)
			}
			curGap = curGap.Next()
		}

		f := func(ftg *TryGroup) {
			for pg := pgl.primeGroups.Front(); pg != nil; pg = pg.Next() {
				pg.Value.(*PrimeGroup).Compare(ftg)
			}
			ch <- struct{}{}
		}

		// compare
		wait := 0
		for i := 0; i < pgl.threads; i++ {
			if tg[i].TryLen > 0 {
				go f(tg[i])
				wait++
			}
		}
		for w := 0; w < wait; w++ {
			<-ch
		}
		for i := 0; i < pgl.threads; i++ {
			if tg[i].TryLen > 0 {
				for ti := 0; ti < tg[i].TryLen; ti++ {
					pgl.Add(tg[i].Try[ti])
				}
			}
		}
	}
}
func GenPrimes(numPrimes int) []int {
	ch := make(chan int, 16)
	pgl := NewPrimeGroupList(ch)
	primes := make([]int, numPrimes)
	go pgl.Generate()
	for i := 0; i < numPrimes; i++ {
		prime := <-ch
		primes[i] = prime
	}
	return primes
}

func GenPrimesMax(max int) []int {
	ch := make(chan int, 16)
	pgl := NewPrimeGroupList(ch)
	var primes []int
	go pgl.Generate()
	for true {
		prime := <-ch
		if prime > max {
			return primes
		}
		primes = append(primes, prime)
	}
	return primes
}

func NthPrime(n int) int {
	ch := make(chan int, 16)
	pgl := NewPrimeGroupList(ch)
	go pgl.Generate()
	var prime int
	for ct := 0; ct < n; ct++ {
		prime = <-ch
	}
	return prime
}
