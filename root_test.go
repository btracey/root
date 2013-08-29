package root

import (
	"math"
	"testing"

	"fmt"
)

type RootFinder func(func(float64) float64, float64, float64, float64) (float64, error)

type TestFunction func(float64) float64

type RootTest struct {
	Fun         TestFunction
	Neg         float64
	Pos         float64
	Tol         float64
	ShouldError bool
	Name        string
}

func GetTestFunctions() []RootTest {
	return []RootTest{
		RootTest{
			Name: "EasyLinear",
			Fun:  func(x float64) float64 { return x - 7 },
			Neg:  -3,
			Pos:  10,
			Tol:  1e-14,
		},
		RootTest{
			Name: "NegInfLinear",
			Fun:  func(x float64) float64 { return x - 7 },
			Neg:  math.Inf(-1),
			Pos:  9.5,
			Tol:  1e-14,
		},
		RootTest{
			Name: "PosInfLinear",
			Fun:  func(x float64) float64 { return x - 7 },
			Neg:  0.1,
			Pos:  math.Inf(1),
			Tol:  1e-14,
		},
		RootTest{
			Name: "BothInfLinear",
			Fun:  func(x float64) float64 { return x - 7 },
			Neg:  math.Inf(-1),
			Pos:  math.Inf(1),
			Tol:  1e-14,
		},
		RootTest{
			Name:        "SwitchedInputs",
			Fun:         func(x float64) float64 { return x - 7 },
			Neg:         10,
			Pos:         -3,
			Tol:         1e-14,
			ShouldError: true,
		},
	}
}

func testRootFinder(t *testing.T, root RootFinder) {
	functions := GetTestFunctions()
	for _, fun := range functions {
		fmt.Println("Starting case " + fun.Name)
		zero, err := root(fun.Fun, fun.Neg, fun.Pos, fun.Tol)
		if err != nil {
			if fun.ShouldError {
				continue
			}
			t.Errorf("Error finding root for case " + fun.Name + ": " + err.Error())
			continue
		}
		if err == nil && fun.ShouldError {
			t.Errorf("No error for case " + fun.Name)
			continue
		}
		if math.Abs(fun.Fun(zero)) > fun.Tol {
			t.Errorf("Zero location not close to zero. Zero found is: %v", zero)
		}
	}
}

func TestBisection(t *testing.T) {
	testRootFinder(t, Bisection)
}
