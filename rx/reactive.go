package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/iterable"
	"github.com/reactivex/rxgo/observable"
	"github.com/reactivex/rxgo/observer"
)

func main() {
	//veryFirstError()
	//veryFirstReactive()
	//veryFirstFlatMap()
	sum()
}

func veryFirstError() {
	watcher := observer.Observer{
		// Register a handler function for every next available item.
		NextHandler: func(item interface{}) {
			fmt.Printf("Processing: %v\n", item)
		},

		// Register a handler for any emitted error.
		ErrHandler: func(err error) {
			fmt.Printf("Encountered error: %v\n", err)
		},

		// Register a handler when a stream is completed.
		DoneHandler: func() {
			fmt.Println("Done!")
		},
	}

	it, _ := iterable.New([]interface{}{1, 2, 3, 4, errors.New("bang"), 5})
	source := observable.From(it)
	sub := source.Subscribe(watcher)

	// wait for the channel to emit a Subscription
	<-sub
}

func veryFirstReactive() {
	score := 9

	onNext := handlers.NextFunc(func(item interface{}) {
		fmt.Println("onNext")
	})

	onDone := handlers.DoneFunc(func() {
		fmt.Println("onDone")
	})

	watcher := observer.New(onNext, onDone)

	sub := observable.Just(1).Subscribe(watcher)

	<-sub

	fmt.Println(score)
}

func veryFirstFlatMap() {
	primeSequence := observable.Just([]int{2, 3, 5, 7, 11, 13})

	<-primeSequence.
		FlatMap(func(primes interface{}) observable.Observable {
			return observable.Create(func(emitter *observer.Observer, disposed bool) {
				for _, prime := range primes.([]int) {
					emitter.OnNext(prime)
				}
				emitter.OnDone()
			})
		}, 1).
		Filter(func(item interface{}) bool {
			return item.(int) > 5
		}).
		First().
		Subscribe(handlers.NextFunc(func(prime interface{}) {
			fmt.Println("Prime -> ", prime)
		}))
}

func sum() {
	var num1, num2 int
	reader := bufio.NewReader(os.Stdin)

	processText := func(text, prefix string, numPtr *int) {
		text = strings.Trim(strings.TrimPrefix(text, prefix), " \n")
		*numPtr, _ = strconv.Atoi(text)
	}

	onNext := handlers.NextFunc(func(item interface{}) {
		if text, ok := item.(string); ok {
			switch {
			case strings.HasPrefix(text, "a:"):
				processText(text, "a:", &num1)
			case strings.HasPrefix(text, "b:"):
				processText(text, "b:", &num2)
			default:
				fmt.Println("Input does not start with prefix \"a:\" or \"b:\"!")
				return
			}
		}

		fmt.Printf("Running sum: %d\n", num1+num2)
	})

	for {
		fmt.Print("type>")

		sub := observable.Start(func() interface{} {
			text, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			return text
		}).Subscribe(onNext)

		<-sub
	}
}
