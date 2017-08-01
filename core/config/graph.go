package config

// Sorting order
const (
	SortInputsFirst = iota + 1
	SortOutputsFirst
)

func reverseList(s []int) (r []int) {
	for _, i := range s {
		i := i
		defer func() { r = append(r, i) }()
	}
	return
}

type graph map[int][]int

func topSortDFS(g graph) (order, cyclic []int) {
	L := make([]int, len(g))
	i := len(L)
	temp := map[int]bool{}
	perm := map[int]bool{}
	var cycleFound bool
	var cycleStart int
	var visit func(int)
	visit = func(n int) {
		switch {
		case temp[n]:
			cycleFound = true
			cycleStart = n
			return
		case perm[n]:
			return
		}
		temp[n] = true
		for _, m := range g[n] {
			visit(m)
			if cycleFound {
				if cycleStart > 0 {
					cyclic = append(cyclic, n)
					if n == cycleStart {
						cycleStart = 0
					}
				}
				return
			}
		}
		delete(temp, n)
		perm[n] = true
		i--
		L[i] = n
	}
	for n := range g {
		if perm[n] {
			continue
		}
		visit(n)
		if cycleFound {
			return nil, cyclic
		}
	}
	return L, nil
}
