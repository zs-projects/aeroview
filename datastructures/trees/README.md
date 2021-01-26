# Compact tree representation 

In this subpackage, we explore an alternative representation for immutable trees. 
The core idea behind the alternative representation is to store the nodes of the tree 
in level order in an array and encode the structure of the tree somewhere else. 

The main characteristics of the data structure will be as follow : 

* Low memory overhead: for a tree with `N` nodes, the overhead is `O(N)`.

* Fast and easy Serialisation/Deserialisation: writes to and reads from disks will be efficient.

* Efficient lookups: 

  - Accessing children is O(1) 
  - Accessing parent is O(Log(N))



# Binary trees  

## Intuition:
Traditionally, we represent such a tree with a recursive data structure that would look the sample below: 

```go
type Node struct {
	value int
	left *Node
	right *Node
} 
```

Representing a tree this way makes it's serialisation and deserialisation difficult. First, we need to 
determine a traversal order for the nodes ( pre-order, post-order, level-order, ...). The first step
allows us to get all the nodes in predetermined order. However, it's is not enough to preserve the 
structure of the tree. Therefore, as a second step, we need to figure out a way to preserve the structure. 

Let's consider an example, we will build an intuition about how to preserve the structure of a tree 
using the example below. Later, we will present a generalization. 
 
```
					5
				   / \
				  /   \
				 /     \
				7       8
			       / \     / \
			      /   \   /   \
			     3     x 6    11
			    / \     / \   / \
		           x   9   7   x 15  x
			  / \           / \
			 x   4         x  17
			    / \
			   2   x
```
First, we start by writng down the level order traversal for this tree, it goes as follow : 

```
5, 7, 8, 3, 6, 11, 9, 7, 15, 4, 17, 2
```

Second, as discussed above, we need to figure out a way to preserve the structure of the tree. The key 
idea to recover the structure is to know for each node of the tree, whether it has a left and a right 
child. Once you know this, you can figure out, for each node in the level order above, the position its left and right childs in the level order list.


Let's build the data structure for the tree above. We will walk the tree in level order and encode the presence of a child with 1 and the absence with 0. This get us the following structure : 

```
Nodes     : 5    7    8    3    6    11   9    7    15   4    17   2
=====================================================================
Structure : 11   10   11   01   10   10   01   00   01   10   00   00
            |    |         |
            |    |         |
            |    |         |-> 3 has no left child and and a right child
            |    |         
            |    |-> 7 has a left child but no right child
            |             
            |-> 5 has a left and a right child.
```

**How to get left and right child positions ?**
Let us take node `3` as an example. We know from the structure above that `3` has only a right child. 
In addition, if we look at all the nodes that come before `3` in the level order traversal and look at their associated values, we can also count the number of children they all have. 


The nodes before `3` in the level order reprensation have 5 children. If we also take into account the 
root node. There are at least 6 positions before the right child of `3`. What do we find if we look at the
position 6 from the level order list? we find `9` which is the right child of `3`.

## Bitsets, Rank and Select

In the example above, we represented the structure of the binary tree as a succession of ones and zeros. Such a structure is very common, it's called a `bitset`. 

There are many usefull operations that we can do on a bitset. More specifically, we will look at two of them here:  

 * `Rank`: `rank(i)` counts the number of ones until the `i-th` position. We used it in the example above.
 * `Select`: `select(i)` return the position of the `i-th` one in the bitset.

From now, we will suppose that we can efficiently perform those operations on a bitset. In fact, it is quite common to find constant time implementations of `Rank` and logarithmic time implementations of `Select`. You can find an implementation in the [rank package](/rank/README.md).


## Binary tree structure using bitsets

Now, it's time to properly explain how to build and use a bitset to encode the structure of a binary tree `T`. 

First, we initialize an empty bitset `B`. Then, we compute level order traversal of the tree `T` and store
it in an array `L`. Finally, we traverse the the level order array `L` and for each node `N` that we
encounter: we update the bitset `B` as follow. We add `1` if the current node `N` has a left child and `0`
otherwise and we do the same for the right child.

Given the bitset `B` and a node `N` that is at position `i` in the level order array `L`, we can compute, left child, right child and parent as follow : 

* **Left child** : if `B[2*i] != 0 ` then position of left child is `Rank(2*i)` 

* **Right child** : if `B[2*i + 1] != 0 ` then position of left child is `Rank(2*i + 1)`

* **Parent** : parent node index is given by then integer part of `Select(i) / 2`

A formal proof of what stated above, is far beyond the scope of this repository. However, if you want to 
see one, you can take a look at this [amazing MIT lecture](https://www.youtube.com/watch?v=3Y2weLDiUWw).


## Space and time analysis

### Time Analysis
#### Building the tree: 
As you have seen above, we traverse the tree once in level order to build it which takes `O(n)`. The optimized version of rank and select will also need to traverse the bitset `S` once. which means that the overal complexity is still `O(n)`

#### Left Child and Right Child: 
The time complexity of left child and right child is the same as the complexity of `Rank`.We use an implemenation of `Rank` that is constant time.

#### Parent: 
As stated above, the time complexity of `Parent` is the same as `Select`. The `Select` implementation we use has `Logarithmic` complexity.

### Space analysis 
In order to store a tree `T` with `n` elements of size `s` bits we need at least: 

- `n` * `s` space for the level order array `L`
- `2` * `n` bits for the bitset `S`.

We could stop here. However, by doing so, we would have to settle for linear time complexity on left child,
right child and parent operations. More specifically, given that we are not storing any additional 
metadata about the structure of the bitset `S`. The only way we can compute `Rank` and `Select`is by doing a linear traversal. 

The implementation presented here, goes a bit further and achieves constant time complexity on the left and right child operations and logarithmic time on the parent operation. In order to do so, we use and optimised implementation of `Rank` and `Select` that will add `o(n)` space overhead.

# K-Ary trees

# Serialization and De-serialization

# How to reuse the provided implementation

# Benchmarks
