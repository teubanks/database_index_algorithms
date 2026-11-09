package tree

// Small capacities make leaf and branch splits visible.
const maxKeys = 3

type node interface {
	nodeKeys() []int
	// insert keeps the left half in place. A non-nil right signals a split;
	// separator is then inserted by the parent.
	insert(key int) (separator int, right node)
}
