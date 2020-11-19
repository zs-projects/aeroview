# Memtable:

The `memtable` is a simple immutable associative data structure ( think HashMap ) that is designed to be
`GC` friendly and memory efficient. `Memtable` stores all the keys and their associated values in two arrays in a "flat" structure and uses interpollation search to find a `given` key and figure out the where it's corresponding value is. 

## Naive implementation:
A first attempt at achieving the design objectives stated above could take the following form.
```go 
type Memtable struct {
	keys         []uint64 // a 64 bit farmhash of the originalKeys
	offsets      []uint32 // offsets has always one more element than `keys`.
	values       []byte
}
```
The structure is indeed `GC` friendly. In fact, in the above design, we only use 3 `heap` allocations. Consequently, the garbage collector only needs to track 3 pointers. Finally, we are only paying an extra `4 bytes` for each key in the data structure. 

Now, given the structure above, how could we retrivie the `value` associated with a `key` ? 

A naive solution could work as follow:

```go 
func (memtable Memtable) get(key string) (value []byte, ok bool) {
    for index, key := range memtable.keys { 
        if key == candidateKey { 
            valueStartIdx := memtable.offsets[index]
            valueStopIdx := memtable.offsets[index + 1]
            return memtable.data[valueStartIdx:valueStopIdx], true
        }
    } 
    return nil, false
}
```

We can see that the naive solution above has an `O(n)` time complexity. We can do better! 

**First improvement**: Let's keep the structure above and make one change. Upon creation of the `Memtable` we can sort the keys and keep them sorted at all time. Then, we can use binary search to find the `key`. Using this first improvement we can improve the time complexity of the get operation and achieve `O(log(n))`

## Improved design:

By changing the way we store the keys, we can even improve the time complexity of the `Get` operation and achieve `O(log(log(n)))`. The algorithm we will use to achieve this improvement is called [interpolation search](https://www.geeksforgeeks.org/interpolation-search/). This algorithm can achieve `O(log(log(n)))`a time complexity if and only if: 

* The keys are sorted.
* The keys follow a uniform distribution.

Sorting the keys is easy as we seen in the previous step. However, since the keys are determined by the memtable user. we can't always garantee they will follow a uniform distribution. Most of the time, the keys won't be uniform. 

Howerver, using a good hashing algorithm, we can address this limitation. In fact, the hash of the keys will follow a uniform distribution. Therefore, we can store the hash of the keys and sort them and use that instead of the original keys. The data structure becomes as follow : 

```go 
type Memtable struct {
	originalKeys []string // a 64 bit farmhash of the originalKeys
	hashedKeys   []uint64 // a 64 bit farmhash of the originalKeys
	offsets      []uint32 // offsets has always one more element than `keys`.
	values       []byte
}
```

Now, when we receive a `key`, we hash it and we use [interpolation search](https://www.geeksforgeeks.org/interpolation-search/) on the `hashedKeys` field to retrieve the `index`. Once we have the index, we get the values back as shown before.

### Interpolation search : 
Now let's take a look at how interpolation search works. First we have to remember the assumptions : 

1. The keys are uniformly distributed 
2. The keys are sorted.

![visual explanation of interpollated search](./interpollated_search.gif)

**How could we estime the value of the middle key in the `keys` array ?**

One strategy could be to add both the first and the last key and devide them by two. This is a reasonnable strategy because of the distribution hypothesis. Here is an intuition on why this makes sens: In order to evenly distribute the keys, we have to make the differences betweens two pair of consecutive keys almost constant. 
TODO: add an illustration

**How could we narrow down the region where a key might be ?**
let's say we are looking for a specific `key`. We can mesure it's distance from the smallest key in set. Now, if we use the  intuition from the previous step. we could come up with a reasonable estimate for where `key` would be. Here is how : 

```go 
func (memtable Memtable) PositionEstimate(key uint64, low, high int) int { 
    // if we already know from somewhere else that the key is between `low` and `high` 
    // here is how we can produce and educated guess on where it could be.
	span := float64(len(m.keys))
	last := float64(m.keys[high])
	first := float64(m.keys[low])
    // Now that we have the slope, we have an estimate 
    // the tells us at what rate the keys are increasing
    slope := span / (last - first)
    nbSteps := slope * float64(key-m.keys[low])
    return lo + int(math.Min(float64(hi-lo-1), float64(nbSteps)))
}
```

Of course, the above estimate is noisy ( it's variance is high, because it is produced with one sample !). However, it still good enough to drive our search. Here is how : 

```go
func (memtable Memtable) SimplifiedInterpolatedSearch(key uint64) (position int, ok bool) {
    low := 0
    high := len(memtable.keys) -1
    for memtable.keys[low] < key { 
        candidateIndex := memtable.PositionEstimate(key, 0, len(memtable.keys) - 1)
        if m.keys[candidateIndex] < hashedKey {
            low = candidateIndex + 1
        } else {
            high = candidateIndex
        }
    }
    if m.keys[low] != key {
		return nil, false
    }
    return low, true
}
```

The trick here, is to produce a noisy estimate and then use it to narrow down the search region and then repreat the process until the key is found.

#### Optimisations: 
There are two optimisations that we make in the practical implementation : 

1. When `high` - `low` gets smaller than 64, we revert to a linear search algorithm.
2. We only compute the `slope` once at the start of the algorithm. We have observed, that even though it could add some iterations to the algorthms, it made each iteration way faster and resulted in a faster implementation.

### Writing and Reading from disk: 
TODO: write about how to serialize/deserialise the tructure

### Limitations: 

The are two main limitations that come with the design of the memtable. First, given that the `Get` operation has to do more work then with a traditionnal `map`. It will be up to an order of magnitude slower. Second, the structure is immutable and adding one element to it requires that we rebuild the whole structure.