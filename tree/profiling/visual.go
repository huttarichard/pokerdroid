package profiling

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/table"
	"github.com/pokerdroid/poker/tree"
)

const (
	cellSize = 80
	padding  = 80
)

var (
	actionColors = map[table.DiscreteAction]color.Color{
		table.DFold:  color.RGBA{R: 100, G: 150, B: 200, A: 255}, // Semi light blue for Fold
		table.DCall:  color.RGBA{R: 100, G: 200, B: 100, A: 255}, // Green for Call
		table.DCheck: color.RGBA{R: 100, G: 200, B: 100, A: 255}, // Green for Check
		table.DAllIn: color.RGBA{R: 150, G: 50, B: 80, A: 255},   // Dark red for AllIn
	}
)

type SaveCallback func(node tree.Node, dc *gg.Context) error

func SaveAllDir(dir string) SaveCallback {
	os.MkdirAll(dir, 0755)

	return func(node tree.Node, dc *gg.Context) error {
		// Ensure directory exists, creating it if necessary
		filename := filepath.Join(dir, fmt.Sprintf("%s.png", tree.GetPath(node)))
		return dc.SavePNG(filename)
	}
}
func SaveVisualization(profile *Profile, save SaveCallback) error {
	matrixSize := 13
	imgSize := cellSize*matrixSize + padding*2

	dist := profile.Nodes.Dist()
	actions := dist.Actions // All discrete actions used

	// Helper to get canonical label for a given matrix position
	formatCoordLabel := func(i, j int) string {
		// Get all combos for this coordinate
		cds := card.CardsInCoordsWithBlockersAt(i, j, nil)
		if len(cds) == 0 {
			return ""
		}

		// Get first combo as representative
		comb := cds[0]
		if len(comb) < 2 {
			return ""
		}

		r1, s1 := comb[0].Rank(), comb[0].Suite()
		r2, s2 := comb[1].Rank(), comb[1].Suite()

		// Ensure higher rank is listed first for readability
		if r2 > r1 {
			r1, r2 = r2, r1
			s1, s2 = s2, s1
		}

		isPair := r1 == r2
		isSuited := s1 == s2

		// Format label based on pair/suited/offsuit
		switch {
		case isPair:
			// e.g. "TT", "KK"
			return fmt.Sprintf("%s%s", r1.String(), r2.String())
		case isSuited:
			// e.g. "AJs", "87s"
			return fmt.Sprintf("%s%ss", r1.String(), r2.String())
		default:
			// e.g. "KTo", "52o"
			return fmt.Sprintf("%s%so", r1.String(), r2.String())
		}
	}

	// Helper to draw bars for a single cell.
	drawCellBars := func(dc *gg.Context, cellX, cellY int, distribution []float64) {
		// Convert (cellX, cellY) to canvas coordinates
		x := float64(padding + cellX*cellSize)
		y := float64(padding + cellY*cellSize)
		w := float64(cellSize)
		h := float64(cellSize)

		// Fill background for the cell
		dc.SetColor(color.RGBA{240, 240, 240, 255})
		dc.DrawRectangle(x, y, w, h)
		dc.Fill()

		// Sum distribution to normalize
		var sum float64
		for _, prob := range distribution {
			sum += prob
		}

		if sum > 0 {
			// Draw each action's portion of the bar
			startX := x
			for i, prob := range distribution {
				if prob <= 0 {
					continue
				}
				frac := float64(prob / sum)
				barW := frac * w

				act := actions[i]
				col, ok := actionColors[act]
				if !ok {
					// fallback color for raises or unknown actions
					col = color.RGBA{160 + uint8(i*10), 30, 30, 255}
				}
				dc.SetColor(col)
				dc.DrawRectangle(startX, y, barW, h) // Fill the entire cell height
				dc.Fill()

				startX += barW
			}
		}

		// Draw combo label with semi-transparent background for better readability
		label := formatCoordLabel(cellX, cellY)
		if label != "" {
			labelWidth, _ := dc.MeasureString(label)
			// Create a small semi-transparent white background just for the text
			dc.SetRGBA255(255, 255, 255, 0) // Semi-transparent white
			dc.DrawRectangle(x+4, y+4, labelWidth+8, 20)
			dc.Fill()

			// Draw the text
			dc.SetColor(color.Black)
			dc.DrawStringAnchored(label, x+8, y+18, 0, 0) // Adjusted position
		}

		// Draw a black frame around the cell - added after everything else
		dc.SetLineWidth(1.0)
		dc.SetColor(color.Black)
		dc.DrawRectangle(x, y, w, h)
		dc.Stroke()
	}

	// Draw legend at the bottom of the image
	drawLegend := func(dc *gg.Context) {
		legendY := float64(imgSize - 30)
		legendX := float64(padding)
		boxSize := 24.0
		gap := 4.0

		dc.SetColor(color.Black)
		dc.DrawStringAnchored("Legend:", legendX, legendY-20, 0, 1)
		legendX += 100

		// Draw a box and label for each action
		for i, act := range actions {
			col, ok := actionColors[act]
			if !ok {
				col = color.RGBA{160 + uint8(i*5), 30, 30, 255}
			}

			// Draw colored square
			dc.SetColor(col)
			dc.DrawRectangle(legendX, legendY-boxSize, boxSize, boxSize)
			dc.Fill()

			// Add label
			dc.SetColor(color.Black)
			actionStr := act.String()
			// Format for readability
			if act > 0 {
				actionStr = fmt.Sprintf("Raise %.1fx", float64(act))
			}
			dc.DrawStringAnchored(actionStr, legendX+boxSize+5, legendY-boxSize/2, 0, 0.5)

			legendX += boxSize + gap + 40 // move to next legend item
		}
	}

	// For each node in the tree, we have a 13×13 Matrix of slice-of-float distribution
	for nodeKey, mat := range dist.Matrix {
		dc := gg.NewContext(imgSize, imgSize)
		dc.SetColor(color.White)
		dc.Clear()

		// Draw outer border around matrix
		dc.SetLineWidth(2.0)
		dc.SetColor(color.Black)
		dc.DrawRectangle(float64(padding-1), float64(padding-1),
			float64(cellSize*matrixSize+2), float64(cellSize*matrixSize+2))
		dc.Stroke()

		// Render each cell
		for i := 0; i < matrixSize; i++ {
			for j := 0; j < matrixSize; j++ {
				distribution := mat[i][j] // e.g. [foldProb, callProb, raiseProb...]
				// Even if distribution is empty, still draw cell
				drawCellBars(dc, j, i, distribution)
			}
		}

		// Draw legend
		drawLegend(dc)

		// Add node path as title
		dc.SetColor(color.Black)
		path := tree.GetPath(nodeKey).String()
		dc.DrawStringAnchored("Node: "+path, float64(imgSize/2), float64(padding/2), 0.5, 0.5)

		// Save via callback
		if err := save(nodeKey, dc); err != nil {
			return err
		}
	}

	return nil
}
