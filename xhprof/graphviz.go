package xhprof

import (
	"errors"
	"fmt"
	"math"
)

func GenerateDotScript(m *PairCallMap, threshold float32, function string, criticalPath bool, right map[string]*Call, left map[string]*Call) (string, error) {
	result := "digraph call_graph {\n"

	maxWidth := float32(5)
	maxHeight := float32(3.5)
	maxFontSize := 35
	maxSizingRatio := float32(20)

	callMap := m.GetCallMap()
	main, ok := callMap["main()"]
	if !ok {
		return "", errors.New("Call map has no main()")
	}
	mainWt := main.WallTime

	var path map[string]bool
	var pathEdges map[string]bool

	if criticalPath {
		path, pathEdges = getCriticalPath(m)
	}

	if function != "" {
		relatedFuncs := getRelatedFuncs(m, function)
		for name, _ := range callMap {
			if _, ok := relatedFuncs[name]; !ok {
				delete(relatedFuncs, name)
			}
		}
	}

	curId := 0
	maxWt := float32(0)
	for name, c := range callMap {
		if function == "" && (c.WallTime/mainWt) < threshold {
			delete(callMap, name)
			continue
		}

		if maxWt == float32(0) || maxWt < c.ExclusiveWallTime {
			maxWt = c.ExclusiveWallTime
		}

		c.graphvizId = curId
		curId += 1
	}

	sizingFactor := float32(0)
	for name, c := range callMap {
		if c.ExclusiveWallTime == 0 {
			sizingFactor = maxSizingRatio
		} else {
			sizingFactor = float32(math.Min(float64(maxWt/c.ExclusiveWallTime), float64(maxSizingRatio)))
		}

		fillColor := ""
		if sizingFactor < 1.5 {
			fillColor = ", style=filled, fillcolor=red"
		}

		if _, ok := path[name]; criticalPath && fillColor == "" && ok {
			fillColor = ", style=filled, fillcolor=yellow"
		}

		fontSize := fmt.Sprintf(", fontsize=%d", int(float32(maxFontSize)/((sizingFactor-1)/10+1)))
		width := fmt.Sprintf(", width=%.1f", maxWidth/sizingFactor)
		height := fmt.Sprintf(", height=%.1f", maxHeight/sizingFactor)

		shape := "box"
		n := ""
		if name == "main()" {
			shape = "octagon"
			n = fmt.Sprintf("Total: %2.2f ms \\nmain()", mainWt/1000)
		} else {
			n = fmt.Sprintf("%s\\nInc: %.3f ms (%.1f%%)", name, c.WallTime/1000, 100*c.WallTime/mainWt)
		}

		var label string
		if left == nil {
			label = fmt.Sprintf(", label=\"%s\\nExcl: %.3f ms (%.1f%%)\\n%d total calls\"", n, c.ExclusiveWallTime/1000, 100*c.ExclusiveWallTime/mainWt, c.Count)
		} else {
			leftC, lOk := left[name]
			rightC, rOk := right[name]

			if lOk && rOk {
				label = fmt.Sprintf(
					", label=\"%s\\nInc: %.3f ms - %.3f ms = %.3f ms\\nExcl: %.3f ms - %.3f ms = %.3f ms\\nCalls: %d - %d = %d\"",
					name,
					leftC.WallTime/1000, rightC.WallTime/1000, c.WallTime/1000,
					leftC.ExclusiveWallTime/1000, rightC.ExclusiveWallTime/1000, c.ExclusiveWallTime/1000,
					leftC.Count, rightC.Count, c.Count,
				)
			} else if lOk {
				label = fmt.Sprintf(
					", label=\"%s\\nInc: %.3f ms - %.3f ms = %.3f ms\\nExcl: %.3f ms - %.3f ms = %.3f ms\\nCalls: %d - %d = %d\"",
					name,
					leftC.WallTime/1000, 0, c.WallTime/1000,
					leftC.ExclusiveWallTime/1000, 0, c.ExclusiveWallTime/1000,
					leftC.Count, 0, c.Count,
				)
			} else {
				label = fmt.Sprintf(
					", label=\"%s\\nInc: %.3f ms - %.3f ms = %.3f ms\\nExcl: %.3f ms - %.3f ms = %.3f ms\\nCalls: %d - %d = %d\"",
					name,
					0, rightC.WallTime/1000, c.WallTime/1000,
					0, rightC.ExclusiveWallTime/1000, c.ExclusiveWallTime/1000,
					0, rightC.Count, c.Count,
				)
			}
		}
		result += fmt.Sprintf("N%d[shape=%s %s%s%s%s%s];\n", c.graphvizId, shape, label, width, height, fontSize, fillColor)
	}

	for name, c := range m.M {
		parent, child := parsePairName(name)

		parentC, ok := callMap[parent]
		if !ok {
			continue
		}

		childC, ok := callMap[child]
		if !ok {
			continue
		}

		if function != "" && parent != function && child != function {
			continue
		}

		label := "1 call"
		if c.Count != 1 {
			label = fmt.Sprintf("%d calls", c.Count)
		}

		headLabel := "0.0%"
		if childC.WallTime > 0 {
			headLabel = fmt.Sprintf("%.1f%%", 100*c.WallTime/childC.WallTime)
		}

		tailLabel := "0.0%"
		if parentC.WallTime > 0 {
			tailLabel = fmt.Sprintf("%.1f%%", 100*c.WallTime/(parentC.WallTime-parentC.ExclusiveWallTime))
		}

		lineWidth := 1
		arrowSize := 1
		if _, ok := pathEdges[name]; criticalPath && ok {
			lineWidth = 10
			arrowSize = 2
		}

		result += fmt.Sprintf(
			"N%d -> N%d[arrowsize=%d, color=grey, style=\"setlinewidth(%d)\", label=\"%s\", headlabel=\"%s\", taillabel=\"%s\" ];\n",
			parentC.graphvizId, childC.graphvizId, arrowSize, lineWidth, label, headLabel, tailLabel,
		)
	}

	result += "\n}"

	return result, nil
}

func GenerateDiffDotScript(m1, m2 *PairCallMap, threshold float32) (string, error) {
	right := m1.GetCallMap()
	left := m2.GetCallMap()
	diff := m2.Subtract(m1)

	return GenerateDotScript(diff, threshold, "", true, right, left)
}

func getCriticalPath(m *PairCallMap) (map[string]bool, map[string]bool) {
	path := make(map[string]bool)
	pathEdges := make(map[string]bool)
	visited := make(map[string]bool)
	childrenMap := m.GetChildrenMap()
	node := "main()"

	for node != "" {
		visited[node] = true
		if children, ok := childrenMap[node]; ok {
			maxChild := ""
			for _, child := range children {
				if _, ok := visited[child]; ok {
					continue
				}

				if maxChild == "" || m.M[pairName(node, child)].WallTime > m.M[pairName(node, maxChild)].WallTime {
					maxChild = child
				}
			}

			if maxChild != "" {
				path[maxChild] = true
				pathEdges[pairName(node, maxChild)] = true
			}

			node = maxChild
		} else {
			node = ""
		}
	}

	return path, pathEdges
}

func getRelatedFuncs(m *PairCallMap, f string) map[string]bool {
	r := make(map[string]bool)

	for name, _ := range m.M {
		parent, child := parsePairName(name)
		if parent == f || child == f {
			r[parent] = true
			r[child] = true
		}
	}

	return r
}
