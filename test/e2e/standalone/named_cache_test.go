/*
 * Copyright (c) 2022, 2023 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * https://oss.oracle.com/licenses/upl.
 */

package standalone

import (
	"github.com/onsi/gomega"
	"github.com/oracle/coherence-go-client/coherence"
	"github.com/oracle/coherence-go-client/coherence/processors"
	. "github.com/oracle/coherence-go-client/test/utils"
	"testing"
	"time"
)

func TestPutWithExpiry(t *testing.T) {
	var (
		g        = gomega.NewWithT(t)
		err      error
		person1  = Person{ID: 1, Name: "Tim"}
		oldValue *Person
	)

	session, err := GetSession()
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	defer session.Close()
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	namedCache := GetNamedCache[int, Person](g, session, "put-with-expiry")

	defer session.Close()

	oldValue, err = namedCache.PutWithExpiry(ctx, person1.ID, person1, time.Duration(5)*time.Second)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(oldValue).To(gomega.BeNil())
	AssertSize[int, Person](g, namedCache, 1)

	// sleep for 6 seconds to allow entry to expire
	time.Sleep(6 * time.Second)

	AssertSize[int, Person](g, namedCache, 0)
}

// TestTouchProcessor tests a touch processor that will update the time of en entry
func TestTouchProcessor(t *testing.T) {
	var (
		g           = gomega.NewWithT(t)
		err         error
		person1     = Person{ID: 1, Name: "Tim"}
		containsKey bool
		oldValue    *Person
	)

	session, err := GetSession()
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	namedCache := GetNamedCache[int, Person](g, session, "touch")

	defer session.Close()

	// "touch" cache has default TTL of 10 seconds
	_, err = namedCache.Put(ctx, person1.ID, person1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(oldValue).To(gomega.BeNil())
	AssertSize[int, Person](g, namedCache, 1)

	// sleep for 6 seconds and the entry should still be there
	time.Sleep(6 * time.Second)

	containsKey, err = namedCache.ContainsKey(ctx, 1)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(containsKey).To(gomega.Equal(true))

	// run the Touch processor which will reset the TTL
	_, err = coherence.Invoke[int, Person, any](ctx, namedCache, 1, processors.Touch())
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// sleep another 6 seconds, which will be approx 12 seconds since original put
	// entry should still exist due to Touch processor
	containsKey, err = namedCache.ContainsKey(ctx, 1)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(containsKey).To(gomega.Equal(true))

	// sleep for 10 seconds and the entry should now be evicted
	time.Sleep(10 * time.Second)

	containsKey, err = namedCache.ContainsKey(ctx, 1)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(containsKey).To(gomega.Equal(false))
}

func TestTestMultipleCallsToNamedCache(t *testing.T) {
	var (
		g            = gomega.NewWithT(t)
		err          error
		person1      = Person{ID: 1, Name: "Tim"}
		personValue1 *Person
		personValue2 *Person
		session      *coherence.Session
	)

	session, err = GetSession()
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	defer session.Close()

	namedCache1, err := coherence.NewNamedCache[int, Person](session, "cache-1")
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	err = namedCache1.Clear(ctx)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// retrieve the named map again, should return the same one
	namedCache2, err := coherence.NewNamedCache[int, Person](session, "cache-1")
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	err = namedCache2.Clear(ctx)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	g.Expect(namedCache2).To(gomega.Equal(namedCache1))

	_, err = namedCache1.Put(ctx, person1.ID, person1)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	personValue1, err = namedCache1.Get(ctx, person1.ID)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	personValue2, err = namedCache2.Get(ctx, person1.ID)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	g.Expect(*personValue1).To(gomega.Equal(*personValue2))

	namedCache3, err := coherence.NewNamedCache[int, Person](session, "cache-2")
	g.Expect(err).NotTo(gomega.HaveOccurred())

	size, err := namedCache3.Size(ctx)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(size).To(gomega.Equal(0))

	// try and retrieve a NamedCache that is for the same cache but different type, this should cause error
	_, err = coherence.NewNamedCache[int, string](session, "cache-2")
	g.Expect(err).To(gomega.HaveOccurred())
}