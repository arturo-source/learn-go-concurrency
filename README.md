# Go concurrency

To start learning about go concurrency, I've chosen my trusted site <https://gobyexample.com/>. There you can see the simplest examples for the standard programming cases in Go. I recommend starting from **Goroutines** and reading all the examples until **Stateful Goroutines**.

Obviously you should also read the Go tour <https://go.dev/tour/concurrency/1>, which have some examples that are useful to start.

The last thing that you should take a quick look is this video <https://youtu.be/LvgVSSpwND8>. Jake Wright talks about WaitGroup, Channels, Deadlock and Channel closing, buffered Channels, Select statement, and the Worker pool pattern, this is a good introduction for Golang concurrency. **Old but gold**.

At this time, I recommend you to stop reading, and practice with some small problems (e.g. [Tree](01.%20Tree) and [Crawler](02.%20Crawler)).

After solving some small problems, you can continue watching the [Rob Pike conference](https://www.youtube.com/watch?v=f6kdp27TYZs) (or you can read [the slides](https://go.dev/talks/2012/concurrency.slide#1), but I recommend the video). Here he explains what is concurrency, some Go concurrency "patterns", and why Go concurrency was made as a part of the language. **Old but gold**.

If you want to get harder, read this post <https://go.dev/blog/pipelines> from the go official blog. This post takes a real example where you can apply concurrency, and breaks it down carefully, explaining why each decision is made. I have modified the problem to do a similar exercise using Go's concurrency: [MD5 files](03.%20MD5%20files).

Race conditions are important to consider in any language that allows multi-threading. [In this post](https://go.dev/blog/race-detector) you'll see a real story of the go team detecting one, with the go `-race` tool. But if you want a more detailed explanation of how to use that tool, read that one: <https://go.dev/doc/articles/race_detector>.

I'll attach some questions I had during my learning process here:

## Channels

### Who should close the channel?

Always **the sender**. Trying to write a closed channel will panic and terminate the program. Trying to read a closed channel will return a nil value (it depends on the channel type: bool -> false, int -> 0, etc.). It's mentioned on [this page](https://go.dev/tour/concurrency/4) of the Go tour.

### Should I close the channel?

**It's not necessary to close the channel**, the garbage collector will do (also mentioned in [the Go tour](https://go.dev/tour/concurrency/4)). Sometimes you will want to close the channel, for example when the reader uses a for-range, that will end the for loop. But the important thing is that **not closing the channel will not produce a memory leak**, you don't have to worry about that.

### How to know if a channel is closed

If you are a sender, you shoud not send data throught a closed channel, but there is not a way to know if channel is closed. If you are a receiver, it is really easy to check `v, isOpen := <-ch`, if the channel is closed `isOpen` will be false. But there is some syntactic sugar that will allow you to receive data until the channel is closed:

```go
func main() {
 c := make(chan int, 10)
 go fibonacci(c)
 for i := range c {
  fmt.Println(i)
 }
}
```

If you use the for-range golang style, this will end when the channel is closed. You will not have to worry about checking if channel is closed.

### What causes panic with channels?

- Use a non-initialized channel (var ch chan int).
- Send data through a closed channel.
- Close an already closed channel.
- Deadlocks: `fatal error: all goroutines are asleep - deadlock!`
  - Receive from a non-closed channel. If no one sends data.
  - Send to a filled channel. If no one receives the data.

Initializing channels is easy, but `var ch chan int` is not initializing. Instead use `var ch chan int = make(chan int)` or `ch := make(chan int)`.

You have to be quiet with closed channels. As said before, you can only send data through non-closed channels. If you send data through a closed channel, it will panic. If you close an already closed channel, it will panic too. (I will add an example of how to close multiple sender channel later...)

Reading a channel can also cause panic, if you don't send data through this channel. The previous example `for i := range c {...}` will panic if the channel `c` is never closed.

But writting can also cause panic, if you don't receive data through this channel. The previous example `fibonacci(c)` will panic if `c` is never read. If the `for` does not exist, when the channel receives the 11th fibonacci result, its buffer will be filled, and it will panic.

## Select

### What is a select?

**Select is a switch-like solution to wait for one of multiple channels**. Sometimes you will want to block the goroutine until one of multiple channels can run. This is the perfect case to use select. It's important to know that select accepts _reading_ case and _writing_ case.

```go
select {
case a <- 10: // a is writing the channel
  // Do something
case x := <- b: // b is reading the channel
  // Do something
}
```

But sometimes you will want to wait for multiple channels in an infinite for loop, to response all of them.

### Be quiet with Select

If you run a `select` in a for loop, with a `default` case, the program may not act as you expected. This is an example **modified** from <https://go.dev/tour/concurrency/6>:

```go
package main

import (
 "fmt"
 "time"
)

func main() {
 tick := time.Tick(100 * time.Millisecond)
 boom := time.After(500 * time.Millisecond)
 for i:= 0; i < 4; i++ {
  select {
  case <-tick:
   fmt.Println("tick.")
  case <-boom:
   fmt.Println("BOOM!")
   return
  default:
   fmt.Println("    .")
  }
 }
}
```

This will cause:

```txt
    .
    .
    .
    .
```

But you will have a problem also with closed channels, that you may not expect:

```go
for {
  select {
  case msg1 := <-a:
    // Do something
  case msg2 := <-b:
    // Do something
  }
}
```

This seems a good way to block a goroutine until `a` or `b` receives data, infinitely. But if you know that `a` can be closed, you may think that the select will only choose `b` case, but **no**. `msg1` will be the `a` nil value (false for bool, 0 for int, etc.) and will be executed infinitely, and it will cause a 100% of use of your CPU. So you should think how to control this case.

## Timers

### How to set timeouts?

This is clearly explained in <https://gobyexample.com/timeouts>. You can use [time package](https://pkg.go.dev/time) from the Go standard packages, which provide some useful functions like `After` or `Tick`. You can combine this functions with `select` to get what you want.

Also, some librarires, like the http standard library, allows you to set timeouts directly when creating the http client, or server.

But the standard way to set timeouts for complex programs isn't the the `time` package, you should use the [context package](https://pkg.go.dev/context).

## Synchronization

### Synchronize multiple goroutines

One of the most useful utilities of golang channels is synchronizing read and write operations. But there are other options to do that.

If you want to throw multiple goroutines in the same function, and wait until the end, the more comprehensive way is [using WaitGroup](https://gobyexample.com/waitgroups). You could also [use a channel](https://gobyexample.com/channel-synchronization) but it is easier and more readable with WaitGroup in this case.

If you want to implement a semaphore you could use a channel, but it can become hard to read again, so the reasonable way is [using Mutex](https://gobyexample.com/mutexes).

### Important about Mutex and WaitGroup

Note that if you want to `wg.Done()` or `mu.Lock()` (for [WaitGroup](https://pkg.go.dev/sync#WaitGroup) and [Mutex](https://pkg.go.dev/sync#Mutex) respectively) you should always use the same, not a copy, using pointers or global variables. As said in the [sync package](https://pkg.go.dev/sync):

```txt
Values containing the types defined in this package should not be copied.
```

## Race conditions

### What are race conditions?

A race condition occurs when two (or more) goroutines try to use the same memory. Two different routines can read the same place of memory safely. But when those two try to read and write, or write and write in the same space of memory, it produces a **race condition**.

### How to avoid race conditions?

Avoiding race conditions is simple and complicated at the same time. You just have to follow the previous commented (do not write and write, or write and read the same address of memory at the same time). It is easy to understand, but in the practice you can create a race condition, even if you know that, because programming becomes harder when you write bigger programs.

Fortunately, golang provides powerful tools to operate with concurrency, and avoid race conditions.

1. The already mentioned **channels** are syntactic sugar for block operations (channels are even more, but also, block operations).
2. The [sync package](https://pkg.go.dev/sync) provides safe **map** implementation. A **wait group** to synchronize routines easily. Also some lower-level utilities like **mutex**.
3. The [sync/atomic package](https://pkg.go.dev/sync/atomic) provides safe operations with basic types: bool, int, etc.

### How to detect race conditions?

A race condition can **NOT** be detected at compile-time. But there are tools to detect them at run-time.

Go implemented a race detector in go1.1, so you can simply add `-race` option when you build, run, test, or install a golang program. But be quiet, because this option makes the program much more slower, and it'll cosume more resources, so use that option only in development enviroments.

```sh
go run -race main.go
```

### Examples of race conditions

Default maps are **not** concurrent safe, even if you initialice the map before accessing it. Next example produces race conditions.

```go
func main() {
 values := map[string]int{"first": 0, "second": 0}
 wg := sync.WaitGroup{}
 wg.Add(len(values))

 for k := range values {
  go func(k string) {
   defer wg.Done()

   for j := 0; j < 10; j++ {
    values[k]++
   }
  }(k)
 }

 wg.Wait()
 fmt.Println(values)
}
```

Slices do not share memory, so they are concurrent safe. But **be quiet**! This only happens when the slice is already filled. If you try to append data to the slice, it can produce a race condition.

```go
func main() {
 values := []int{0, 0}
 wg := sync.WaitGroup{}
 wg.Add(len(values))

 for i := range values {
  go func(i int) {
   defer wg.Done()

   for j := 0; j < 10; j++ {
    values[i]++
   }
  }(i)
 }

 wg.Wait()
 fmt.Println(values)
}
```

Before go1.22 there was kind of a "bug" that shared memory of the loop index. Next example produced race conditions, but now, it doesn't, so you are allowed to do it.

```go
func main() {
 var wg sync.WaitGroup
 wg.Add(10)

 for i := 0; i < 10; i++ {
  go func() {
   defer wg.Done()
   fmt.Println(i)
  }()
 }

 wg.Wait()
}
```

## Deadlocks

TODO

## Go concurrency patterns

TODO

## More questions

### Do goroutines and channels produce memory leaks?

Short answer: **channels no, goroutines yes**. As said before, channels are managed by the garbage collector, you should close channels if you want to communicate that you will not continue sending data. It's not the same for the goroutines, if you block a goroutine, don't stop an infinite for loop, or whatever that prevents the goroutine finishing, the garbage collector will not be able to close that goroutine. Let's see a simple example:

```go
package main

import (
 "errors"
 "fmt"
 "runtime"
 "time"
)

func NumGoroutines(when string) {
 fmt.Println(when, runtime.NumGoroutine())
}

func main() {
 defer NumGoroutines("After end:")
 NumGoroutines("Beginning:")
 errc := make(chan error) // channel len 0

 go func() {
  // Do some stuff that produces error
  err := errors.New("some error")
  errc <- err
 }() // If nobody reads errc (<-errc), goroutine keeps blocked and never ends...
 NumGoroutines("After starting a goroutine:")

 // The code continues...
 time.Sleep(time.Second)
}
```

With `runtime.NumGoroutine()` we can count active goroutines. With `defer` we execute that after the `main` end. Then, we start a goroutine that produces an error, we could expect that this error will be read, but if nobody reads errc, this goroutine will never end. That code produces the next output:

```txt
Beginning: 1
After starting a goroutine: 2
After end: 2
```

This problem would be solved using a buffer to errc `errc := make(chan error, 1)`. This don't block the goroutine, then it's finished, and you don't have to worry about `errc` because garbage collector will free `errc` memory when it detects it'll never be used.

## To do

- I'd like to take a look at <https://github.com/golang/go/wiki/LearnConcurrency>.
- What is the best way to control errors? In channels you only send one type of data.
