/*
 * Copyright (c) 2022, 2023 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * https://oss.oracle.com/licenses/upl.
 */

package standalone

import (
	"context"
	"fmt"
	"github.com/onsi/gomega"
	"github.com/oracle/coherence-go-client/coherence"
	"github.com/oracle/coherence-go-client/coherence/extractors"
	"github.com/oracle/coherence-go-client/coherence/filters"
	"github.com/oracle/coherence-go-client/coherence/processors"
	. "github.com/oracle/coherence-go-client/test/utils"
	"log"
	"os"
	"sync"
	"testing"
)

// includeLongRunning indicates if to include long-running tests
const includeLongRunning = "INCLUDE_LONG_RUNNING"

var ctx = context.Background()

func TestBasicCrudOperationsVariousTypes(t *testing.T) {
	var (
		g       = gomega.NewWithT(t)
		err     error
		session *coherence.Session
	)

	session, err = GetSession()
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	defer session.Close()

	RunKeyValueTest[int, string](g, getNewNamedMap[int, string](g, session, "c1"), 1, "Tim")
	RunKeyValueTest[int, Person](g, getNewNamedMap[int, Person](g, session, "c2"), 1, Person{ID: 1, Name: "Tim"})
	RunKeyValueTest[int, float32](g, getNewNamedMap[int, float32](g, session, "c3"), 1, float32(1.123))
	RunKeyValueTest[int, float64](g, getNewNamedMap[int, float64](g, session, "c4"), 1, 1.123)
	RunKeyValueTest[int, int](g, getNewNamedMap[int, int](g, session, "c5"), 1, 1)
	RunKeyValueTest[int, int16](g, getNewNamedMap[int, int16](g, session, "c7"), 1, 10)
	RunKeyValueTest[int, int32](g, getNewNamedMap[int, int32](g, session, "c8"), 1, 1333)
	RunKeyValueTest[int, int64](g, getNewNamedMap[int, int64](g, session, "c9"), 1, 1333)
	RunKeyValueTest[int, bool](g, getNewNamedMap[int, bool](g, session, "c10"), 1, false)
	RunKeyValueTest[int, bool](g, getNewNamedMap[int, bool](g, session, "c11"), 1, true)
	RunKeyValueTest[int, byte](g, getNewNamedMap[int, byte](g, session, "c12"), 1, byte(22))
	RunKeyValueTest[string, Person](g, getNewNamedMap[string, Person](g, session, "c13"), "k1", Person{ID: 1, Name: "Tim"})
	RunKeyValueTest[string, string](g, getNewNamedMap[string, string](g, session, "c14"), "k1", "value1")
	RunKeyValueTest[int, Person](g, getNewNamedMap[int, Person](g, session, "c15"), 1,
		Person{ID: 1, Name: "Tim", HomeAddress: Address{Address1: "a1", Address2: "a2", City: "Perth", State: "WA", PostCode: 6028}})
	RunKeyValueTest[int, []string](g, getNewNamedMap[int, []string](g, session, "c16"), 1,
		[]string{"a", "b", "c"})
	RunKeyValueTest[int, map[int]string](g, getNewNamedMap[int, map[int]string](g, session, "c17"), 1,
		map[int]string{1: "one", 2: "two", 3: "three"})

	RunKeyValueTest[int, string](g, getNewNamedCache[int, string](g, session, "c1"), 1, "Tim")
	RunKeyValueTest[int, Person](g, getNewNamedCache[int, Person](g, session, "c2"), 1, Person{ID: 1, Name: "Tim"})
	RunKeyValueTest[int, float32](g, getNewNamedCache[int, float32](g, session, "c3"), 1, float32(1.123))
	RunKeyValueTest[int, float64](g, getNewNamedCache[int, float64](g, session, "c4"), 1, 1.123)
	RunKeyValueTest[int, int](g, getNewNamedCache[int, int](g, session, "c5"), 1, 1)
	RunKeyValueTest[int, int16](g, getNewNamedCache[int, int16](g, session, "c7"), 1, 10)
	RunKeyValueTest[int, int32](g, getNewNamedCache[int, int32](g, session, "c8"), 1, 1333)
	RunKeyValueTest[int, int64](g, getNewNamedCache[int, int64](g, session, "c9"), 1, 1333)
	RunKeyValueTest[int, bool](g, getNewNamedCache[int, bool](g, session, "c10"), 1, false)
	RunKeyValueTest[int, bool](g, getNewNamedCache[int, bool](g, session, "c11"), 1, true)
	RunKeyValueTest[int, byte](g, getNewNamedCache[int, byte](g, session, "c12"), 1, byte(22))
	RunKeyValueTest[string, Person](g, getNewNamedCache[string, Person](g, session, "c13"), "k1", Person{ID: 1, Name: "Tim"})
	RunKeyValueTest[string, string](g, getNewNamedCache[string, string](g, session, "c14"), "k1", "value1")
	RunKeyValueTest[int, Person](g, getNewNamedCache[int, Person](g, session, "c15"), 1,
		Person{ID: 1, Name: "Tim", HomeAddress: Address{Address1: "a1", Address2: "a2", City: "Perth", State: "WA", PostCode: 6028}})
	RunKeyValueTest[int, []string](g, getNewNamedCache[int, []string](g, session, "c16"), 1,
		[]string{"a", "b", "c"})
	RunKeyValueTest[int, map[int]string](g, getNewNamedCache[int, map[int]string](g, session, "c17"), 1,
		map[int]string{1: "one", 2: "two", 3: "three"})
}

// getNewNamedMap returns a map for a session and asserts err is nil
func getNewNamedMap[K comparable, V any](g *gomega.WithT, session *coherence.Session, name string) coherence.NamedMap[K, V] {
	namedMap, err := coherence.NewNamedMap[K, V](session, name)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	return namedMap
}

// getNewNamedCache returns a cache for a session and asserts err is nil
func getNewNamedCache[K comparable, V any](g *gomega.WithT, session *coherence.Session, name string) coherence.NamedCache[K, V] {
	namedCache, err := coherence.NewNamedCache[K, V](session, name)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	return namedCache
}

// TestBasicOperationsAgainstMapAndCache runs all tests against NamedMap and NamedCache
func TestBasicOperationsAgainstMapAndCache(t *testing.T) {
	g := gomega.NewWithT(t)
	session, err := GetSession()
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	defer session.Close()

	testCases := []struct {
		testName string
		nameMap  coherence.NamedMap[int, Person]
		test     func(t *testing.T, namedCache coherence.NamedMap[int, Person])
	}{
		{"NamedMapCrudTest", GetNamedMap[int, Person](g, session, "people-map"), RunTestBasicCrudOperations},
		{"NamedCacheCrudTest", GetNamedCache[int, Person](g, session, "people-cache"), RunTestBasicCrudOperations},
		{"NamedMapRunTestGetOrDefault", GetNamedMap[int, Person](g, session, "get-or-default-map"), RunTestGetOrDefault},
		{"NamedCacheRunTestGetOrDefault", GetNamedCache[int, Person](g, session, "get-or-default-cache"), RunTestGetOrDefault},
		{"NamedMapRunTestContainsKey", GetNamedMap[int, Person](g, session, "contains-key-map"), RunTestContainsKey},
		{"NamedCacheRunTestContainsKey", GetNamedCache[int, Person](g, session, "contains-key-cache"), RunTestContainsKey},
		{"NamedMapRunTestPutIfAbsent", GetNamedMap[int, Person](g, session, "put-if-absent-map"), RunTestPutIfAbsent},
		{"NamedCacheRunTestPutIfAbsent", GetNamedCache[int, Person](g, session, "put-of-absent-cache"), RunTestPutIfAbsent},
		{"NamedMapRunTestClearAndIsEmpty", GetNamedMap[int, Person](g, session, "clear-map"), RunTestClearAndIsEmpty},
		{"NamedCacheRunTestClearAndIsEmpty", GetNamedCache[int, Person](g, session, "clear-cache"), RunTestClearAndIsEmpty},
		{"NamedMapRunTestTruncateAndDestroy", GetNamedMap[int, Person](g, session, "truncate-map"), RunTestTruncateAndDestroy},
		{"NamedCacheRunTestTruncateAndDestroy", GetNamedCache[int, Person](g, session, "truncate-cache"), RunTestTruncateAndDestroy},
		{"NamedMapRunTestReplace", GetNamedMap[int, Person](g, session, "replace-map"), RunTestReplace},
		{"NamedCacheRunTestReplace", GetNamedCache[int, Person](g, session, "replace-cache"), RunTestReplace},
		{"NamedMapRunTestReplaceMapping", GetNamedMap[int, Person](g, session, "replace-mapping-map"), RunTestReplaceMapping},
		{"NamedCacheRunTestReplaceMapping", GetNamedCache[int, Person](g, session, "replace-mapping-cache"), RunTestReplaceMapping},
		{"NamedMapRunTestRemoveMapping", GetNamedMap[int, Person](g, session, "remove-mapping-map"), RunTestRemoveMapping},
		{"NamedCacheRunTestRemoveMapping", GetNamedCache[int, Person](g, session, "remove-mapping-cache"), RunTestRemoveMapping},
		{"NamedMapRunTestPutAll", GetNamedMap[int, Person](g, session, "remove-mapping-map"), RunTestPutAll},
		{"NamedCacheRunTestPutAll", GetNamedCache[int, Person](g, session, "remove-mapping-cache"), RunTestPutAll},
		{"NamedMapRunTestContainsValue", GetNamedMap[int, Person](g, session, "contains-value-map"), RunTestContainsValue},
		{"NamedCacheRunTestContainsValue", GetNamedCache[int, Person](g, session, "contains-value-cache"), RunTestContainsValue},
		{"NamedMapRunTestContainsEntry", GetNamedMap[int, Person](g, session, "contains-entry-map"), RunTestContainsEntry},
		{"NamedCacheRunTestContainsEntry", GetNamedCache[int, Person](g, session, "contains-entry-cache"), RunTestContainsEntry},
		{"NamedMapRunTestValuesFilter", GetNamedMap[int, Person](g, session, "values-filter-map"), RunTestValuesFilter},
		{"NamedCacheRunTestValuesFilter", GetNamedCache[int, Person](g, session, "values-filter-cache"), RunTestValuesFilter},
		{"NamedMapRunTestEntrySetFilter", GetNamedMap[int, Person](g, session, "entryset-filter-map"), RunTestEntrySetFilter},
		{"NamedCacheRunTestEntrySetFilter", GetNamedCache[int, Person](g, session, "entryset-filter-cache"), RunTestEntrySetFilter},
		{"NamedMapRunTestKeySetFilter", GetNamedMap[int, Person](g, session, "keyset-map"), RunTestKeySetFilter},
		{"NamedCacheRunTestKeySetFilter", GetNamedCache[int, Person](g, session, "keyset-cache"), RunTestKeySetFilter},
		{"NamedMapRunTestGetAll", GetNamedMap[int, Person](g, session, "getall-filter-map"), RunTestGetAll},
		{"NamedCacheRunTestGetAll", GetNamedCache[int, Person](g, session, "getall-filter-cache"), RunTestGetAll},
		{"NamedMapRunTestInvokeAll", GetNamedMap[int, Person](g, session, "invokeall-keys-map"), RunTestInvokeAllKeys},
		{"NamedCacheRunTestInvokeAll", GetNamedCache[int, Person](g, session, "invokeall-keys-cache"), RunTestInvokeAllKeys},
		{"NamedMapRunTestKeySet", GetNamedMap[int, Person](g, session, "keyset-map"), RunTestKeySet},
		{"NamedCacheRunTestKeySet", GetNamedCache[int, Person](g, session, "keyset-cache"), RunTestKeySet},
		{"NamedMapRunTestEntrySet", GetNamedMap[int, Person](g, session, "entryset-map"), RunTestEntrySet},
		{"NamedCacheRunTestEntrySet", GetNamedCache[int, Person](g, session, "entryset-cache"), RunTestEntrySet},
		{"NamedMapRunTestValues", GetNamedMap[int, Person](g, session, "values-map"), RunTestValues},
		{"NamedCacheRunTestValues", GetNamedCache[int, Person](g, session, "values-cache"), RunTestValues},
		{"NamedMapRunTestEntrySetGoRoutines", GetNamedMap[int, Person](g, session, "go-map"), RunTestEntrySetGoRoutines},
		{"NamedCacheRunTestEntrySetGoRoutines", GetNamedCache[int, Person](g, session, "go-cache"), RunTestEntrySetGoRoutines},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			tc.test(t, tc.nameMap)
		})
	}
}

func RunTestBasicCrudOperations(t *testing.T, namedMap coherence.NamedMap[int, Person]) {
	var (
		g         = gomega.NewWithT(t)
		result    *Person
		oldValue  *Person
		err       error
		person1   = Person{ID: 1, Name: "Tim"}
		oldPerson *Person
	)

	oldPerson, err = namedMap.Put(ctx, person1.ID, person1)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(oldPerson).To(gomega.BeNil())

	result, err = namedMap.Get(ctx, person1.ID)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(result).To(gomega.Not(gomega.BeNil()))
	AssertPersonResult(g, *result, person1)

	// update the name to "Timothy"
	person1.Name = "Timothy"
	oldValue, err = namedMap.Put(ctx, person1.ID, person1)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(oldValue).To(gomega.Not(gomega.BeNil()))
	g.Expect(oldValue.Name).To(gomega.Equal("Tim"))

	result, err = namedMap.Get(ctx, person1.ID)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	AssertPersonResult(g, *result, person1)

	oldValue, err = namedMap.Remove(ctx, person1.ID)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	AssertPersonResult(g, *oldValue, person1)
}

func TestTestMultipleCallsToNamedMap(t *testing.T) {
	var (
		g            = gomega.NewWithT(t)
		err          error
		person1      = Person{ID: 1, Name: "Tim"}
		personValue1 *Person
		personValue2 *Person
	)

	session, err := GetSession()
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	defer session.Close()

	namedMap1, err := coherence.NewNamedMap[int, Person](session, "map-1")
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	err = namedMap1.Clear(ctx)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// retrieve the named map again, should return the same one
	namedMap2, err := coherence.NewNamedMap[int, Person](session, "map-1")
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	err = namedMap2.Clear(ctx)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	g.Expect(namedMap2).To(gomega.Equal(namedMap1))

	_, err = namedMap1.Put(ctx, person1.ID, person1)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	personValue1, err = namedMap1.Get(ctx, person1.ID)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	personValue2, err = namedMap2.Get(ctx, person1.ID)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	g.Expect(*personValue1).To(gomega.Equal(*personValue2))

	namedMap3, err := coherence.NewNamedMap[int, Person](session, "map-2")
	g.Expect(err).NotTo(gomega.HaveOccurred())

	AssertSize(g, namedMap3, 0)

	// try and retrieve a NamedMap that is for the same cache but different type, this should cause error
	_, err = coherence.NewNamedMap[int, string](session, "map-2")
	g.Expect(err).To(gomega.HaveOccurred())
}

func RunTestGetOrDefault(t *testing.T, namedMap coherence.NamedMap[int, Person]) {
	var (
		g        = gomega.NewWithT(t)
		person1  = Person{ID: 10, Name: "John"}
		result   *Person
		oldValue *Person
		err      error
	)

	// should be able to get default when Value does not exist
	result, err = namedMap.GetOrDefault(ctx, 10, person1)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(result).To(gomega.Not(gomega.BeNil()))
	g.Expect(result.ID).To(gomega.Equal(person1.ID))
	g.Expect(result.Name).To(gomega.Equal(person1.Name))

	// put a Value and this Value should be retrieved in subsequent calls
	oldValue, err = namedMap.Put(ctx, 10, person1)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(oldValue).To(gomega.BeNil())

	result, err = namedMap.GetOrDefault(ctx, 10, Person{ID: 111, Name: "Not this one"})
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(result).To(gomega.Not(gomega.BeNil()))
	AssertPersonResult(g, *result, person1)
}

func RunTestContainsKey(t *testing.T, namedMap coherence.NamedMap[int, Person]) {
	var (
		g        = gomega.NewWithT(t)
		person1  = Person{ID: 1, Name: "Tim"}
		oldValue *Person
		err      error
		found    bool
	)

	// should not contain a Key for an entry that does not exist
	found, err = namedMap.ContainsKey(ctx, person1.ID)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(found).To(gomega.BeFalse())

	// add a new entry
	oldValue, err = namedMap.Put(ctx, person1.ID, person1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(oldValue).To(gomega.BeNil())

	// ensure that the ContainsKey is true
	found, err = namedMap.ContainsKey(ctx, person1.ID)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(found).To(gomega.BeTrue())
}

func RunTestContainsEntry(t *testing.T, namedMap coherence.NamedMap[int, Person]) {
	var (
		g        = gomega.NewWithT(t)
		person1  = Person{ID: 1, Name: "Tim"}
		person2  = Person{ID: 2, Name: "Tim2"}
		err      error
		found    bool
		oldValue *Person
	)

	// should not contain a entry for an entry that does not exist
	found, err = namedMap.ContainsEntry(ctx, person1.ID, person1)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(found).To(gomega.BeFalse())

	// add a new entry
	oldValue, err = namedMap.Put(ctx, person1.ID, person1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(oldValue).To(gomega.BeNil())

	// ensure that the ContainsEntry is true
	found, err = namedMap.ContainsEntry(ctx, person1.ID, person1)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(found).To(gomega.BeTrue())

	// ensure that the if the entry is different, then this should fail
	found, err = namedMap.ContainsEntry(ctx, person1.ID, person2)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(found).To(gomega.BeFalse())

	// ensure that the if the key is different, then this should fail
	found, err = namedMap.ContainsEntry(ctx, person2.ID, person2)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(found).To(gomega.BeFalse())
}

func RunTestContainsValue(t *testing.T, namedMap coherence.NamedMap[int, Person]) {
	var (
		g        = gomega.NewWithT(t)
		person1  = Person{ID: 1, Name: "Tim"}
		err      error
		found    bool
		oldValue *Person
	)

	// should not contain a value with no cache entries
	found, err = namedMap.ContainsValue(ctx, person1)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(found).To(gomega.BeFalse())

	// add a new entry
	oldValue, err = namedMap.Put(ctx, person1.ID, person1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(oldValue).To(gomega.BeNil())

	// ensure that the ContainsKey is true
	found, err = namedMap.ContainsValue(ctx, person1)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(found).To(gomega.BeTrue())
}

func RunTestPutIfAbsent(t *testing.T, namedMap coherence.NamedMap[int, Person]) {
	var (
		g        = gomega.NewWithT(t)
		person1  = Person{ID: 1, Name: "Tim"}
		person2  = Person{ID: 1, Name: "Timothy"}
		err      error
		found    bool
		result   *Person
		oldValue *Person
	)

	// put if absent should return nil if Value is not present
	oldValue, err = namedMap.PutIfAbsent(ctx, person1.ID, person1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(oldValue).To(gomega.BeNil())

	// Key should exist
	found, err = namedMap.ContainsKey(ctx, person1.ID)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(found).To(gomega.BeTrue())

	// try to put an updated person2. The entry should not be updated
	// and the existing entry should be returned
	result, err = namedMap.PutIfAbsent(ctx, person2.ID, person2)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(result).To(gomega.Not(gomega.BeNil()))

	// assert the Value returned is the existing entry
	AssertPersonResult(g, *result, person1)

	// ensure a Get for the person1.ID returns the original person
	// and not the attempted update
	result, err = namedMap.Get(ctx, person1.ID)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(result).To(gomega.Not(gomega.BeNil()))
	AssertPersonResult(g, *result, person1)
}

func RunTestClearAndIsEmpty(t *testing.T, namedMap coherence.NamedMap[int, Person]) {
	var (
		g        = gomega.NewWithT(t)
		person1  = Person{ID: 1, Name: "Tim"}
		err      error
		isEmpty  bool
		oldValue *Person
	)

	isEmpty, err = namedMap.IsEmpty(ctx)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(isEmpty).To(gomega.BeTrue())

	_, err = namedMap.Put(ctx, person1.ID, person1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(oldValue).To(gomega.BeNil())

	isEmpty, err = namedMap.IsEmpty(ctx)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(isEmpty).To(gomega.BeFalse())

	err = namedMap.Clear(ctx)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	isEmpty, err = namedMap.IsEmpty(ctx)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(isEmpty).To(gomega.BeTrue())
}

func RunTestTruncateAndDestroy(t *testing.T, namedMap coherence.NamedMap[int, Person]) {
	var (
		g        = gomega.NewWithT(t)
		person1  = Person{ID: 1, Name: "Tim"}
		err      error
		oldValue *Person
	)

	_, err = namedMap.Put(ctx, person1.ID, person1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(oldValue).To(gomega.BeNil())

	err = namedMap.Truncate(ctx)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	err = namedMap.Destroy(ctx)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
}

func RunTestReplace(t *testing.T, namedMap coherence.NamedMap[int, Person]) {
	var (
		g             = gomega.NewWithT(t)
		person1       = Person{ID: 1, Name: "Tim"}
		personReplace = Person{ID: 1, Name: "Timothy"}
		err           error
		oldValue      *Person
	)

	// no Value for Key exists so will not replace or return old Value
	oldValue, err = namedMap.Replace(ctx, 1, person1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(oldValue).To(gomega.BeNil())

	// add an entry that we will replace further down
	oldValue, err = namedMap.Put(ctx, 1, person1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(oldValue).To(gomega.BeNil())

	AssertSize(g, namedMap, 1)

	// this should work as it's mapped to any Value
	oldValue, err = namedMap.Replace(ctx, 1, personReplace)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(oldValue).To(gomega.Not(gomega.BeNil()))
	AssertPersonResult(g, *oldValue, person1)
}

func RunTestReplaceMapping(t *testing.T, namedMap coherence.NamedMap[int, Person]) {
	var (
		g             = gomega.NewWithT(t)
		person1       = Person{ID: 1, Name: "Tim"}
		personReplace = Person{ID: 1, Name: "Timothy"}
		personNew     = Person{ID: 1, Name: "Timothy Jones"}
		err           error
		result        bool
		personValue   *Person
		oldValue      *Person
	)

	// no Value for Key exists so will not replace and should return false
	result, err = namedMap.ReplaceMapping(ctx, 1, personReplace, person1)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(result).To(gomega.Equal(false))

	// add an entry that we will replace further down
	oldValue, err = namedMap.Put(ctx, 1, person1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(oldValue).To(gomega.BeNil())

	AssertSize(g, namedMap, 1)

	// value exists but doesn't match so should return false
	result, err = namedMap.ReplaceMapping(ctx, 1, personReplace, person1)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(result).To(gomega.Equal(false))

	// now try replacing where exists and matches
	result, err = namedMap.ReplaceMapping(ctx, 1, person1, personNew)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(result).To(gomega.Equal(true))

	// get the value and check that matches
	personValue, err = namedMap.Get(ctx, 1)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(personValue).To(gomega.Not(gomega.BeNil()))
	g.Expect(*personValue).To(gomega.Equal(personNew))
}

func RunTestRemoveMapping(t *testing.T, namedMap coherence.NamedMap[int, Person]) {
	var (
		g       = gomega.NewWithT(t)
		person1 = Person{ID: 1, Name: "Tim"}
		person2 = Person{ID: 1, Name: "Tim2"}
		err     error
		removed bool
	)

	// remove a mapping that doesn't exist
	removed, err = namedMap.RemoveMapping(ctx, 1, person1)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(removed).Should(gomega.Equal(false))

	// add a Key with a Value that will not match
	_, err = namedMap.Put(ctx, 1, person2)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	AssertSize(g, namedMap, 1)

	// remove a mapping that doesn't match
	removed, err = namedMap.RemoveMapping(ctx, 1, person1)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(removed).Should(gomega.Equal(false))
	AssertSize(g, namedMap, 1)

	// set the Key to a Value that will match
	_, err = namedMap.Put(ctx, 1, person1)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	removed, err = namedMap.RemoveMapping(ctx, 1, person1)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(removed).Should(gomega.Equal(true))
	AssertSize(g, namedMap, 0)
}

var peopleData = map[int]Person{
	1: {ID: 1, Name: "Tim", Age: 50},
	2: {ID: 2, Name: "Andrew", Age: 44},
	3: {ID: 3, Name: "Helen", Age: 20},
	4: {ID: 4, Name: "Alexa", Age: 12},
}

func RunTestPutAll(t *testing.T, namedMap coherence.NamedMap[int, Person]) {
	var (
		g     = gomega.NewWithT(t)
		err   error
		found bool
		size  int
	)

	err = namedMap.PutAll(ctx, peopleData)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	size, err = namedMap.Size(ctx)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(size).To(gomega.Equal(len(peopleData)))

	for k := range peopleData {
		found, err = namedMap.ContainsKey(ctx, k)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(found).To(gomega.BeTrue())
	}

}

func RunTestValuesFilter(t *testing.T, namedMap coherence.NamedMap[int, Person]) {
	var (
		g       = gomega.NewWithT(t)
		results = make([]Person, 0)
	)

	// populate the cache
	populatePeople(g, namedMap)

	ch := namedMap.ValuesFilter(ctx, filters.Always())
	for se := range ch {
		g.Expect(se.Err).ShouldNot(gomega.HaveOccurred())
		g.Expect(se.Value).ShouldNot(gomega.BeNil())
		results = append(results, se.Value)
	}
	g.Expect(len(results)).To(gomega.Equal(len(peopleData)))

	// reset the results
	results = make([]Person, 0)

	ch2 := namedMap.ValuesFilter(ctx, filters.GreaterEqual(extractors.Extract[int]("age"), 20))
	for se := range ch2 {
		g.Expect(se.Err).ShouldNot(gomega.HaveOccurred())
		g.Expect(se.Value).ShouldNot(gomega.BeNil())
		results = append(results, se.Value)
	}

	g.Expect(len(results)).To(gomega.Equal(3))
}

func RunTestEntrySetFilter(t *testing.T, namedMap coherence.NamedMap[int, Person]) {
	var (
		g       = gomega.NewWithT(t)
		results = make([]Person, 0)
	)

	// populate the cache
	populatePeople(g, namedMap)

	ch := namedMap.EntrySetFilter(ctx, filters.Always())
	for se := range ch {
		g.Expect(se.Err).ShouldNot(gomega.HaveOccurred())
		g.Expect(se.Value).ShouldNot(gomega.BeNil())
		results = append(results, se.Value)
	}
	g.Expect(len(results)).To(gomega.Equal(len(peopleData)))

	// reset the results
	results = make([]Person, 0)

	ch2 := namedMap.EntrySetFilter(ctx, filters.GreaterEqual(extractors.Extract[int]("age"), 20))
	for se := range ch2 {
		g.Expect(se.Err).ShouldNot(gomega.HaveOccurred())
		g.Expect(se.Value).ShouldNot(gomega.BeNil())
		results = append(results, se.Value)
	}

	g.Expect(len(results)).To(gomega.Equal(3))
}

func RunTestKeySetFilter(t *testing.T, namedMap coherence.NamedMap[int, Person]) {
	var (
		g       = gomega.NewWithT(t)
		results = make([]int, 0)
	)

	// populate the cache
	populatePeople(g, namedMap)

	ch := namedMap.KeySetFilter(ctx, filters.Always())
	for se := range ch {
		g.Expect(se.Err).ShouldNot(gomega.HaveOccurred())
		g.Expect(se.Key).ShouldNot(gomega.BeNil())
		results = append(results, se.Key)
	}
	g.Expect(len(results)).To(gomega.Equal(len(peopleData)))

	// reset the results
	results = make([]int, 0)

	ch2 := namedMap.KeySetFilter(ctx, filters.GreaterEqual(extractors.Extract[int]("age"), 20))
	for se := range ch2 {
		g.Expect(se.Err).ShouldNot(gomega.HaveOccurred())
		g.Expect(se.Key).ShouldNot(gomega.BeNil())
		results = append(results, se.Key)
	}

	g.Expect(len(results)).To(gomega.Equal(3))
}

func RunTestGetAll(t *testing.T, namedMap coherence.NamedMap[int, Person]) {
	var (
		g       = gomega.NewWithT(t)
		results = make([]int, 0)
	)

	// populate the cache
	populatePeople(g, namedMap)

	ch := namedMap.GetAll(ctx, []int{1, 3})
	for se := range ch {
		g.Expect(se.Err).ShouldNot(gomega.HaveOccurred())
		g.Expect(se.Key).ShouldNot(gomega.BeNil())
		results = append(results, se.Key)
	}
	g.Expect(len(results)).To(gomega.Equal(2))

	// reset the results
	results = make([]int, 0)

	ch2 := namedMap.GetAll(ctx, []int{333})
	for se := range ch2 {
		g.Expect(se.Err).ShouldNot(gomega.HaveOccurred())
		g.Expect(se.Key).ShouldNot(gomega.BeNil())
		results = append(results, se.Key)
	}
	g.Expect(len(results)).To(gomega.Equal(0))
}

func RunTestInvokeAllKeys(t *testing.T, namedMap coherence.NamedMap[int, Person]) {
	var (
		g       = gomega.NewWithT(t)
		results = make([]int, 0)
	)

	// populate the cache
	populatePeople(g, namedMap)

	// run a processor to increment the age of people with key 1 and 2
	ch := coherence.InvokeAllKeys[int, Person, int](ctx, namedMap, []int{1, 2}, processors.Increment("age", 1, true))

	for se := range ch {
		g.Expect(se.Err).ShouldNot(gomega.HaveOccurred())
		g.Expect(se.Value).ShouldNot(gomega.BeNil())
		results = append(results, se.Value)
	}
	g.Expect(len(results)).To(gomega.Equal(2))
	// the results are the keys that were updates
	g.Expect(containsValue[int](results, 1)).Should(gomega.BeTrue())
	g.Expect(containsValue[int](results, 2)).Should(gomega.BeTrue())

	// reset and run for filter
	results = make([]int, 0)

	// run a processor to increment the age of people who are older than 1
	ch2 := coherence.InvokeAllFilter[int, Person, int](ctx, namedMap,
		filters.Greater(extractors.Extract[int]("age"), 1), processors.Increment("age", 1, true))

	for se := range ch2 {
		g.Expect(se.Err).ShouldNot(gomega.HaveOccurred())
		g.Expect(se.Value).ShouldNot(gomega.BeNil())
		results = append(results, se.Value)
	}

	// should match all entries
	g.Expect(len(results)).To(gomega.Equal(4))
}

func RunTestKeySet(t *testing.T, namedMap coherence.NamedMap[int, Person]) {
	var (
		g           = gomega.NewWithT(t)
		results     = make([]int, 0)
		insertCount = 400_000
		result      *int
	)

	if !includeLongRunningTests() {
		t.Log("Skipping long running tests")
		return
	}

	err := namedMap.Clear(ctx)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	// test with empty cache to ensure we receive the ErrDone straight away
	iter := namedMap.KeySet(ctx)
	_, err = iter.Next()
	g.Expect(err).To(gomega.Equal(coherence.ErrDone))

	// test with single entry which will force only 1 page to be returned
	_, err = namedMap.Put(ctx, 1, Person{ID: 1, Name: "Tim", Age: 54})
	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	iter = namedMap.KeySet(ctx)
	result, err = iter.Next()
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(*result).To(gomega.Equal(1))
	_, err = iter.Next()
	g.Expect(err).To(gomega.Equal(coherence.ErrDone))

	// test with enough data to use multiple requests
	addManyPeople(g, namedMap, 1, insertCount)

	iter = namedMap.KeySet(ctx)
	for {
		result, err = iter.Next()

		if err == coherence.ErrDone {
			break
		}

		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		results = append(results, *result)
	}

	g.Expect(len(results)).To(gomega.Equal(insertCount))
	_ = namedMap.Clear(ctx)
}

func RunTestEntrySet(t *testing.T, namedMap coherence.NamedMap[int, Person]) {
	var (
		g           = gomega.NewWithT(t)
		results     = make([]coherence.Entry[int, Person], 0)
		err         error
		value       *coherence.Entry[int, Person]
		insertCount = 400_000
	)

	if !includeLongRunningTests() {
		t.Log("Skipping long running tests")
		return
	}

	err = namedMap.Clear(ctx)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	// test with empty cache to ensure we receive the ErrDone straight away
	iter := namedMap.EntrySet(ctx)
	_, err = iter.Next()
	g.Expect(err).To(gomega.Equal(coherence.ErrDone))

	// test with single entry which will force only 1 page to be returned
	_, err = namedMap.Put(ctx, 1, Person{ID: 1, Name: "Tim", Age: 54})
	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	iter = namedMap.EntrySet(ctx)
	value, err = iter.Next()
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(value.Value.ID).To(gomega.Equal(1))
	_, err = iter.Next()
	g.Expect(err).To(gomega.Equal(coherence.ErrDone))

	// test with enough data to use multiple requests
	addManyPeople(g, namedMap, 1, insertCount)

	iter = namedMap.EntrySet(ctx)

	for {
		value, err = iter.Next()
		if err == coherence.ErrDone {
			break
		}
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(value.Key).ShouldNot(gomega.BeNil())
		g.Expect(value.Value).ShouldNot(gomega.BeNil())
		results = append(results, *value)
	}

	g.Expect(len(results)).To(gomega.Equal(insertCount))
	_ = namedMap.Clear(ctx)
}

func RunTestValues(t *testing.T, namedMap coherence.NamedMap[int, Person]) {
	var (
		g           = gomega.NewWithT(t)
		results     = make([]Person, 0)
		err         error
		value       *Person
		insertCount = 400_000
	)

	if !includeLongRunningTests() {
		t.Log("Skipping long running tests")
		return
	}

	err = namedMap.Clear(ctx)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	// test with empty cache to ensure we receive the ErrDone straight away
	iter := namedMap.Values(ctx)
	_, err = iter.Next()
	g.Expect(err).To(gomega.Equal(coherence.ErrDone))

	// test with single entry which will force only 1 page to be returned
	_, err = namedMap.Put(ctx, 1, Person{ID: 1, Name: "Tim", Age: 54})
	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	iter = namedMap.Values(ctx)
	value, err = iter.Next()
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(value.ID).To(gomega.Equal(1))
	_, err = iter.Next()
	g.Expect(err).To(gomega.Equal(coherence.ErrDone))

	// test with enough data to use multiple requests
	addManyPeople(g, namedMap, 1, insertCount)

	iter = namedMap.Values(ctx)

	for {
		value, err = iter.Next()
		if err == coherence.ErrDone {
			break
		}
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(value.ID).Should(gomega.BeNumerically(">", 0))
		results = append(results, *value)
	}

	g.Expect(len(results)).To(gomega.Equal(insertCount))
	_ = namedMap.Clear(ctx)
}

func RunTestEntrySetGoRoutines(t *testing.T, namedMap coherence.NamedMap[int, Person]) {
	var (
		g           = gomega.NewWithT(t)
		err         error
		insertCount = 500_000
		wg          sync.WaitGroup
		count       [4]int
	)

	if !includeLongRunningTests() {
		t.Log("Skipping long running tests")
		return
	}

	err = namedMap.Clear(ctx)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	countLen := len(count)
	wg.Add(countLen)

	// test with enough data to use multiple requests
	addManyPeople(g, namedMap, 1, insertCount)

	iter := namedMap.EntrySet(ctx)

	// run 4 go routines reading data off the same iterator to ensure the values are
	// correctly iterated through.
	for i := 0; i < countLen; i++ {
		go func(index int) {
			defer wg.Done()
			for {
				value1, err1 := iter.Next()
				if err1 == coherence.ErrDone {
					break
				}
				g.Expect(err1).ShouldNot(gomega.HaveOccurred())
				g.Expect(value1.Key).ShouldNot(gomega.BeNil())
				g.Expect(value1.Value).ShouldNot(gomega.BeNil())
				count[index]++
			}
			t.Log("Index", index, "count", count[index])
		}(i)
	}

	wg.Wait()

	g.Expect(count[0] + count[1] + count[2] + count[3]).To(gomega.Equal(insertCount))

}

func addManyPeople(g *gomega.WithT, namedMap coherence.NamedMap[int, Person], startKey, count int) {
	var (
		err error
	)

	buffer := make(map[int]Person, 0)
	for i := startKey; i < count+startKey; i++ {
		buffer[i] = Person{
			ID:   i,
			Name: fmt.Sprintf("Person %d", i),
			Age:  10 + i%50,
		}

		if i%1000 == 0 {
			err = namedMap.PutAll(ctx, buffer)
			g.Expect(err).ShouldNot(gomega.HaveOccurred())
			buffer = make(map[int]Person, 0)
		}
	}

	// write any leftover buffer
	if len(buffer) > 0 {
		err = namedMap.PutAll(ctx, buffer)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// populatePeople populates the people map
func populatePeople(g *gomega.WithT, namedMap coherence.NamedMap[int, Person]) {
	var (
		err  error
		size int
	)
	err = namedMap.PutAll(ctx, peopleData)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	size, err = namedMap.Size(ctx)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(size).To(gomega.Equal(len(peopleData)))
}

// containsValue returns true if the value is contains within the array
func containsValue[V comparable](values []V, value V) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func includeLongRunningTests() bool {
	if val := os.Getenv(includeLongRunning); val != "" {
		return true
	}
	return false
}