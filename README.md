A "random binary forest" is a hybrid between kd-trees and random forests.
For nearest neighbors this ends up being similar to [Minhash Forests](https://github.com/ekzhu/datasketch) and to
[Spotify's annoy library](https://github.com/spotify/annoy).

We build a forest of roughly-binary search trees, with each tree being
built as follows: pick a random subset of features at each split, look for
the "best" feature (see below), split on that feature, and then recurse.

We want the split to be close to the median for the best search speeds (as
this will give us trees that are almost binary), but we want to maximize
variance for accuracy-optimization (e.g. if we have two features
A = [5, 5, 5, 6, 6, 6] and B = [0, 0, 0, 10, 10, 10], then we want to choose
B so that noisy data is less likely to fall on the wrong side of the split).

These two goals can conflict, so right now we just use a simple split
function that splits closest to the median. This has the added advantage that
you don't need to normalize features to have similar distributions.

We have another split function that takes variance into account, but this is
currently unused.
