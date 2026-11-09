package tree

import (
	"fmt"
	"io"
)

// Leaves keep full keys so compressed-path searches can verify exact equality.
type leaf struct{ key string }

func (l *leaf) insert(key string, depth int) node {
	if key == l.key {
		return l
	}
	shared := 0
	for depth+shared < len(key) && depth+shared < len(l.key) && key[depth+shared] == l.key[depth+shared] {
		shared++
	}
	b := &branch{prefix: key[depth : depth+shared], children: &node4{}}
	b.attach(l.key, depth+shared, l)
	b.attach(key, depth+shared, &leaf{key: key})
	return b
}

func (l *leaf) lookup(key string, depth int, trace io.Writer) bool {
	found := l.key == key
	if found {
		fmt.Fprintf(trace, "leaf %q: found\n", l.key)
	} else {
		fmt.Fprintf(trace, "leaf %q: not found\n", l.key)
	}
	return found
}

func (l *leaf) print(out io.Writer, indent string) { fmt.Fprintf(out, "%sleaf %q\n", indent, l.key) }
