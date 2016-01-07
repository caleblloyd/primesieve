package primesieve

import (
	"container/list"
	"container/ring"
	"runtime"
)

const PRIME_GROUP_SIZE = 65536
const TRY_SIZE = 2048

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
	Start      int
	End        int
	capped     bool
	primes     []int
	primesList *list.List
}

func NewPrimeGroup() *PrimeGroup {
	return &PrimeGroup{
		primesList: list.New(),
	}
}

func (pg *PrimeGroup) Add(prime int) bool {
	added := true
	if !pg.capped {
		if pg.primesList.Len() < PRIME_GROUP_SIZE-1 {
			if pg.primesList.Len() == 0 {
				pg.Start = prime
			}
			pg.primesList.PushBack(prime)
		} else {
			pg.End = prime
			pg.primesList.PushBack(prime)
			pg.primes = pg.listInternal()
			pg.primesList = nil
			pg.capped = true
		}
	} else {
		added = false
	}
	return added
}

func (pg *PrimeGroup) listInternal() []int {
	primes := make([]int, pg.primesList.Len())
	i := 0
	for e := pg.primesList.Front(); e != nil; e = e.Next() {
		primes[i] = e.Value.(int)
		i++
	}
	return primes
}

func (pg *PrimeGroup) List() []int {
	if !pg.capped && len(pg.primes) != pg.primesList.Len() {
		pg.primes = pg.listInternal()
	}
	return pg.primes
}

func (pg *PrimeGroup) Compare(tryPrimes []int, passedList *list.List, doneCh chan struct{}) {
	for _, try := range tryPrimes {
		pass := true
		var lastPrime int
		for _, prime := range pg.List() {
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
			passedList.PushBack(try)
		}
	}
	doneCh <- struct{}{}
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
		generate:    TRY_SIZE * threads,
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
	curGap := wheel2357()
	pgl.Add(2)
	pgl.Add(3)
	pgl.Add(5)
	pgl.Add(7)
	pgl.Add(11)
	gapTotal := 11
	for true {
		// generate potential primes
		tryList := list.New()
		max := pgl.End * 2
		for true {
			gapVal := curGap.Value.(int)
			gapTotal += gapVal
			if gapTotal > max {
				gapTotal -= gapVal
				break
			}
			tryList.PushBack(gapTotal)
			curGap = curGap.Next()
			if tryList.Len() >= pgl.generate {
				break
			}
		}
		// compare
		pg := pgl.primeGroups.Front()
		for true {

			nextTry := list.New()
			tryLen := tryList.Len()
			tryLenOrig := tryLen
			tryPtr := tryList.Front()
			for true {
				whole := tryLen / pgl.threads
				remain := tryLen % pgl.threads
				if whole > TRY_SIZE {
					whole = TRY_SIZE
					remain = 0
					tryLen -= TRY_SIZE * pgl.threads
				} else {
					tryLen = 0
				}

				passedLists := make([]*list.List, pgl.threads)
				for t := 0; t < pgl.threads; t++ {
					passedLists[t] = list.New()
					size := whole
					if remain > 0 {
						size++
						remain--
					}
					try := make([]int, size)
					for tryI := 0; tryI < size; tryI++ {
						try[tryI] = tryPtr.Value.(int)
						tryPtr = tryPtr.Next()
					}
					go pg.Value.(*PrimeGroup).Compare(try, passedLists[t], ch)
				}
				for t := 0; t < pgl.threads; t++ {
					<-ch
				}
				for _, pl := range passedLists {
					for p := pl.Front(); p != nil; p = p.Next() {
						nextTry.PushBack(p.Value.(int))
					}
				}
				if tryLen == 0 {
					break
				}
			}

			pg = pg.Next()
			if pg == nil {
				// End of iteration, everything left is primes
				for p := nextTry.Front(); p != nil; p = p.Next() {
					pgl.Add(p.Value.(int))
				}
				pgl.generate = TRY_SIZE * pgl.threads * (tryLenOrig / nextTry.Len())
				break
			}
			tryList = nextTry
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
