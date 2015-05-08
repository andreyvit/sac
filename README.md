**sac**
=======

Small Array Container - chained reusable associative arrays
-----------------------------------------------------------

This package defines a datastructure that may be used as a small recyclable key/value store.
As a key/value store, it deals with 5 operations:

* Get()		 - tries to retrieve a Value. Returns an error if key is absent.
* Put()		 - inserts a Value for a given Key.
* Delete()
* Clear()	 - Delete all elements
* Length()	 - returns the number of elements stored

It is safe for concurrent use.

The emphasis put on the small size is due to the fact that the operations have a linear,
if not quadratic algorithmic time complexity. 
It means that the more elements stored, the slower the operations will be. 
But it does not matter if you need to store few elements and it may even be faster.
About a dozen or so elements is the use-case designed for.

The rationale is that because we have control on this datastructure, we can delete its elements
when not needed and put it back into a pool instead of it being garbage collected.
For a webserver, it means that we have a reusable key/value store that can be shared
amongst middleware without letting it become garbage at the end of request handling.

*Trivia* : "sac" is the French word for "bag".

**TODO** 
Create a generator to allow people to specialize Key and Value types statically.
Note that type specialization is not enough if you want the datastructure to fit
nicely in the cache.
The cacheline padding will have to be rethought for instance.

API
---

###Instance creation
``` go

// Let's create a new pool of sac objects
Pool := sac.Pool()

// Now, let's create a new sac object.
storage := sac.New(Pool)
```

###Put
``` go

// Let's insert something
// Key and Value are of type interface{}

storage.Put("Ultron", "Singularity")
store.Put(42, "Smurf")
store.Put("Ben","Fed")

```

###Get
``` go
// Let's retrieve something
res,err:= storage.Get("Ultron") // returns interface{}("Singularity")

```

###Delete
``` go
// Let's delete something
storage.Delete("Ben")

```

###Clear
``` go
// Let's delete everything
storage.Clear()

```

###Length
``` go
// Let's check the number of items inside the sac
n := storage.Length() // returns 0 since all clear.

```
For completeness, please refer to the package [documentation].

[documentation]:https://godoc.org/github.com/atdiar/sac
