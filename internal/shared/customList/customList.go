package customList

type LinkedList struct {
	Root *Node
}

type Node struct {
	Next     *Node
	Data     *NodeData
	InAction bool
}

type NodeData struct {
	Value       float64
	Operation   rune
	IsOperation bool
}

func NewLinkedList() *LinkedList {
	return &LinkedList{
		nil,
	}
}

func (list *LinkedList) Add(data *NodeData) error {
	list.Root = &Node{
		Next: list.Root,
		Data: data,
	}
	return nil
}

func (list *LinkedList) Replace(newData *NodeData, cur *Node, newNext *Node) error {
	cur.Data = newData
	cur.Next = newNext
	return nil
}
