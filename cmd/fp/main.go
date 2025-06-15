package main

import (
	"fmt"
	"log"

	"golang.org/x/exp/constraints"
)

// Applicative is a functional programming concept that generalizes functors.
// It allows you to apply functions that are themselves wrapped in a context (such as a container or computation)
// to values that are also wrapped in that context. In Go, an Applicative[T] is often represented as a function
// returning a value of type T (e.g., func() T), enabling you to compose computations in a lazy and type-safe way.
// Applicatives make it possible to combine multiple independent computations and apply functions to their results
// without having to unwrap the values from their context.
type Applicative[T any] func() T

func Ap[A any, B any](fa Applicative[func(A) B], a Applicative[A]) Applicative[B] {
	return func() B {
		return fa()(a())
	}
}

type Task func()

func NewTask() Task {
	return func() {}
}

func (t Task) Then(fn func()) Task {
	return func() {
		t()
		fn()
	}
}

func (t Task) Execute() {
	t()
}

func Pipe[A, B any, F ~func() A, G ~func(A) B](f F, g G) func() B {
	return func() B {
		return g(f())
	}
}

type Ordered[T constraints.Ordered] func() T

func NewAccumulator[T constraints.Ordered](x T) Ordered[T] {
	return func() T {
		log.Println("computing")
		return x
	}
}

func (a Ordered[T]) ToApplicative() Applicative[T] {
	return func() T {
		return a()
	}
}

func (a Ordered[T]) Accumulate4(accumulator Ordered[T]) Ordered[T] {
	return func() T {
		return a() + accumulator()
	}
}

func (a Ordered[T]) Accumulate(y T) Ordered[T] {
	return func() T {
		return a() + y
	}
}

func (a Ordered[T]) Accumulate2(y Ordered[T]) Ordered[T] {
	return func() T {
		return a() + y()
	}
}

func (a Ordered[T]) Accumulate3(y Accumulator[T, Ordered[T]]) Ordered[T] {
	return func() T {
		return a() + y.Compute()
	}
}

func (a Ordered[T]) Compute() T {
	return a()
}

type Func[T any] func() T

func Fmap[A any, B any, F ~func() A](fa F, f func(A) B) func() B {
	return func() B {
		return f(fa())
	}
}

type Accumulator[T any, Computer interface{ Compute() T }] interface {
	Accumulate(y T) Computer
	Compute() T
}

func main() {

	t := NewTask()
	t = t.Then(func() {
		fmt.Print("Hello ")
	}).Then(func() {
		fmt.Print("World")
	}).Then(func() {
		fmt.Print("!\n")
	})

	t.Execute()

	f := func(a int) func(int) int {
		return func(b int) int {
			return a * b
		}
	}

	a := NewAccumulator[int](0)
	a = a.Accumulate(3)
	a = a.Accumulate2(NewAccumulator(2).Accumulate(4))
	a = a.Accumulate3(NewAccumulator(1).Accumulate(4))
	b := NewAccumulator[int](-1)
	b = b.Accumulate(4)
	a = a.Accumulate4(b)

	log.Println("starting to compute stuff - everything so far was lazy")
	fmt.Println(a())

	var v Accumulator[int, Ordered[int]]
	v = a

	fmt.Println(v.Accumulate(5).Compute())

	acc := NewAccumulator[int](2)
	acc = acc.Accumulate(3)
	timesFive := Fmap(acc, f)

	//This will not cause anything - because this is lazy!
	Fmap(NewAccumulator(3), func(i int) int {
		a := 0
		return i / a
	})

	five := NewAccumulator(2)
	value := Ap(timesFive, five.ToApplicative())
	fmt.Println(value())

	p := Pipe(
		func() string {
			return "hi"
		},
		func(s string) struct{} {
			fmt.Println(s, "bye!")
			return struct{}{}
		},
	)
	fmt.Println(p())
}
