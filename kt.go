package kt

import (
	"bytes"
	"fmt"
	"sort"
)

type Node struct {
	Child, Sib, Parent, Fail *Node
	In                       byte
	Depth                    int
	Output                   []int
	Id                       int
}
type queue []*Node
type NodeAction func(*Node, ...interface{})
type Match struct {
	Position, Pattern int
}
type MatchSlice []Match

var nodeId = 1

func (v *Node) findChild(c byte) *Node {
	child := v.Child
	for child != nil && child.In != c {
		child = child.Sib
	}
	return child
}
func (q *queue) add(n *Node) {
	*q = append(*q, n)
}
func (q *queue) get() *Node {
	n := (*q)[0]
	*q = (*q)[1:]
	return n
}
func (n *Node) String() string {
	w := new(bytes.Buffer)
	writeTree(n, w)
	return (w.String())
}
func (root *Node) Search(t []byte, p []string) []Match {
	t = append(t, 0)
	matches := make([]Match, 0)
	var match Match
	v := root
	j := 0
	for j < len(t)-1 {
		for c := v.findChild(t[j]); c != nil; c = v.findChild(t[j]) {
			if len(c.Output) > 0 {
				for _, o := range c.Output {
					match.Position = j - len(p[o]) + 1
					match.Pattern = o
					matches = append(matches, match)
				}
			}
			v = c
			j++
		}
		if v.Parent == nil {
			j++
		} else {
			v = v.Fail
		}

	}
	sort.Sort(MatchSlice(matches))
	return matches
}
func (m MatchSlice) Len() int { return len(m) }
func (m MatchSlice) Less(i, j int) bool {
	return m[i].Position < m[j].Position
}
func (m MatchSlice) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func NewKeywordTree(patterns []string) *Node {
	root := new(Node)
	for i, pattern := range patterns {
		v := root
		var j int
		for j = 0; j < len(pattern); j++ {
			c := v.findChild(pattern[j])
			if c != nil {
				v = c
			} else {
				break
			}
		}
		for k := j; k < len(pattern); k++ {
			child := new(Node)
			child.Parent = v
			child.In = pattern[k]
			child.Depth = v.Depth + 1
			child.Id = nodeId
			nodeId++
			if v.Child == nil {
				v.Child = child
			} else {
				cp := v.Child
				for cp.Sib != nil {
					cp = cp.Sib
				}
				cp.Sib = child
			}
			v = child
		}
		v.Output = append(v.Output, i)
	}
	root.Fail = root
	v := root.Child
	for v != nil {
		v.Fail = root
		v = v.Sib
	}
	BreadthFirst(root, setFailureLink, root)
	BreadthFirst(root, addOutput)
	return root
}
func NodeCount() int {
	return nodeId
}
func BreadthFirst(v *Node, fn NodeAction, args ...interface{}) {
	q := new(queue)
	for v != nil {
		fn(v, args...)
		q.add(v)
		v = v.Sib
		for v != nil {
			fn(v, args...)
			q.add(v)
			v = v.Sib
		}
		for v == nil && len(*q) > 0 {
			v = q.get()
			v = v.Child
		}
	}
}
func setFailureLink(v *Node, args ...interface{}) {
	root := args[0].(*Node)
	if v.Depth > 1 {
		v.Fail = root
		w := v.Parent
		for w.Parent != nil {
			w = w.Fail
			c := w.findChild(v.In)
			if c != nil {
				v.Fail = c
				break
			}
		}
	}
}
func addOutput(v *Node, args ...interface{}) {
	if v.Parent == nil {
		return
	}
	for w := v.Fail; w.Parent != nil; w = w.Fail {
		v.Output = append(v.Output, w.Output...)
	}
}
func writeTree(v *Node, w *bytes.Buffer) {
	if v == nil {
		return
	}
	if v.Parent != nil && v.Parent.Child.Id != v.Id {
		fmt.Fprint(w, ",")
	}
	if v.Child != nil {
		fmt.Fprint(w, "(")
	}
	writeTree(v.Child, w)
	label(v, w)
	writeTree(v.Sib, w)
	if v.Parent != nil && v.Sib == nil {
		fmt.Fprint(w, ")")
	}
	if v.Parent == nil {
		fmt.Fprint(w, ";")
	}
}
func label(v *Node, w *bytes.Buffer) {
	fmt.Fprintf(w, "%d[", v.Id+1)
	if v.Parent != nil {
		fmt.Fprintf(w, "%c", v.In)
	}
	fmt.Fprintf(w, "->%d", v.Fail.Id+1)
	if len(v.Output) > 0 {
		fmt.Fprintf(w, "{%d", v.Output[0]+1)
		for i := 1; i < len(v.Output); i++ {
			fmt.Fprintf(w, ",%d", v.Output[i]+1)
		}
		fmt.Fprintf(w, "}")
	}
	fmt.Fprintf(w, "]")
}
