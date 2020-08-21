# RBF

A random binary forest is a fast method to find nearest neighbors.
It's a hybrid of a [k-d tree](https://en.wikipedia.org/wiki/K-d_tree)
and a [random forest](https://en.wikipedia.org/wiki/Random_forest).
It's very similar to [Spotify's annoy library](https://github.com/spotify/annoy)
and slightly similar to [Minhash Forests](https://github.com/ekzhu/datasketch).


## Download
```bash
$ go get github.com/moygit/rbf
```


## Example
```go
package main

import (
	"fmt"
	"github.com/moygit/rbf"
)

func main() {
	// Which point is (1,1) closest to, (0,0) or (10,10)?
	points := [][]byte{{0, 0}, {10, 10}}
	queryPoint := []byte{1, 1}

	// Build the forest and query it:
	var numTrees, depth, leafSize, numFeaturesToCompare int32 = 1, 2, 1, 1
	forest := rbf.TrainForest(points, numTrees, depth, leafSize, numFeaturesToCompare)
	count, results := forest.FindPointAllResults(queryPoint)
	nearest := points[results[numTrees-1][count-1]]

	fmt.Printf("Number of results: %d\n", count)                 // 1 point returned
	fmt.Printf("Nearest point to %v: %v\n", queryPoint, nearest) // Nearest point to (1,1) is (0,0)
}
```


## How it works

We build a forest of roughly-binary search trees, with each tree being
built as follows: pick a random subset of features at each split, look for
the "best" feature (see below), split on that feature, and then recurse.

We want the split to be close to the median for the best search speeds (as
this will give us trees that are almost binary), but we want to maximize
variance for accuracy-optimization (e.g. if we have two features
A = [5, 5, 5, 6, 6, 6]^T and B = [0, 0, 0, 10, 10, 10]^T, then we want to choose
B so that noisy data is less likely to fall on the wrong side of the split).

These two goals can conflict, so right now we just use a simple split
function that splits closest to the median. We have another split function that
takes variance into account, but this is currently unused.

This median-splitting has some pros and cons, outlined below.


## Median-splitting pros and cons

RBFs' current simple-split uses the median, not a metric. The more complex split,
currently unused, factors variance into the split for the reasons described above,
but it still doesn't use a metric.

This has the advantage that you don't need to normalize features and ensure they have
similar distributions and scales, but it has two disadvantages: first, as mentioned
above, not using variance makes the split more sensitive to noise in the data. TL;DR
of the second issue: not using a metric means we're being slightly loose with the
notion of ``nearest``.

Details of the latter issue: suppose you have a geographic distance feature and a
time-of-day feature. They have very different scales, so ``distance`` means something
very different with respect to these two features. You can scale them both down to
[0, 1], but now suppose you also have a person's height or blood-pressure as a
feature, or the price of an item. You can scale these down too, but they have very
different distributions than the first two. Nearest neighbors algorithms are very
sensitive to these differences, so now maybe you want histogram-normalization---and
the simple solution is starting to become messy. RBFs find *nearest-quantile*
neighbors, so they don't care about scale and distribution. To be clear: the
distribution could actually matter, so this isn't a free lunch, it's just a quick
lunch.


## Note on data science usage

(Summary: incomplete Python-callable version [here](https://github.com/moygit/c_rbf).)

I wrote this in Go expecting that it would be callable from Python (this was my first
Go project), but it turns out that [gopy](https://github.com/go-python/gopy) is
somewhat limited because the two languages' garbage-collectors collide. I have a
pre-alpha [C port](https://github.com/moygit/c_rbf) which *is* callable from Python.
But I ended up running this in a microservice, so I didn't need the Python interop
myself and I never finished the C port.
