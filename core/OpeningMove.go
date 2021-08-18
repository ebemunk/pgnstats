package core

import "sync/atomic"

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

func RecordOpening(ptr *OpeningMove, san string) *OpeningMove {
	// atomic.AddUint32(&ptr.Count, 1)
	openingMove := ptr.Find(san)
	if openingMove != nil {
		atomic.AddUint32(&openingMove.Count, 1)
		ptr = openingMove
	} else {
		openingMove = &OpeningMove{
			1, san, make([]*OpeningMove, 0),
		}
		ptr.Children = append(ptr.Children, openingMove)
		ptr = ptr.Children[len(ptr.Children)-1]
	}

	return ptr
}
