# Go Design Principles for Great Go Development

## Before we begin
- CLONE ME!!! [github.com/ronna-s/gc-eu-25](https://github.com/ronna-s/gc-eu-25)
- If you haven't already responded, please let me know how much experience you have with Go [here](https://app.sli.do/event/3eTt4zjeY4JEp7moQPpqVc).
- Install Go [here](https://go.dev/dl/)


## Lesson 1: The Go Type-System and Interfaces

### Task: build P&P Game
P&P™ stands for Platforms and Programmers.

P&P is a game in which a band of developers tries to take on the ultimate villain: PRODUCTION.

Today you are going to build parts of the game engine, add characters (developers), and define how they interact with PRODUCTION.

Game starts:

The band of developers initially consists of a test character named after your PM, Sir Tan Lee Knot.
Sir Tan Lee Knot has the following skills:
- Pays wages.

The test band of developers initially consists of a test character named after your PM, Sir Tan Lee Knot. You can find it in `cmd/pnp/main.go`.

- The band starts with 10 gold coins, which is used on every turn by the PM to pay the wages, it can be used to buy objects such as coffee or banana (you'll understand why in a sec), pizza, etc._
- PRODUCTION is "calm".

Run the game by executing `go run cmd/pnp/main.go` to see what you have in action.

#### New requirement 0:
- The main loop should skip dead players.
- The game should end if all players are dead (fired or quit).

#### New requirement 1:
- Every band has a minion - implemented in `pkg/pnp/minion.go`. Minions love PRODUCTION, the ultimate villain, and will do anything to serve it. 
- A minion has only one skill - to cause bugs. A minion cannot learn new skills. A minion can be distracted using a banana.
- Add a minion to every game in the game constructor.
- Since a minion cannot die (it's a minion), we need to figure out a way to plug the minion into the game. The minion must not have an `Alive()  bool` method.
- **Note: do not add the `Alive() bool` method to the minion.**
- When your minion is plugged into the game, you will notice that when it's the minion's turn - the title is messed up. Fix it.

#### New requirements 2:
- If all the players are minions, the game needs to end.
- So, we need to tell if a minion is a minion.

#### You're free now to define your own game:
- You can add a gopher who can add features with a chance of hurting production, or fix bugs to make production happy.
- Feature generate gold coins.
- A player dies at random if production reaches legacy.
- Add a satisfaction system to your game.
- A player dies (quits) if their satisfaction reaches 0.
- Add as many more players as you'd like (for instance, a manager who can buy Pizza for everyone - it's supported in the engine already).
- Add leveling up of the players where they have even more skills.
- Whatever strikes your fancy.
- Decide when the game is won - for instance, when you reach 100 gold coins.

## Lesson 2: Generics

Generics are a pretty late edition to the Go language, but few people know that they were already with us all this time,
for instance: `append` takes a slice of a type T and append items of type T to it, resulting in another slice of type T.
The problem was that the language had generics and could use generic types, but we couldn't define our own.

### When to use them?
It was a long-time consensus that "real gophers" don't need generics, so much so that around the time the generics draft of 2020 was released, many gophers still said that they are not likely to use them.

Let's first understand the point that they were trying to make.

Consider [this code](https://gist.github.com/Xaymar/7c82ed127c8f1def53075f414a7df153), made using C++.
We see here generic code (templates) that allows an event to add functions (listeners) to its subscribers.
Let's ignore for a second that this code adds functions, not objects and let's assume it did take in objects with the function `Handle(e Event)`.
We don't need generics in Go to make this work because interfaces are implicit. As we saw already in C++ an object has to be aware of its implementations, this is why to allow plugging-in of functionality we have to use generics in C++ (and in Java).

In Go this code would look something like [this](https://go.dev/play/p/Tqm_Hb0vcZb):

```go
package main

import "fmt"

type Listener interface {
	Handle(Event)
}

type Event struct {
	Lis []Listener
}

func (e *Event) Add(l Listener) {
	e.Lis = append(e.Lis, l)
}

func main() {
	var l Listener
	var e Event
	e.Add(l)
	fmt.Println(e)
}
```

**We didn't need generics at all!**

However, there are cases in Go where we have to use generics and until recently we used code generation for.
Those cases are when the behavior is derived from the type or leaks into the type's behavior:

For example:
The linked list
```go
// https://go.dev/play/p/ZpAqvVFAIDZ
package main

import "fmt"

type Node[T any] struct { // any is builtin for interface{}
  Value T
  Next  *Node[T]
}

func main() {
  n1 := Node[int]{1, nil}
  n2 := Node[int]{3, &n1}
  fmt.Println(n2.Value, n2.Next.Value)
}
```
Example 2 - [Addition](https://go.dev/play/p/dmeQEVxpyAq)
```go
package main

import "fmt"

type A int

// Add takes any type with underlying type int 
func Add[T ~int](i, j T) T { 
  return i + j
}

func main() {
  var i, j A
  fmt.Println(Add(i, j))
}
```
Of course, you might not be likely to use linked lists in your day to day, but you are likely to use:
1. Repositories, database models, data structures that are type specific, etc.
2. Event handlers and processors that act differently based on the type.
3. The [concurrent map in the sync package](https://pkg.go.dev/sync#Map) which uses the empty interface.
4. [The heap container](https://pkg.go.dev/container/heap#example-package-IntHeap)

The common thread to these examples is that, whereas before generics we had to trade generalizing certain behavior for type safety (or generate code to do so), now we can have both.

So, how does it work, exactly?

To use generic types in Go, we have to tell the compiler something about the type that we expect, using a constraint.
Constraints are defined using interfaces.
1. If our code supports any type - we can use the `any` keyword (stand-in for the empty interface `interface{}`).
2. If our codes expects a type with a subset of behaviors, we use an interface with the methods that we need (very similarly to using regular interfaces).
3. If our code expects a type with an underlying type of a certain type, we use the `~` operator. e.g. `interface{~int}`.
4. If our code expects an exact type - we use the type name. e.g. `inteface{string}`.
5. We can also union types using the `|` operator. e.g. `interface{int|string}`. This constraint will allow using `+` on strings and ints alike.

Sometimes we need to get more creative with generics and express dependencies between types. For instance, when the pointer to a type implements a constraint (the interface), but we also need the type itself.
For instance, we can expect a type T and another type `interface{~[]T}` which is any type with an underlying type of a slice of T.
We can express a type which is a function that returns T like so: `interface{~func() T}`.'

Consider the following example where we need to populate a variable of type T, but the interface is implemented by its pointer.
PT is defined to be a pointer to T (a dependency on the previous generic type) and we also provide the interface methods that it implements (by embedding proto.Message):
```go
import "google.golang.org/protobuf/proto"

func DoSomething[T any, PT interface {
	proto.Message
	*T
}]() {

	var t T
	var protoMessage PT = &t
	// do something to populate t
}
```

### Task
Introduce a leaderboard to the game in `pkg/pnp/leaderboard.go` that will keep track of the players' scores.
We will define our own generic heap to keep the top leaderboard sorted by score.

**Justice to the heap!**

## Lesson 3: Functional Programming

### What is functional programming?
Functional programming is the practice of building software by composing pure functions, avoiding shared state, mutable data, and side-effects. It is declarative rather than imperative, and application state flows through pure functions. E.g.:

```go
fn1(fn2(fn3(someArgs)))
```
Instead of:

```go
res := Proc1(someArgs)
Proc2(otherArgs, res)
sideEffects := ReadSideEffects()
Proc3(sideEffects)
```
Functional programming is a programming paradigm that treats computation as the evaluation of mathematical functions and avoids changing state and mutable data. In Go, we can use functional programming techniques such as first-class functions, higher-order functions, and closures.
The introduction of generics in Go allows us to write more generic and reusable data types representing functions (pure or otherwise).
This makes functional programming more accessible in Go, as we can now define functions that operate on different types without losing type safety.

A limitation of generics with functional programming in Go has to do with the fact that [Go doesn't support methods taking type parameters](https://github.com/golang/go/issues/49085) ([explanation](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods)).
This means that when we opt to work with functional programming, for best design, we should avoid methods, but rather think about how to apply functions to our types.

### Why use functional programming?
Functional programming is a powerful tool for building resilient and maintainable software.

It's hard to reason about a program or function when its result depends on
state that can change. As a very simple example, consider this impure function:
```
var g int

func addg(x int) int { return x + g}
```
What does `addg(1)` return? The answer depends on the value of `g`, which itself
can depend on anything in the program.
By contrast, it's easy to figure out what the pure version returns:
```
func add(x, y int) int { return x + y}
```
Everything you need is right there in the function itself.

Functional programming is also kind of cool.

### Is Go a functional programming language?

Not really. Go is a multi-paradigm language that supports procedural, object-oriented, and concurrent programming. However, Go has some functional programming features, such as first-class functions, higher-order functions, and closures.

- Go doesn't provide lazy evaluation, which is a feature of some functional programming languages that allows expressions to be evaluated only when needed. (So we have to do lazy evaluation ourselves).
- Go doesn't provide immutable data structures, which are common in functional programming languages.
- Go doesn't have tail call optimization, which is a feature of functional programming languages that allows recursive functions to be optimized to avoid stack overflow (and also generally good for performance).

### What is TCO (Tail Call Optimization) and why is it important?
It's the practice of overriding the current stack frame with the next one (we just jump to the last function in the current function), instead of adding a new one.
It's useful to allow iterations using recursions without causing a stack overflow.
Go doesn't support tail call optimization for [various reasons](https://github.com/golang/go/issues/22624) ([see also](https://groups.google.com/g/golang-nuts/c/nOS2FEiIAaM/m/miAg83qEn-AJ)), despite some suggestions otherwise online.
For recursive calls, since the Go stack is managed on the heap, we can still run recursions with a larger stack (but we should prefer loops instead).

### See also:
- A very cohesive [Functional Programming package for Go by IBM](https://github.com/IBM/FP-GO) and the [video](https://www.youtube.com/watch?v=Jif3jL6DRdw) about it.
- An Introduction to Functional Programming in Go - Eleanor McHugh [video](https://www.youtube.com/watch?v=OKlhUv8R1ag).
- Monadic operartions overview (in Haskell - but explains basic concepts) [post](https://www.adit.io/posts/2013-04-17-functors,_applicatives,_and_monads_in_pictures.html)


### Task
Our PM asks that we remove the AsciiArt() method from the Player interface, because it is not used in many engines that people plug into the game.
In order to do that, we will check at runtime if player implements the AsciiArt() method, and if it does, we will call it, otherwise we will check if it has a String() string method, otherwise check if it has a `Name() string`,
and if not we will default to `Player #<number>`.

To do that we will use functional programming techniques to chain the calls together, for each that doesn't work, we will call the next one until we get a working string.

Note: We don't want to call any consecutive functions, this computation should be lazy.

**Bonus:** Make the code that reads the scores from the leaderboard in `pkg/pnp/leaderboard.go` functional, using `Left()` and `Right()` functions from github.com/IBM/fp-go package (included in the modules).

## Lesson 4: Concurrency and Testing 

1. We are going to fix some code to make it testable in pkg/concurrency.
2. We have a rogue goroutine in our game engine that nobody (i.e. me) bothered to keep track of and test, we are going to introduce some tests to ensure that it terminates properly.
2. We are going to introduce graceful shutdown to our game, so that we can stop the game and all the goroutines gracefully while ensuring that if the game ends unexpectedly - the leaderboard remains up to date. (ahhmm, there are no transactions in this code, so really no guarantees but you get the idea.)
