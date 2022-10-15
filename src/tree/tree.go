package tree

type node struct {
	value int
	left  *node
	right *node
}

type Tree struct {
	root *node
	size int
}

func (t *Tree) Add(val int) {
	newNode := &node{value: val}
	if t.root == nil {
		t.root = newNode
		return
	}
	traverseAdd(t.root, val)
	t.size++
}

func (t *Tree) InOrderTraversal() []int {
	if t.root == nil {
		return []int{}
	}
	return traverse(t.root)
}

func traverse(n *node) []int {
	if n == nil {
		return []int{}
	}
	arr := traverse(n.left)
	arr = append(arr, n.value)
	arr = append(arr, traverse(n.right)...)
	return arr
}

func (t *Tree) CalculateHash() int {
	hash := 1
	for _, val := range t.InOrderTraversal() {
		newVal := val + 2
		hash = (hash*newVal + newVal) % 1000
	}
	return hash
}

func traverseAdd(cur *node, val int) {
	if val == cur.value {
		panic("duplicate key detected")
	} else if val < cur.value {
		if cur.left == nil {
			cur.left = &node{value: val}
			return
		}
		traverseAdd(cur.left, val)
	} else {
		if cur.right == nil {
			cur.right = &node{value: val}
			return
		}
		traverseAdd(cur.right, val)
	}
}
