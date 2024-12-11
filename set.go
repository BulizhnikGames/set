package set

import (
	"iter"
	"math"
)

type float interface {
	~float32 | ~float64
}

type signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type unsigned interface {
	~uint | ~uint8 | ~int16 | ~int32 | ~int64 | ~uintptr
}

type Ordered interface {
	float | signed | unsigned | ~string
}

type color bool

const (
	red   color = true
	black color = false
)

type node[T Ordered] struct {
	value  T
	color  color
	left   *node[T]
	right  *node[T]
	parent *node[T]
}

type Set[T Ordered] struct {
	root *node[T]
	size int
}

// New returns an empty set.
func New[T Ordered]() *Set[T] {
	return &Set[T]{}
}

// Len returns the number of elements in the set.
func (s *Set[T]) Len() int {
	return s.size
}

// Contains returns true if x is in the set.
func (s *Set[T]) Contains(x T) bool {
	n := s.root
	for n != nil {
		if x < n.value {
			n = n.left
		} else if x > n.value {
			n = n.right
		} else {
			return true
		}
	}
	return false
}

// Insert inserts the value x into the set if not present.
func (s *Set[T]) Insert(x T) {
	// Standard BST insert
	var parent *node[T]
	n := s.root
	for n != nil {
		parent = n
		if x < n.value {
			n = n.left
		} else if x > n.value {
			n = n.right
		} else {
			// Element already in the set
			return
		}
	}

	newNode := &node[T]{value: x, color: red, parent: parent}
	if parent == nil {
		// Tree was empty
		s.root = newNode
	} else if x < parent.value {
		parent.left = newNode
	} else {
		parent.right = newNode
	}

	s.insertFixup(newNode)
	s.size++
}

// insertFixup restores Red-Black Tree properties after insertion.
func (s *Set[T]) insertFixup(z *node[T]) {
	for z.parent != nil && z.parent.color == red {
		if z.parent == z.parent.parent.left {
			y := z.parent.parent.right
			if y != nil && y.color == red {
				z.parent.color = black
				y.color = black
				z.parent.parent.color = red
				z = z.parent.parent
			} else {
				if z == z.parent.right {
					z = z.parent
					s.leftRotate(z)
				}
				z.parent.color = black
				z.parent.parent.color = red
				s.rightRotate(z.parent.parent)
			}
		} else {
			y := z.parent.parent.left
			if y != nil && y.color == red {
				z.parent.color = black
				y.color = black
				z.parent.parent.color = red
				z = z.parent.parent
			} else {
				if z == z.parent.left {
					z = z.parent
					s.rightRotate(z)
				}
				z.parent.color = black
				z.parent.parent.color = red
				s.leftRotate(z.parent.parent)
			}
		}
	}
	s.root.color = black
}

func (s *Set[T]) leftRotate(x *node[T]) {
	y := x.right
	x.right = y.left
	if y.left != nil {
		y.left.parent = x
	}
	y.parent = x.parent
	if x.parent == nil {
		s.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

func (s *Set[T]) rightRotate(x *node[T]) {
	y := x.left
	x.left = y.right
	if y.right != nil {
		y.right.parent = x
	}
	y.parent = x.parent
	if x.parent == nil {
		s.root = y
	} else if x == x.parent.right {
		x.parent.right = y
	} else {
		x.parent.left = y
	}
	y.right = x
	x.parent = y
}

// Remove deletes x from the set if it exists.
func (s *Set[T]) Remove(x T) {
	z := s.root
	for z != nil {
		if x < z.value {
			z = z.left
		} else if x > z.value {
			z = z.right
		} else {
			s.deleteNode(z)
			return
		}
	}
}

// deleteNode removes a given node from the red-black tree.
func (s *Set[T]) deleteNode(z *node[T]) {
	var x, y *node[T]
	y = z
	originalColor := y.color
	if z.left == nil {
		x = z.right
		s.rbTransplant(z, z.right)
	} else if z.right == nil {
		x = z.left
		s.rbTransplant(z, z.left)
	} else {
		// Find successor
		y = s.minNode(z.right)
		originalColor = y.color
		x = y.right
		if y.parent == z {
			if x != nil {
				x.parent = y
			}
		} else {
			s.rbTransplant(y, y.right)
			y.right = z.right
			y.right.parent = y
		}
		s.rbTransplant(z, y)
		y.left = z.left
		y.left.parent = y
		y.color = z.color
	}
	s.size--
	if originalColor == black && x != nil {
		s.deleteFixup(x)
	}
}

// rbTransplant replaces subtree u with subtree v in the tree.
func (s *Set[T]) rbTransplant(u, v *node[T]) {
	if u.parent == nil {
		s.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}
	if v != nil {
		v.parent = u.parent
	}
}

// deleteFixup restores Red-Black properties after deletion.
func (s *Set[T]) deleteFixup(x *node[T]) {
	for x != s.root && x.color == black {
		if x == x.parent.left {
			w := x.parent.right
			if w.color == red {
				w.color = black
				x.parent.color = red
				s.leftRotate(x.parent)
				w = x.parent.right
			}
			if (w.left == nil || w.left.color == black) && (w.right == nil || w.right.color == black) {
				w.color = red
				x = x.parent
			} else {
				if w.right == nil || w.right.color == black {
					if w.left != nil {
						w.left.color = black
					}
					w.color = red
					s.rightRotate(w)
					w = x.parent.right
				}
				w.color = x.parent.color
				x.parent.color = black
				if w.right != nil {
					w.right.color = black
				}
				s.leftRotate(x.parent)
				x = s.root
			}
		} else {
			w := x.parent.left
			if w.color == red {
				w.color = black
				x.parent.color = red
				s.rightRotate(x.parent)
				w = x.parent.left
			}
			if (w.right == nil || w.right.color == black) && (w.left == nil || w.left.color == black) {
				w.color = red
				x = x.parent
			} else {
				if w.left == nil || w.left.color == black {
					if w.right != nil {
						w.right.color = black
					}
					w.color = red
					s.leftRotate(w)
					w = x.parent.left
				}
				w.color = x.parent.color
				x.parent.color = black
				if w.left != nil {
					w.left.color = black
				}
				s.rightRotate(x.parent)
				x = s.root
			}
		}
	}
	x.color = black
}

// Min returns the smallest element in the set.
func (s *Set[T]) Min() (T, bool) {
	if s.root == nil {
		var zero T
		return zero, false
	}
	m := s.minNode(s.root)
	return m.value, true
}

// Max returns the largest element in the set.
func (s *Set[T]) Max() (T, bool) {
	if s.root == nil {
		var zero T
		return zero, false
	}
	m := s.maxNode(s.root)
	return m.value, true
}

func (s *Set[T]) minNode(n *node[T]) *node[T] {
	for n.left != nil {
		n = n.left
	}
	return n
}

func (s *Set[T]) maxNode(n *node[T]) *node[T] {
	for n.right != nil {
		n = n.right
	}
	return n
}

func (s *Set[T]) All() iter.Seq[T] {
	return func(yield func(v T) bool) {
		maxSize := int(math.Floor(math.Log2(float64(s.size+2)/float64(5))) + 2)
		stack := make([]*node[T], 0, maxSize)
		n := s.root
		for len(stack) > 0 || n != nil {
			if n != nil {
				stack = append(stack, n)
				n = n.left
			} else {
				n = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if !yield(n.value) {
					return
				}
				n = n.right
			}
		}
	}
}
