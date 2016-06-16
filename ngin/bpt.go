package ngin

import (
	"bytes"
	"fmt"
	"strings"
	"unsafe"
)

/*
	3984 size of...
	struct {
		numk int
		keys    [55][64]byte
		ptrs    [56]*[]byte
		parent  *int
		leaf  struct{}
	}
*/

const M = 32 // 56

/*
func NodeToPtr(n *node) unsafe.Pointer {
	return unsafe.Pointer(&n)
}

func RecordToPtr(r *record) unsafe.Pointer {
	return unsafe.Pointer(&r)
}

func PtrToNode(ptr unsafe.Pointer) *node {
	return (*node)(unsafe.Pointer(ptr))
}

func PtrToRecord(ptr unsafe.Pointer) *record {
	return (*record)(unsafe.Pointer(ptr))
}

const (
	node_t byte = iota
	leaf_t
	data_t
)

func create(typ byte) unsafe.Pointer {
	switch typ {
	case node_t, leaf_t:
		n := &node{}
		if typ == leaf_t {
			n.leaf = struct{}{}
		}
		return unsafe.Pointer(&n)
	case data_t:
		r := &record{}
		return unsafe.Pointer(&r)
	default:
		return nil
	}
}
*/

type key_t []byte

func compare(a, b key_t) int {
	return bytes.Compare(a, b)
}

type val_t []byte

type ptr_t unsafe.Pointer

func asNode(p ptr_t) *node {
	return (*node)(unsafe.Pointer(p))
}

func asRecord(p ptr_t) *record {
	return (*record)(unsafe.Pointer(p))
}

// node represents a tree's node
type node struct {
	numk   int
	keys   [M - 1]key_t
	ptrs   [M]ptr_t
	parent *node
	leaf   struct{}
	next   *node
}

func (n *node) hasKey(k key_t) int {
	for i := 0; i < n.numk; i++ {
		if compare(k, n.keys[i]) == 0 {
			return i
		}
	}
	return -1
}

// leaf node record
type record struct {
	key key_t
	val val_t
}

// tree represents the main b+tree
type tree struct {
	root *node
}

// Has returns a boolean indicating weather or not the
// provided key and associated record / value exists.
func (t *tree) Has(key key_t) bool {
	return t.Get(key) != nil
}

// Add inserts a new record using provided key.
// It only inserts if the key does not already exist.
func (t *tree) Add(key key_t, val val_t) {
	// create record ptr for given value
	ptr := &record{key, val}

	// if the tree is empty
	if t.root == nil {
		t.root = startNewtree(key, ptr)
		return
	}
	// tree already exists, lets see what we
	// get when we try to find the correct leaf
	leaf := findLeaf(t.root, key)
	// ensure the leaf does not contain the key
	if leaf.hasKey(key) > -1 {
		return
	}
	// tree already exists, and ready to insert into
	if leaf.numk < M-1 {
		insertIntoLeaf(leaf, ptr.key, ptr)
		return
	}
	// otherwise, insert, split, and balance... returning updated root
	t.root = insertIntoLeafAfterSplitting(t.root, leaf, ptr.key, ptr)
}

// Set is mainly used for re-indexing
// as it assumes the data to already
// be contained the tree/index. it will
// overwrite duplicate keys, as it does
// not check to see if the key exists...
func (t *tree) Set(key []byte, val []byte) {
	// if the tree is empty, start a new one
	if t.root == nil {
		t.root = startNewtree(key, &record{key, val})
		return
	}

	// tree already exists, lets see what we
	// get when we try to find the correct leaf
	leaf := findLeaf(t.root, key)
	// ensure the leaf does not contain the key
	if i := leaf.hasKey(key); i > -1 {
		asRecord(leaf.ptrs[i]).val = val
		return
	}

	// create record ptr for given value
	ptr := &record{key, val}
	// if the leaf has room, then insert key and record
	if leaf.numk < M-1 {
		insertIntoLeaf(leaf, ptr.key, ptr)
		return
	}
	// otherwise, insert, split, and balance... returning updated root
	t.root = insertIntoLeafAfterSplitting(t.root, leaf, ptr.key, ptr)
}

/*
 *	inserting internals
 */

// first insertion, start a new tree
func startNewtree(key key_t, ptr *record) *node {
	root := &node{leaf: struct{}{}}
	root.keys[0] = key
	root.ptrs[0] = ptr_t(&ptr)
	root.ptrs[M-1] = nil
	root.parent = nil
	root.numk++
	return root
}

// creates a new root for two sub-trees and inserts the key into the new root
func insertIntoNewRoot(left *node, key key_t, right *node) *node {
	root := &node{}
	root.keys[0] = key
	root.ptrs[0] = ptr_t(&left)
	root.ptrs[1] = ptr_t(&right)
	root.numk++
	root.parent = nil
	left.parent = root
	right.parent = root
	return root
}

// insert a new node (leaf or internal) into tree, return root of tree
func insertIntoParent(root, left *node, key key_t, right *node) *node {
	if left.parent == nil {
		return insertIntoNewRoot(left, key, right)
	}
	leftIndex := getLeftIndex(left.parent, left)
	if left.parent.numk < M-1 {
		return insertIntoNode(root, left.parent, leftIndex, key, right)
	}
	return insertIntoNodeAfterSplitting(root, left.parent, leftIndex, key, right)
}

// helper->insert_into_parent
// used to find index of the parent's ptr to the
// node to the left of the key to be inserted
func getLeftIndex(parent, left *node) int {
	var leftIndex int
	for leftIndex <= parent.numk && asNode(parent.ptrs[leftIndex]) != left {
		leftIndex++
	}
	return leftIndex
}

/*
 *	Inner node insert internals
 */

// insert a new key, ptr to a node
func insertIntoNode(root, n *node, leftIndex int, key key_t, right *node) *node {
	copy(n.ptrs[leftIndex+2:], n.ptrs[leftIndex+1:])
	copy(n.keys[leftIndex+1:], n.keys[leftIndex:])
	n.ptrs[leftIndex+1] = ptr_t(&right)
	n.keys[leftIndex] = key
	n.numk++
	return root
}

// insert a new key, ptr to a node causing node to split
func insertIntoNodeAfterSplitting(root, oldNode *node, leftIndex int, key key_t, right *node) *node {

	var i, j int
	var tmpKeys [M]key_t
	var tmpPtrs [M + 1]ptr_t

	for i, j = 0, 0; i < oldNode.numk+1; i, j = i+1, j+1 {
		if j == leftIndex+1 {
			j++
		}
		tmpPtrs[j] = oldNode.ptrs[i]
	}

	for i, j = 0, 0; i < oldNode.numk; i, j = i+1, j+1 {
		if j == leftIndex {
			j++
		}
		tmpKeys[j] = oldNode.keys[i]
	}

	tmpPtrs[leftIndex+1] = ptr_t(&right)
	tmpKeys[leftIndex] = key

	split := cut(M)

	newNode := &node{}
	oldNode.numk = 0

	for i = 0; i < split-1; i++ {
		oldNode.ptrs[i] = tmpPtrs[i]
		oldNode.keys[i] = tmpKeys[i]
		oldNode.numk++
	}

	oldNode.ptrs[i] = tmpPtrs[i]

	prime := tmpKeys[split-1]

	for i, j = i+1, 0; i < M; i, j = i+1, j+1 {
		newNode.ptrs[j] = tmpPtrs[i]
		newNode.keys[j] = tmpKeys[i]
		newNode.numk++
	}

	newNode.ptrs[j] = tmpPtrs[i]

	// free tmps...
	for i = 0; i < M; i++ {
		tmpKeys[i] = nil
		tmpPtrs[i] = nil
	}
	tmpPtrs[M] = nil

	newNode.parent = oldNode.parent

	for i = 0; i <= newNode.numk; i++ {
		n := asNode(newNode.ptrs[i])
		n.parent = newNode
	}
	return insertIntoParent(root, oldNode, prime, newNode)
}

/*
 *	Leaf node insert internals
 */

// inserts a new key and *record into a leaf, then returns leaf
func insertIntoLeaf(leaf *node, key key_t, ptr *record) {
	var i, insertionPoint int
	for insertionPoint < leaf.numk && compare(leaf.keys[insertionPoint], key) == -1 {
		insertionPoint++
	}
	for i = leaf.numk; i > insertionPoint; i-- {
		leaf.keys[i] = leaf.keys[i-1]
		leaf.ptrs[i] = leaf.ptrs[i-1]
	}
	leaf.keys[insertionPoint] = key
	leaf.ptrs[insertionPoint] = ptr_t(&ptr)
	leaf.numk++
}

// inserts a new key and *record into a leaf, so as
// to exceed the order, causing the leaf to be split
func insertIntoLeafAfterSplitting(root, leaf *node, key key_t, ptr *record) *node {
	// perform linear search to find index to insert new record
	var insertionIndex int
	for insertionIndex < M-1 && compare(leaf.keys[insertionIndex], key) == -1 {
		insertionIndex++
	}
	var tmpKeys [M]key_t
	var tmpPtrs [M]ptr_t
	var i, j int
	// copy leaf keys & ptrs to temp
	// reserve space at insertion index for new record
	for i, j = 0, 0; i < leaf.numk; i, j = i+1, j+1 {
		if j == insertionIndex {
			j++
		}
		tmpKeys[j] = leaf.keys[i]
		tmpPtrs[j] = leaf.ptrs[i]
	}
	tmpKeys[insertionIndex] = key
	tmpPtrs[insertionIndex] = ptr_t(&ptr)

	leaf.numk = 0
	// index where to split leaf
	split := cut(M - 1)
	// over write original leaf up to split point
	for i = 0; i < split; i++ {
		leaf.ptrs[i] = tmpPtrs[i]
		leaf.keys[i] = tmpKeys[i]
		leaf.numk++
	}
	// create new leaf
	newLeaf := &node{leaf: struct{}{}}

	// writing to new leaf from split point to end of giginal leaf pre split
	for i, j = split, 0; i < M; i, j = i+1, j+1 {
		newLeaf.ptrs[j] = tmpPtrs[i]
		newLeaf.keys[j] = tmpKeys[i]
		newLeaf.numk++
	}
	// freeing tmps...
	for i = 0; i < M; i++ {
		tmpPtrs[i] = nil
		tmpKeys[i] = nil
	}
	newLeaf.ptrs[M-1] = leaf.ptrs[M-1]
	leaf.ptrs[M-1] = ptr_t(&newLeaf)
	for i = leaf.numk; i < M-1; i++ {
		leaf.ptrs[i] = nil
	}
	for i = newLeaf.numk; i < M-1; i++ {
		newLeaf.ptrs[i] = nil
	}
	newLeaf.parent = leaf.parent
	newKey := newLeaf.keys[0]
	return insertIntoParent(root, leaf, newKey, newLeaf)
}

// Get returns the record for
// a given key if it exists
func (t *tree) Get(key []byte) []byte {
	n := findLeaf(t.root, key)
	if n == nil {
		return nil
	}
	var i int
	for i = 0; i < n.numk; i++ {
		if compare(n.keys[i], key) == 0 {
			break
		}
	}
	if i == n.numk {
		return nil
	}
	r := asRecord(n.ptrs[i])
	return r.val
}

func find(root *node, key []byte) *record {
	var n *node = findLeaf(root, key)
	if n == nil {
		return nil
	}
	var i int
	for i = 0; i < n.numk; i++ {
		if compare(n.keys[i], key) == 0 {
			break
		}
	}
	if i == n.numk {
		return nil
	}
	return asRecord(n.ptrs[i])
}

/*
 *	Get node internals
 */

func findLeaf(root *node, key []byte) *node {
	var c *node = root
	if c == nil {
		return c
	}
	var i int
	for c.leaf != struct{}{} {
		i = 0
		for i < c.numk {
			if compare(key, c.keys[i]) >= 0 {
				i++
			} else {
				break
			}
		}
		c = asNode(c.ptrs[i])
	}
	return c
}

// binary search utility
func search(n *node, key key_t) int {
	lo, hi := 0, n.numk
	for lo < hi {
		md := (lo + hi) >> 1
		if compare(key, n.keys[md]) >= 0 {
			lo = md + 1
		} else {
			hi = md - 1
		}
	}
	return lo
}

// breadth-first-search algorithm, kind of
func (t *tree) BFS() {
	if t.root == nil {
		return
	}
	c, h := t.root, 0
	for c.leaf != struct{}{} {
		c = asNode(c.ptrs[0])
		h++
	}
	fmt.Printf(`[`)
	for h >= 0 {
		for i := 0; i < M; i++ {
			if i == M-1 && c.ptrs[M-1] != nil {
				fmt.Printf(` -> `)
				c = asNode(c.ptrs[M-1])
				i = 0
				continue
			}
			fmt.Printf(`[%s]`, c.keys[i])
		}
		fmt.Println()
		h--
	}
	fmt.Printf(`]\n`)
}

// finds the first leaf in the tree (lexicographically)
func findFirstLeaf(root *node) *node {
	if root == nil {
		return root
	}
	c := root
	for c.leaf != struct{}{} {
		c = asNode(c.ptrs[0])
	}
	return c
}

// Del deletes a record by key
func (t *tree) Del(key key_t) {
	record := t.Get(key)
	leaf := findLeaf(t.root, key)
	if record != nil && leaf != nil {
		// remove from tree, and rebalance
		t.root = deleteEntry(t.root, leaf, key, record)
	}
}

/*
 * Delete internals
 */

// helper for delete methods... returns index of
// a nodes nearest sibling to the left if one exists
func getNeighborIndex(n *node) int {
	for i := 0; i <= n.parent.numk; i++ {
		if asNode(n.parent.ptrs[i]) == n {
			return i - 1
		}
	}
	panic("Search for nonexistent ptr to node in parent.")
}

func removeEntryFromNode(n *node, key key_t, ptr ptr_t) *node {
	var i, numPtrs int
	// remove key and shift over keys accordingly
	for compare(n.keys[i], key) != 0 {
		i++
	}
	for i++; i < n.numk; i++ {
		n.keys[i-1] = n.keys[i]
	}
	// remove ptr and shift other ptrs accordingly
	// first determine the number of ptrs
	if n.leaf == struct{}{} {
		numPtrs = n.numk
	} else {
		numPtrs = n.numk + 1
	}
	i = 0
	for n.ptrs[i] != ptr {
		i++
	}

	for i++; i < numPtrs; i++ {
		n.ptrs[i-1] = n.ptrs[i]
	}
	// one key has been removed
	n.numk--
	// set other ptrs to nil for tidiness; remember leaf
	// nodes use the last ptr to point to the next leaf
	if n.leaf == struct{}{} {
		for i := n.numk; i < M-1; i++ {
			n.ptrs[i] = nil
		}
	} else {
		for i := n.numk + 1; i < M; i++ {
			n.ptrs[i] = nil
		}
	}
	return n
}

// deletes an entry from the tree; removes record, key, and ptr from leaf and rebalances tree
func deleteEntry(root, n *node, key key_t, ptr ptr_t) *node {
	var primeIndex, capacity int
	var neighbor *node
	var prime []byte

	// remove key, ptr from node
	n = removeEntryFromNode(n, key, ptr)

	if n == root {
		return adjustRoot(root)
	}

	var minKeys int
	// case: delete from inner node
	if n.leaf == struct{}{} {
		minKeys = cut(M - 1)
	} else {
		minKeys = cut(M) - 1
	}
	// case: node stays at or above min order
	if n.numk >= minKeys {
		return root
	}

	// case: node is bellow min order; coalescence or redistribute
	neighborIndex := getNeighborIndex(n)
	if neighborIndex == -1 {
		primeIndex = 0
	} else {
		primeIndex = neighborIndex
	}
	prime = n.parent.keys[primeIndex]
	if neighborIndex == -1 {
		neighbor = asNode(n.parent.ptrs[1])
	} else {
		neighbor = asNode(n.parent.ptrs[neighborIndex])
	}
	if n.leaf == struct{}{} {
		capacity = M
	} else {
		capacity = M - 1
	}

	// coalescence
	if neighbor.numk+n.numk < capacity {
		return coalesceNodes(root, n, neighbor, neighborIndex, prime)
	}
	return redistributeNodes(root, n, neighbor, neighborIndex, primeIndex, prime)
}

func adjustRoot(root *node) *node {
	// if non-empty root key and ptr
	// have already been deleted, so
	// nothing to be done here
	if root.numk > 0 {
		return root
	}
	var newRoot *node
	// if root is empty and has a child
	// promote first (only) child as the
	// new root node. If it's a leaf then
	// the while tree is empty...
	if root.leaf != struct{}{} {
		newRoot = asNode(root.ptrs[0])
		newRoot.parent = nil
	} else {
		newRoot = nil
	}
	root = nil // free root
	return newRoot
}

// merge (underflow)
func coalesceNodes(root, n, neighbor *node, neighborIndex int, prime []byte) *node {
	var i, j, neighborInsertionIndex, nEnd int
	var tmp *node
	// swap neight with node if nod eis on the
	// extreme left and neighbor is to its right
	if neighborIndex == -1 {
		tmp = n
		n = neighbor
		neighbor = tmp
	}
	// starting index for merged pointers
	neighborInsertionIndex = neighbor.numk
	// case nonleaf node, append k_prime and the following ptr.
	// append all ptrs and keys for the neighbors.
	if n.leaf != struct{}{} {
		// append k_prime (key)
		neighbor.keys[neighborInsertionIndex] = prime
		neighbor.numk++
		nEnd = n.numk
		i = neighborInsertionIndex + 1
		j = 0
		for j < nEnd {
			neighbor.keys[i] = n.keys[j]
			neighbor.ptrs[i] = n.ptrs[j]
			neighbor.numk++
			n.numk--
			i++
			j++
		}
		neighbor.ptrs[i] = n.ptrs[j]
		for i = 0; i < neighbor.numk+1; i++ {
			tmp = asNode(neighbor.ptrs[i])
			tmp.parent = neighbor
		}
	} else {
		// in a leaf; append the keys and ptrs.
		i = neighborInsertionIndex
		j = 0
		for j < n.numk {
			neighbor.keys[i] = n.keys[j]
			neighbor.ptrs[i] = n.ptrs[j]
			neighbor.numk++
			i++
			j++
		}
		neighbor.ptrs[M-1] = n.ptrs[M-1]
	}
	root = deleteEntry(root, n.parent, prime, n)
	n = nil // free n
	return root
}

// merge / redistribute
func redistributeNodes(root, n, neighbor *node, neighborIndex, primeIndex int, prime []byte) *node {
	var i int
	var tmp *node
	// case: node n has a neighnor to the left
	if neighborIndex != -1 {
		if n.leaf != struct{}{} {
			n.ptrs[n.numk+1] = n.ptrs[n.numk]
		}
		for i = n.numk; i > 0; i-- {
			n.keys[i] = n.keys[i-1]
			n.ptrs[i] = n.ptrs[i-1]
		}
		if n.leaf != struct{}{} {
			n.ptrs[0] = neighbor.ptrs[neighbor.numk]
			tmp = asNode(n.ptrs[0])
			tmp.parent = n
			neighbor.ptrs[neighbor.numk] = nil
			n.keys[0] = prime
			n.parent.keys[primeIndex] = neighbor.keys[neighbor.numk-1]
		} else {
			n.ptrs[0] = neighbor.ptrs[neighbor.numk-1]
			neighbor.ptrs[neighbor.numk-1] = nil
			n.keys[0] = neighbor.keys[neighbor.numk-1]
			n.parent.keys[primeIndex] = n.keys[0]
		}
	} else {
		// case: n is left most child (n has no left neighbor)
		if n.leaf == struct{}{} {
			n.keys[n.numk] = neighbor.keys[0]
			n.ptrs[n.numk] = neighbor.ptrs[0]
			n.parent.keys[primeIndex] = neighbor.keys[1]
		} else {
			n.keys[n.numk] = prime
			n.ptrs[n.numk+1] = neighbor.ptrs[0]
			tmp = asNode(n.ptrs[n.numk+1])
			tmp.parent = n
			n.parent.keys[primeIndex] = neighbor.keys[0]
		}
		for i = 0; i < neighbor.numk-1; i++ {
			neighbor.keys[i] = neighbor.keys[i+1]
			neighbor.ptrs[i] = neighbor.ptrs[i+1]
		}
		if n.leaf != struct{}{} {
			neighbor.ptrs[i] = neighbor.ptrs[i+1]
		}
	}
	n.numk++
	neighbor.numk--
	return root
}

func destroytreeNodes(n *node) {
	if n == nil {
		return
	}
	if n.leaf == struct{}{} {
		for i := 0; i < n.numk; i++ {
			n.ptrs[i] = nil
		}
	} else {
		for i := 0; i < n.numk+1; i++ {
			destroytreeNodes(asNode(n.ptrs[i]))
		}
	}
	n = nil // free
}

// All returns all of the values in the tree (lexicographically)
func (t *tree) All() [][]byte {
	leaf := findFirstLeaf(t.root)
	if leaf == nil {
		return nil
	}
	var vals [][]byte
	for {
		for i := 0; i < leaf.numk; i++ {
			if leaf.ptrs[i] != nil {
				// get record from leaf
				rec := asRecord(leaf.ptrs[i])
				// get doc and append to docs
				vals = append(vals, rec.val)
			}
		}
		// we're at the end, no more leaves to iterate
		if leaf.ptrs[M-1] == nil {
			break
		}
		// increment/follow pointer to next leaf node
		leaf = asNode(leaf.ptrs[M-1])
	}
	return vals
}

// Count returns the number of records in the tree
func (t *tree) Count() int {
	if t.root == nil {
		return -1
	}
	c := t.root
	for c.leaf != struct{}{} {
		c = asNode(c.ptrs[0])
	}
	var size int
	for {
		size += c.numk
		if c.ptrs[M-1] != nil {
			c = asNode(c.ptrs[M-1])
		} else {
			break
		}
	}
	return size
}

// Close destroys all the nodes of the tree
func (t *tree) Close() {
	destroytreeNodes(t.root)
}

// cut will return the proper
// split point for a node
func cut(length int) int {
	if length%2 == 0 {
		return length / 2
	}
	return length/2 + 1
}

/*
 * Printing methods
 */

var queue *node = nil

func enQueue(n *node) {
	var c *node
	if queue == nil {
		queue = n
		queue.next = nil
	} else {
		c = queue
		for c.next != nil {
			c = c.next
		}
		c.next = n
		n.next = nil
	}
}

func deQueue() *node {
	var n *node = queue
	queue = queue.next
	n.next = nil
	return n
}

func pathToRoot(root, child *node) int {
	var length int
	var c *node = child
	for c != root {
		c = c.parent
		length++
	}
	return length
}

func (t *tree) String() string {
	var i, rank, newRank int
	if t.root == nil {
		return "[]"
	}
	queue = nil
	var tree string
	enQueue(t.root)
	tree = "[["
	for queue != nil {
		n := deQueue()
		if n.parent != nil && n == asNode(n.parent.ptrs[0]) {
			newRank = pathToRoot(t.root, n)
			if newRank != rank {
				rank = newRank
				f := strings.LastIndex(tree, ",")
				tree = tree[:f] + tree[f+1:]
				tree += "],["
			}
		}
		tree += "["
		var keys []string
		for i = 0; i < n.numk; i++ {
			keys = append(keys, fmt.Sprintf("%q", n.keys[i]))
		}
		tree += strings.Join(keys, ",")
		if n.leaf != struct{}{} {
			for i = 0; i <= n.numk; i++ {
				enQueue(asNode(n.ptrs[i]))
			}
		}
		tree += "],"
	}
	f := strings.LastIndex(tree, ",")
	tree = tree[:f] + tree[f+1:]
	tree += "]]"
	return tree
}

func Btoi(b []byte) int64 {
	return int64(b[7]) |
		int64(b[6])<<8 |
		int64(b[5])<<16 |
		int64(b[4])<<24 |
		int64(b[3])<<32 |
		int64(b[2])<<40 |
		int64(b[1])<<48 |
		int64(b[0])<<56
}

func Itob(i int64) []byte {
	return []byte{
		byte(i >> 56),
		byte(i >> 48),
		byte(i >> 40),
		byte(i >> 32),
		byte(i >> 24),
		byte(i >> 16),
		byte(i >> 8),
		byte(i),
	}
}
