package tree

// Degree 2 keeps splits visible: at most 3 keys and 4 children per node.
const degree = 2
const maxKeys = 2*degree - 1

// A child (and the root) can be either a leaf or a branch.
type node interface {
	nodeKeys() []int
	insertNonFull(key int)
	split() (separator int, right node)
}
