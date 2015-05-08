**sac**
=======

Small Array Container - chained, reusable, associative arrays
-------------------------------------------------------------

This package defines a datastructure that may be used as a small recyclable key/value store.
As a key/value store, it deals with 5 operations:

* Get()   		tries to retrieve a Value. Returns an error if key is absent.  

* Put()   		inserts a Value for a given Key.  

* Delete()   

* Clear()   	Delete all elements  

* Length()   	returns the number of elements stored  

It is safe for concurrent use.

The operations have a linear, or quadratic time complexity. 
It means that the more elements stored, the slower the operations will be. 
About a dozen or so elements is the use-case designed for.
For performance data, look at the test file which includes synthetic benchmarks.

###What can this be used for ?
Typical usage is as a per http Request storage datastructure. 
After the request servicing, it will be cleared
so that it can be reused.  


*TODO*: Create a generator to allow people to specialize Key and Value types statically.  

*Trivia*: "sac" is the French word for "bag".

API
---

###Instantiation
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
storage.Put(42, "Smurf")
storage.Put("Ben","FedReserve")

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
