# primesieve

A sieve of eratosthenes with wheel factorization optimization prime number generator, written in Go

I've attempted some performance improvements by attempting to keep loops small enough to fit into the CPU Cache

Usage
-----

	package main

	import(
		"fmt"
		"github.com/caleblloyd/primesieve"
	)

	func main(){
		fmt.Println(primesieve.GenPrimes(10))
		// [2 3 5 7 11 13 17 19 23 29]
		fmt.Println(primesieve.GenPrimesMax(19))
		// [2 3 5 7 11 13 17 19]
		fmt.Println(primesieve.NthPrime(10))
		// 29
	}

Performance
-----------

Performance is decent, but could definitely be better.  On a 2015 Macbook Pro Retina running Go 1.5, performance was:

- 1,000 primes in ~15ms
- 10,000 primes in ~35ms
- 100,000 primes in ~275ms
- 1,000,000 primes in ~3.75s


Contributions
-------------

I welcome any performance improving Pull Requests or suggestions in the Issues section

Originally inspired by https://github.com/rlmcpherson/primesieve, credit goes out to them for the wheel2357 factorization function also
