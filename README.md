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
It was a long time consensus that "real gophers" don't need generics, so much so that around the time the generics draft of 2020 was released, many gophers still said that they are not likely to use them.

Let's understand first the point that they were trying to make.

Consider [this code](https://gist.github.com/Xaymar/7c82ed127c8f1def53075f414a7df153), made using C++.
We see here generic code (templates) that allows an event to add functions (listeners) to its subscribers.
Let's ignore for a second that this code adds functions, not objects and let's assume it did take in objects with the function `Handle(e Event)`.
We don't need generics in Go to make this work because interfaces are implicit. As we saw already in C++ an object has to be aware of it's implementations, this is why to allow plugging-in of functionality we have to use generics in C++ (and in Java).

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
Those cases are when the behavior is derived from the type or leaks to the type's behavior:

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
func Add[T ~int](i T, j T) T { 
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
4. [The heap](https://pkg.go.dev/container/heap#example-package-IntHeap)

The common thread to these examples is that before generics we had to trade generalizing certain behavior for type safety (or generate code to do so), now we can have both.

So, how does it work, exactly?

To use generic types in Go, we have to tell the compiler something about the type that we expect, using a constraint.
Constraints are defined using interfaces.
1. If our code supports any type - we can use the `any` keyword (stand-in for the empty interface `interface{}`).
2. If our codes expects a type with a subset of behaviors, we use an interface with the methods that we need (very similarly to using regular interfaces).
3. If our code expects a type with an underlying type of a certain type, we use the `~` operator. e.g. `interface{~int}`.
4. If our code expects an exact type - we use the type name. e.g. `inteface{string}`.
5. We can also union types using the `|` operator. e.g. `interface{int|string}`. This constraint will allow using + on strings and ints alike.

Sometimes we need to get more creative with generics and express dependencies between types. For instance, when the pointer to a type implements a constraint (the interface), but we also need the type itself.
For instance, we can expect a type T and another type `interface{~[]T}` which is any type with an underlying type of a slice of T.
We can expect a type which is a function that returns T like so: `interface{~func() T}`.'

Consider the following example where we need to populate a variable of type T, but the interface is implemented by its pointer.
PT is defined to be a pointer to T (a dependency on the previous generic type) and we provide also the interface methods that it implements (by embedding proto.Message).:
```go
import (
	"google.golang.org/protobuf/proto"
)

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
Define the generic heap in heap/heap.go so that cmd/top.go compiles and the tests in heap_test.go are successfully run.


## Lesson 3: Functional Programming

## Lesson 4: Concurrency and Testing 