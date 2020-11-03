package main

import "math"

type MyInt int

type sieveElement struct {
	number  int
	isPrime bool
}

func (m MyInt) IsEven() bool {
	return m%2 == 0
}

func (m MyInt) IsOdd() bool {
	return !m.IsEven()
}

func (m MyInt) IsPrime() bool {

	if m <= 3 {
		return true
	}

	s := make([]sieveElement, m-1, m-1)
	v := 2

	// Init Sieve
	for i := range s {
		s[i] = sieveElement{
			v,
			true,
		}

		v++
	}

	n := int(math.Ceil(math.Sqrt(float64(m))))
	for i := 0 ; i <= n ; i++ {
		current := s[i]
		for j := i + 1; j < len(s); j++ {
			if s[j].number%current.number == 0 {
				s[j].isPrime = false
			}
		}
	}

	return s[m-2].isPrime
}
