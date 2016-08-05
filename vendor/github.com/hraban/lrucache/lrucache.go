// Copyright © 2012, 2013 Lrucache contributors, see AUTHORS file
//
// The license for this file is described in the LICENSE file.

// Light-weight in-memory LRU (object) cache library for Go.
//
// To use this library, first create a cache:
//
//      c := lrucache.New(1234)
//
// Then, optionally, define a type that implements some of the interfaces:
//
//      type cacheableInt int
//
//      func (i cacheableInt) OnPurge(why lrucache.PurgeReason) {
//          fmt.Printf("Purging %d\n", i)
//      }
//
// Finally:
//
//     for i := 0; i < 2000; i++ {
//         c.Set(strconv.Itoa(i), cacheableInt(i))
//     }
//
// This will generate the following output:
//
//     Purging 0
//     Purging 1
//     ...
//     Purging 764
//     Purging 765
//
// Note:
//
// * The unit of item sizes is not defined; whatever it is, once the sum
// exceeds the maximum cache size, elements start getting purged until it
// drops below the threshold again.
//
// * These integers are passed by value. Caching pointers is, of course, Okay,
// but be careful when caching a memory location that holds two different
// values at different points in time; updating the value of a pointer after
// caching it will change the cached value.
//
package lrucache

import (
	"errors"
)

// A function that generates a fresh entry on "cache miss". See the OnMiss
// method.
type OnMissHandler func(string) (Cacheable, error)

type Cache struct {
	maxSize int64
	size    int64
	entries map[string]*cacheEntry
	// Cache operations are pushed down this channel to the main cache loop
	opChan chan operation
	// most recently used entry
	mostRU *cacheEntry
	// least recently used entry
	leastRU *cacheEntry
	// If not nil, invoked for every cache miss.
	onMiss OnMissHandler
}

// Anything can be cached!
type Cacheable interface{}

// Optional interface for cached objects. If this interface is not implemented,
// an element is assumed to have size 1.
type SizeAware interface {
	// See Cache.MaxSize() for an explanation of the semantics. Please report a
	// constant size; the cache does not expect objects to change size while
	// they are cached. Items are trusted to report their own size accurately.
	Size() int64
}

func getSize(x Cacheable) int64 {
	if s, ok := x.(SizeAware); ok {
		return s.Size()
	}
	return 1
}

// Reasons for a cached element to be deleted from the cache
type PurgeReason int

const (
	// Cache is growing too large and this is the least used item
	CACHEFULL PurgeReason = iota
	// This item was explicitly deleted using Cache.Delete(id)
	EXPLICITDELETE
	// A new element with the same key is stored (usually indicates an update)
	KEYCOLLISION
)

// Optional interface for cached objects
type NotifyPurge interface {
	// Called once when the element is purged from cache. The argument
	// indicates why.
	//
	// Example use-case: a session cache where sessions are not stored in a
	// database until they are purged from the memory cache. As long as the
	// memory cache is large enough to hold all of them, they expire before the
	// cache grows too large and no database connection is ever needed. This
	// OnPurge implementation would store items to a database iff reason ==
	// CACHEFULL.
	//
	// Called from within a private goroutine, but never called concurrently
	// with other elements' OnPurge(). The entire cache is blocked until this
	// function returns. By all means, feel free to launch a fresh goroutine
	// and return immediately.
	OnPurge(why PurgeReason)
}

// Requests that are passed to the cache managing goroutine
type operation struct {
	// Cache is explicitly passed with every request so the main loop does not
	// need to keep a reference to the cache around. This allows garbage
	// collection to kick in (and eventually end the main loop through a
	// finalizer) when the last external reference to the cache goes out of
	// scope.
	c *Cache
	// any of the types defined below starting with req...
	req interface{}
}

type reqSet struct {
	id      string
	payload Cacheable
}

// Reply to a Get request
type replyGet struct {
	val Cacheable
	err error
}

type reqGet struct {
	id string
	// If the key is found the value is pushed down this channel after which it
	// is closed immediately. If the value is not found, OnMiss is called. If
	// that does not work (OnMiss is not defined, or it returns nil) the
	// error is set to ErrNotFound. Otherwise the result is set to whatever
	// OnMiss returned. One way or another, exactly one value is pushed down
	// this channel, after which it is closed.
	reply chan<- replyGet
}

type reqDelete string

type reqOnMissFunc OnMissHandler

type reqMaxSize int64

type reqGetSize chan<- int64

type cacheEntry struct {
	payload Cacheable
	id      string
	// youngest older entry (age being usage) (DLL pointer)
	older *cacheEntry
	// oldest younger entry (age being usage) (DLL pointer)
	younger *cacheEntry
}

// Only call c.OnPurge() if c implements NotifyPurge.
func safeOnPurge(c Cacheable, why PurgeReason) {
	if t, ok := c.(NotifyPurge); ok {
		t.OnPurge(why)
	}
	return
}

func removeEntry(c *Cache, e *cacheEntry) {
	delete(c.entries, e.id)
	if e.older == nil {
		c.leastRU = e.younger
	} else {
		e.older.younger = e.younger
	}
	if e.younger == nil {
		c.mostRU = e.older
	} else {
		e.younger.older = e.older
	}
	c.size -= getSize(e.payload)
	return
}

// Purge the least recently used from the cache
func purgeLRU(c *Cache) {
	safeOnPurge(c.leastRU.payload, CACHEFULL)
	removeEntry(c, c.leastRU)
	return
}

// Trim the cache until its size <= max size
func trimCache(c *Cache) {
	if c.maxSize <= 0 {
		return
	}
	for c.size > c.maxSize {
		purgeLRU(c)
	}
	return
}

// Not safe for use in concurrent goroutines
func directSet(c *Cache, req reqSet) {
	// Overwrite old entry
	if old, ok := c.entries[req.id]; ok {
		safeOnPurge(old.payload, KEYCOLLISION)
		removeEntry(c, old)
	}
	e := cacheEntry{payload: req.payload, id: req.id}
	c.entries[req.id] = &e
	size := getSize(e.payload)
	if size == 0 {
		return
	}
	if c.leastRU == nil { // aka "if this is the first entry..."
		// init DLL
		c.leastRU = &e
		c.mostRU = &e
		e.younger = nil
		e.older = nil
	} else {
		// e is younger than the old "most recently used"
		c.mostRU.younger = &e
		e.older = c.mostRU
		c.mostRU = &e
	}
	c.size += size
	trimCache(c)
	return
}

// Not safe for use in concurrent goroutines
func directDelete(c *Cache, req reqDelete) {
	id := string(req)
	e, ok := c.entries[id]
	if ok {
		safeOnPurge(e.payload, EXPLICITDELETE)
		if getSize(e.payload) != 0 {
			removeEntry(c, e)
		}
	}
	return
}

// Handle a cache miss from outside the main goroutine
func handleCacheMiss(c *Cache, req reqGet) {
	var val Cacheable
	var err error = ErrNotFound
	if c.onMiss != nil {
		val, err = c.onMiss(req.id)
		if err == nil {
			if val != nil {
				c.Set(req.id, val)
			} else {
				err = ErrNotFound
			}
		}
	}
	req.reply <- replyGet{val, err}
	close(req.reply)
	return
}

// Not safe for use in concurrent goroutines
func directGet(c *Cache, req reqGet) {
	e, ok := c.entries[req.id]
	if !ok {
		go handleCacheMiss(c, req)
		return
	}
	req.reply <- replyGet{e.payload, nil}
	close(req.reply)
	if e.younger == nil {
		// I'm already the fresh kid on the block
		return
	}
	// Put element at the start of the LRU list
	if e.older != nil {
		e.older.younger = e.younger // the only reason this is a *D*LL
	} else {
		// pfew! just in time... c.leastRU is the pointer of death
		c.leastRU = e.younger // (ok and this)
	}
	// some pointer mumbo jumbo
	e.younger.older = e.older
	e.older = c.mostRU  // my elder is whoever used to be youngest
	c.mostRU = e        // I'm the newest one now!
	e.younger = nil     // nobody's younger than me
	e.older.younger = e // eeeeee
	return
}

// Consume an operation from the channel and process it. Returns false if the
// channel was closed and the main loop should stop.
//
// Implemented as a separate function to ensure all local variables go out of
// scope when this main loop iteration is complete.
//
// Imagine this function accepted an operation directly and the mainLoop
// function were implemented as follows:
//
//     for op := range opchan {
//         mainLoopBody(op)
//     }
//
// This blocks on the read from opchan, but it is not immediately clear if the
// operation from the last iteration (haha) is cleared / garbage collected while
// this read is blocking. Because the operation struct contains a reference to
// the Cache, if that doesn't happen the entire cache will not be garbage
// collected.
func mainLoopBody(opchan <-chan operation) bool {
	op, ok := <-opchan
	if !ok {
		return false
	}
	c := op.c // careful: don't keep this one around!
	switch req := op.req.(type) {
	case reqSet:
		directSet(c, req)
	case reqDelete:
		directDelete(c, req)
	case reqGet:
		directGet(c, req)
	case reqOnMissFunc:
		c.onMiss = OnMissHandler(req)
	case reqMaxSize:
		c.maxSize = int64(req)
		trimCache(c)
	case reqGetSize:
		req <- c.size
		close(req)
	default:
		panic("Illegal cache operation")
	}
	return true
}

// does not keep any reference to the cache so it can be garbage collected
func mainLoop(opchan <-chan operation) {
	for mainLoopBody(opchan) {
	}
}

func (c *Cache) Init(maxsize int64) {
	c.maxSize = maxsize
	c.opChan = make(chan operation)
	c.entries = map[string]*cacheEntry{}
	go mainLoop(c.opChan)
	return
}

// Store this item in cache. Panics if the cacheable is nil.
func (c *Cache) Set(id string, p Cacheable) {
	if p == nil {
		panic("Cacheable value must not be nil")
	}
	c.opChan <- operation{c, reqSet{payload: p, id: id}}
	return
}

var ErrNotFound = errors.New("Key not found in cache")

func (c *Cache) Get(id string) (Cacheable, error) {
	replychan := make(chan replyGet)
	req := reqGet{id: id, reply: replychan}
	c.opChan <- operation{c, req}
	reply := <-replychan
	return reply.val, reply.err
}

func (c *Cache) Delete(id string) {
	c.opChan <- operation{c, reqDelete(id)}
}

// Clean up the resources associated with this goroutine. Stops the main loop
// and allows the cache to be garbage collected once all user references are
// gone.
func (c *Cache) Close() error {
	finalizeCache(c)
	return nil
}

// Used to populate the cache if an entry is not found.  Say you're looking
// for entry "bob". But there is no such entry in your cache! Do you always
// handle that in the same way? Get "bob" from disk or S3? Then this function
// is for you! Make this your "persistent storage lookup" function, hook it up
// to your cache right here and it will be called automatically next time
// you're looking for bob. The advantage is that you can expect Get() calls to
// resolve.
//
// The Get() call invoking this OnMiss will always return whatever value is
// returned from the OnMiss handler, error or not.
//
// If the function returns a non-nil error, that error is directly returned
// from the Get() call that caused it to be invoked.  Otherwise, if the function
// return value is not nil, it is stored in cache.
//
// Call with f is nil to clear.
//
// Return (nil, nil) to indicate the specific key could not be found. It will
// be treated as a Get() to an unknown key without an OnMiss handler set.
//
// Called from a separate goroutine, which does not block other operations on
// the cache. This means that calling Get concurrently with the same key, before
// the first OnMiss call returns, will invoke another OnMiss call; the last one
// to return will have its value stored in the cache. To avoid this, wrap the
// OnMiss handler in a NoConcurrentDupes.
func (c *Cache) OnMiss(f OnMissHandler) {
	c.opChan <- operation{c, reqOnMissFunc(f)}
}

// Feel free to change this whenever. The units are not bytes but just whatever
// unit it is that your cache entries return from Size(). If (roughly) all
// cached items are going to be (roughly) the same size it makes sense to
// return 1 from Size() and set maxSize to the maximum number of elements you
// want to allow in cache. To remove the limit altogether set a maximum size of
// 0. No elements will be purged with reason CACHEFULL until the next call to
// MaxSize.
//
// (reading this back I have to admit, once again, that I was wrong. obviously,
// if you're gonna be returning 1 from Size() might as well not specify the
// method at all because 1 is the default size assumed for objects that don't
// have a Size() method. but it explains the idea nicely so I'll leave it in.)
func (c *Cache) MaxSize(i int64) {
	c.opChan <- operation{c, reqMaxSize(i)}
}

func (c *Cache) Size() int64 {
	reply := make(chan int64)
	c.opChan <- operation{c, reqGetSize(reply)}
	return <-reply
}

func finalizeCache(c *Cache) {
	close(c.opChan)
}

// Create and initialize a new cache, ready for use. This library is designed to
// allow garbage collecting of the cache without explicit .Close once the
// reference returned by this function is not held by anyone anymore. However,
// testing shows that this does not work consistently across different systems
// (yet).  So, for now, to be sure the cache is cleared up, remember to .Close
// after use.
func New(maxsize int64) *Cache {
	var mem Cache
	c := &mem
	c.Init(maxsize)
	// Go's SetFinalizer cannot be unit tested, so basically it's a joke.
	//runtime.SetFinalizer(c, finalizeCache)
	return c
}

// Shared cache for configuration-less use

var sharedCache = New(0)

// Get an element from the shared cache.
func Get(id string) (Cacheable, error) {
	return sharedCache.Get(id)
}

// Put an object in the shared cache (requires no configuration).
func Set(id string, c Cacheable) {
	sharedCache.Set(id, c)
	return
}

// Delete an item from the shared cache.
func Delete(id string) {
	sharedCache.Delete(id)
	return
}

// A shared cache is available immediately for all users of this library. By
// default, there is no size limit. Use this function to change that.
func MaxSize(size int64) {
	sharedCache.MaxSize(size)
}
