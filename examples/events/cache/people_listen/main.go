/*
 * Copyright (c) 2022, 2023 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * https://oss.oracle.com/licenses/upl.
 */

/*
Package main listens for all events on a NamedMap or NamedCache.
*/
package main

import (
	"context"
	"fmt"
	"github.com/oracle/coherence-go-client/coherence"
	"log"
	"time"
)

type Person struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Age        int    `json:"age"`
	InsertTime int64  `json:"insertTime"`
}

func (p Person) String() string {
	return fmt.Sprintf("Person{id=%d, name=%s, age=%d, insertTime=%v}", p.ID, p.Name, p.Age, p.InsertTime)
}

func main() {
	ctx := context.Background()

	// create a new Session
	session, err := coherence.NewSession(ctx, coherence.WithPlainText())
	if err != nil {
		panic(err)
	}

	defer session.Close()

	// create a new NamedMap of Person with key int
	namedMap, err := coherence.NewNamedMap[int, Person](session, "people")
	if err != nil {
		panic(err)
	}

	fmt.Println("Adding listener for all events")
	// Create a listener and add to the cache
	listener := NewAllEventsListener[int, Person]()
	if err = namedMap.AddListener(ctx, listener.listener); err != nil {
		panic(err)
	}

	defer func(ctx context.Context, namedMap coherence.NamedMap[int, Person], listener *AllEventsListener[int, Person]) {
		if err := namedMap.RemoveListener(ctx, listener.listener); err != nil {
			panic(fmt.Sprintf("cannot remove listener %v, %v", listener.listener, err))
		}
	}(ctx, namedMap, listener)

	time.Sleep(time.Duration(10000000) * time.Second)
}

type AllEventsListener[K comparable, V any] struct {
	listener coherence.MapListener[K, V]
}

func NewAllEventsListener[K comparable, V any]() *AllEventsListener[K, V] {
	exampleListener := AllEventsListener[K, V]{
		listener: coherence.NewMapListener[K, V](),
	}

	exampleListener.listener.OnAny(func(e coherence.MapEvent[K, V]) {
		var (
			newValue *V
			oldValue *V
			key      *K
			err      error
		)
		key, err = e.Key()
		if err != nil {
			panic("unable to deserialize key")
		}

		if e.Type() == coherence.EntryInserted || e.Type() == coherence.EntryUpdated {
			newValue, err = e.NewValue()
			if err != nil {
				panic("unable to deserialize new value")
			}
		}
		if e.Type() == coherence.EntryDeleted || e.Type() == coherence.EntryUpdated {
			oldValue, err = e.OldValue()
			if err != nil {
				panic("unable to deserialize old value")
			}
		}

		log.Printf("***EVENT=%v: key=%v, oldValue=%v, newValue=%v\n", e.Type(), *key, oldValue, newValue)
	})

	return &exampleListener
}