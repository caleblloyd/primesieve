# primesieve

A sieve of eratosthenes with wheel factorization optimization prime number generator, written in Go

This is pure Go implementation of the algorithm outlined at http://primesieve.org/segmented_sieve.html

Usage
-----

	package main

	import(
		"fmt"
		"github.com/caleblloyd/primesieve"
	)

	func main(){
		fmt.Println(primesieve.ListN(10))
		// [2 3 5 7 11 13 17 19 23 29]

		fmt.Println(primesieve.ListMax(19))
		// [2 3 5 7 11 13 17 19]

		fmt.Println(primesieve.PrimeN(10))
		// 19

		fmt.Println(primesieve.PrimeMax(25))
		// 23

		fmt.Println(primesieve.Count(29))
		// 10

		c := make(chan int)
		primesieve.Channel(c)
		for i:=0; i<10; i++{
			fmt.Printf("%d ", <-c)
		}
		fmt.Println()
		// 2 3 5 7 11 13 17 19 23 29
	}

Performance
-----------

On a 2015 Macbook Pro Retina running Go 1.5, calling primesieve.Count(x):

| x      	| Prime Count 	|  Time 	|
|--------	|-------------	|------:	|
| 10^7   	| 664,579     	| 0.03s 	|
| 10^8   	| 5,761,455   	| 0.3s  	|
| 10^9   	| 50,847,534  	| 3.2s  	|
| 2^31-1 	| 105,097,565 	| 7.3s  	|
| 2^32   	| 203,280,221 	| 14.8s 	|


Contributions
-------------

I welcome any performance improving Pull Requests or suggestions in the Issues section

Originally inspired by https://github.com/rlmcpherson/primesieve, credit goes out to
them for the wheel2357 factorization. Revised to use Segmented sieve of
Eratosthenes algorithm outlined at http://primesieve.org/segmented_sieve.html
