package root

import (
	"errors"
	"fmt"
	"math"
)

/*
const (
	Phi    = 1.618033988749894848204586834
	Resphi = 2 - Phi
)

// GoldenSection is a root-finding method. Assumes there's a unique minimum
// between minLoc and maxLoc
func GoldenSection(f func(float64) float64, minLoc, minVal, maxLoc, maxVal, tol float64) {

}

func boundedSearch(f func(float64) float64, a, fA, b, fB, c, fC, tol float64) {
	var newLoc float64
	if maxLoc-midLoc > midLoc-minLoc {
		newLoc = b + Resphi*(c-b)
	} else {
		newLoc = b - Resphi*(b-a)
	}
	if math.Abs(c-a) < tol*(math.Abs(b)+math.Abs(a)) {
		return (c + a) / 2
	}
}
*/

// Bisection finds the location between neg and pos where f(x) is within
// tol of zero. Neg is the location of a point where f(x) < 0 and pos
// is the location of a point where f(x) > 0.
// Neg and/or pos may be infinity. If exactly one of them is infinity,
// a search will be done in that direction to find a point where f(x)
// has the correct sign (for example, if pos is infinity, it will search
// until it finds a location where f(x)> 0). If both of them are negative
// it will search to find two points with opposite sign and then perform a
// bisection between those two points.
// If no point is found within 100 iterations, it will return an error.
// This could mean that the function is not continuous, the function is
// not deterministic, or could imply user error due to the function values
// at neg and pos having the same sign.
// If infs are used the code assumes that the function is monotonic in that direction
func Bisection(f func(float64) float64, neg, pos, tol float64) (float64, error) {
	negIsInf := math.IsInf(neg, 0)
	posIsInf := math.IsInf(pos, 0)
	var err error
	if negIsInf || posIsInf {
		//fmt.Println("In inf case")
		// See if they're both infinity
		if negIsInf && posIsInf {
			neg, pos, err = bisectionBothInf(f, tol, negIsInf, posIsInf)
		} else {
			// Only one of them is infinity.
			switch {
			// GOT THE LOGIC WRONG FOR ALL THE CASES
			case math.IsInf(neg, 1):
				// Searching toward positive infinity for a positive value
				neg, pos, err = infSearch(f, tol, true, pos, 1, 1, 0)
			case math.IsInf(neg, -1):
				// Searching toward negative infinity for a positive value
				neg, pos, err = infSearch(f, tol, true, pos, -1, 1, 0)
			case math.IsInf(pos, 1):
				//	fmt.Println("In pos inf case")
				// Searching toward positive infinity for a negative value
				neg, pos, err = infSearch(f, tol, false, neg, 1, 1, 0)
			case math.IsInf(pos, -1):
				// Searching toward negatiev infinity for a negative value
				neg, pos, err = infSearch(f, tol, false, neg, 1, 1, 0)
			}
		}
	}
	if err != nil {
		return math.NaN(), err
	}
	if neg == pos {
		// Usually a sign that the inf search found a good value
		return neg, nil
	}
	return boundedBisection(f, neg, pos, tol, 0)
}

func bisectionBothInf(f func(float64) float64, tol float64, negIsInf, posIsInf bool) (neg, pos float64, err error) {
	// Check that they aren't the same infinity
	if math.IsInf(neg, 1) && math.IsInf(pos, 1) {
		return 0, 0, errors.New("Both locations are positive infinity")
	}
	if math.IsInf(neg, -1) && math.IsInf(pos, -1) {
		return 0, 0, errors.New("Both locations are negative infinity")
	}

	// Find initial guesses
	firstLoc := 0.0
	first := f(firstLoc)
	if math.IsNaN(first) {
		return math.NaN(), math.NaN(), errors.New("NaN function value at 0")
	}
	secondLoc := 1.0 + firstLoc
	second := f(secondLoc + firstLoc)
	if math.IsNaN(second) {
		return math.NaN(), math.NaN(), errors.New("NaN function value at 1.0")
	}

	//fmt.Println(firstLoc, first, secondLoc, second)
	// TODO: Add in some NaN checking
	var iter int
	for first == second {
		iter++
		if iter == 100 {
			return math.NaN(), math.NaN(), errors.New("Couldn't find two points with different function values")
		}
		secondLoc = firstLoc + 2*float64(iter)
		second = f(secondLoc)
		if math.IsNaN(first) {
			return math.NaN(), math.NaN(), fmt.Errorf("NaN function value at %v", secondLoc)
		}
		if second != first {
			firstLoc = firstLoc + 2*float64(iter-1)
			break
		}
		secondLoc = firstLoc - 2*float64(iter)
		if math.IsNaN(first) {
			return math.NaN(), math.NaN(), fmt.Errorf("NaN function value at %v", secondLoc)
		}
		if second != first {
			firstLoc = firstLoc - 2*float64(iter-1)
			break
		}
		// TODO: Add something about updating the bound better when we finally find a different point
	}

	// See if either of these are within tol of zero
	if math.Abs(first) < tol {
		return firstLoc, firstLoc, nil
	}
	if math.Abs(second) < tol {
		return secondLoc, secondLoc, nil
	}

	// See if these points have opposite signs. If so, perform the bisection
	firstGTZero := first > 0
	secondGTZero := second > 0
	if !firstGTZero && secondGTZero {
		neg = firstLoc
		pos = secondLoc
		return neg, pos, nil
	}
	if firstGTZero && !secondGTZero {
		neg = firstLoc
		pos = secondLoc
		return neg, pos, nil
	}
	// If not, we need to do a search until we find a point of opposite sign
	firstGTSecond := first > second

	switch {
	case firstGTZero && firstGTSecond:
		// First is greater than the second. Start the search at 1 and head in the
		// positive direction
		return infSearch(f, tol, true, secondLoc, 1, 1, 0)
	case firstGTZero && !firstGTSecond:
		// First is less than second, so search in negative direction
		return infSearch(f, tol, true, firstLoc, -1, 1, 0)
	case !firstGTZero && firstGTSecond:
		// First is closer to zero than second, so search in negative direction
		return infSearch(f, tol, false, first, -1, 1, 0)
	case !firstGTZero && !firstGTSecond:
		// Second is closer to zero than the first, so search in positive direction
		return infSearch(f, tol, false, second, 1, 1, 0)
	default:
		panic("Should never be here")
	}
}

// initLoc sets the initial direction
func infSearch(f func(float64) float64, tol float64, needNeg bool, initLoc, dir, step float64, iter int) (neg, pos float64, err error) {
	if iter > 100 {
		return math.NaN(), math.NaN(), errors.New("Couldn't find two points with different sign")
	}
	/*
		fmt.Println("In inf search")
		fmt.Println("initLoc=", initLoc)
		fmt.Println("dir=", dir)
		fmt.Println("step=", step)
	*/
	newLoc := initLoc + step*dir
	//	fmt.Println("New loc = ", newLoc)
	newVal := f(newLoc)
	//	fmt.Println("New val = ", newVal)

	// See if the new val is close to zero
	if math.Abs(newVal) < tol {
		return newVal, newVal, nil
	}

	if newVal < 0 && needNeg {
		// Found negative point of opposite sign
		return newLoc, initLoc, nil
	}
	if newVal > 0 && !needNeg {
		// Found positive point of oppsoite sign
		return initLoc, newLoc, nil
	}
	// Otherwise, haven't found the point. Continue the search with double the step
	// The new location is closer to the root than the original location (because of
	// monotonicity assumption)
	return infSearch(f, tol, needNeg, newLoc, dir, step*2, iter+1)
}

// Assumes a root between min and max
func boundedBisection(f func(float64) float64, neg, pos, tol float64, iter int) (float64, error) {
	mid := (neg + pos) / 2
	midVal := f(mid)
	if math.IsNaN(midVal) {
		return math.NaN(), fmt.Errorf("NaN function value encountered at %v", midVal)
	}
	if math.Abs(midVal) < tol {
		return mid, nil
	}
	/*
		fmt.Println("neg=", neg)
		fmt.Println("pos=", pos)
		fmt.Println("mid=", mid)
	*/
	if iter > 100 {
		return mid, errors.New("Did not converge after 100 iterations")
	}
	if neg == pos {
		return mid, errors.New("Bounds match but did not find zero")
	}
	if midVal < 0 {
		return boundedBisection(f, mid, pos, tol, iter+1)
	}
	return boundedBisection(f, neg, mid, tol, iter+1)
}
