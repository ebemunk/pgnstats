package core

//OpeningMove is a node that represents a move
type OpeningMove struct {
	Count    uint32
	San      string
	Children []*OpeningMove
}

//Find returns the child OpeningMove that matches the san
func (om *OpeningMove) Find(san string) *OpeningMove {
	for _, m := range om.Children {
		if m.San == san {
			return m
		}
	}

	return nil
}

//Prune recursively prunes OpeningMove trees depending on threshold count
func (om *OpeningMove) Prune(threshold int) {
	var toRemove []int

	for i, m := range om.Children {
		if m.Count < uint32(threshold) {
			toRemove = append(toRemove, i)
		}
	}

	for n, i := range toRemove {
		nx := i - n
		om.Children, om.Children[len(om.Children)-1] = append(om.Children[:nx], om.Children[nx+1:]...), nil
	}

	for _, m := range om.Children {
		m.Prune(threshold)
	}
}
