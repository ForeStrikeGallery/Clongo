// B-tree index implementation

/* 
 1. In-memory implementation of B-Tree
 2. A way to map B-tree keys with files on the disk, and to do this quickly 
 3. Maintain consistency with B-tree nodes and files on the disk  
*/

package main

const (
	Order = 4 
) 

type BTreeIndex struct {
	Root *BTreeNode
	Order int
}

type BTreeNode struct {
	Leaf bool 
	Keys []string
	Values [][]string
	Children []*BTreeNode 
}

func NewBTreeIndex() *BTreeIndex {
	return &BTreeIndex{
		Root: nil,
		Order: Order,
	}
}

func NewBTreeNode() *BTreeNode {
	return &BTreeNode {
		Leaf: false,
		Keys: []string{},
		Values: [][]string{},
		Children: []*BTreeNode{},
	}
}

func (n *BTreeNode) InsertNonFull(key string, value string) {
	i := len(n.Keys) - 1

	if n.Leaf {
		n.Keys = append(n.Keys, "")

		for i >= 0 && key < n.Keys[i] {
			n.Keys[i+1] = n.Keys[i]
			i--
		}

		if key == n.Keys[i] {
			n.Values[i] = append(n.Values[i], value)
			return 
		}
		
		n.Keys[i+1] = key 
		n.Values[i+1] = []string{value}
	} else {
		for i >= 0 && key < n.Keys[i] {
			i--
		}		

		if len(n.Children[i].Keys) == 2*Order-1 {
			n.splitChild(i)
			
			if key > n.Keys[i] {
				i++
			}

			n.Children[i].InsertNonFull(key, value)
		}	
	}
}

func (n *BTreeNode) splitChild(i int) {
	child := n.Children[i]
	newChild := NewBTreeNode()
	newChild.Leaf = child.Leaf 
	n.Children = append(n.Children, nil) 

	copy(n.Children[i+2:], n.Children[i+1:])
	n.Children[i+1] = newChild 

	n.Keys = append(n.Keys, "")
	copy(n.Keys[i+1:], n.Keys[i:])
	n.Keys[i] = child.Keys[Order-1]

	n.Values = append(n.Values, nil)
	copy(n.Values[i+1:], n.Values[i:])
	n.Values[i] = child.Values[Order - 1]

	newChild.Keys = make([]string, len(child.Keys[int(order):]))
    copy(newChild.Keys, child.Keys[Order:])
	child.Keys = child.Keys[:Order-1]

	newChild.Values = make([][]string, len(child.Values[int(order):]))
	copy(newChild.Values, child.Values[Order:])
	child.Values = child.Values[:Order-1]
}

func (n *BTreeNode) Search(key string) []string {	
	i := 0
	for i < len(n.Keys) && key > n.Keys[i] {
		i++
	}

	if i < len(n.Keys) && key == n.Keys[i] {
		return n.Values[i]
	}

	if n.Leaf {
		return []string{}
	}

	return n.Children[i].Search(key)	
}








