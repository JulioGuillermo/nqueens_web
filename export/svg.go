package export

import (
	"fmt"
	"io"
	"strings"
)

const (
	CellSize = 15

	FmtCell          = `<rect style="fill:#0088ff" width="%d" height="%d" x="%d" y="%d"/>`
	FmtQueen         = `<rect style="fill:#ffff00;stroke:#ff9800;stroke-width:1" width="%d" height="%d" x="%d" y="%d" ry="%f"/>`
	FmtBoardLayer    = `<g transform="translate(5,20)">`
	FmtBoardLayerEnd = `</g>`

	FmtTitle  = `<text transform="translate(5,10)" style="font-size:8px;line-height:1.25;font-family:sans;white-space:pre;fill:#0088ff">NQueens_%dx%d => %s</text>`
	FmtColNum = `<text transform="translate(%d,18)" style="font-size:3px;line-height:1.25;font-family:sans;white-space:pre;fill:#0088ff">%d</text>`
	FmtSvg    = `<svg width="%dmm" height="%dmm" viewBox="0 0 %d %d" xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape" xmlns="http://www.w3.org/2000/svg"><g>`
	FmtSvgEnd = `</g></svg>`
)

func renderBoard(w io.StringWriter, size int) error {
	var err error
	defer fmt.Println()
	for i := 0; i < size; i++ {
		fmt.Printf("\033[LExporting board => %.2f%%", float32(i)*100.0/float32(size))
		for j := 0; j < size; j++ {
			if i%2 == j%2 {
				continue
			}
			_, err = w.WriteString(fmt.Sprintf(FmtCell, CellSize, CellSize, i*CellSize, j*CellSize))
			if err != nil {
				return err
			}
		}
	}
	fmt.Printf("\033[LExporting board => 100%%")
	return nil
}

func renderQueen(x, y int) string {
	qSize := CellSize - 2
	x = x*CellSize + 1
	y = y*CellSize + 1
	rad := float64(qSize) / 2.0
	return fmt.Sprintf(FmtQueen, qSize, qSize, x, y, rad)
}

func renderQueens(w io.StringWriter, genome []int) error {
	defer fmt.Println()
	var err error
	size := len(genome)
	for i, j := range genome {
		fmt.Printf("\033[LExporting board => %.2f%%", float32(i)*100.0/float32(size))
		_, err = w.WriteString(renderQueen(i, j))
		if err != nil {
			return err
		}
	}
	fmt.Printf("\033[LExporting queens => 100%%")
	return nil
}

func renderBoardLayer(w io.StringWriter, size int, genome []int) error {
	_, err := w.WriteString(FmtBoardLayer)
	if err != nil {
		return err
	}
	err = renderBoard(w, size)
	if err != nil {
		return err
	}
	err = renderQueens(w, genome)
	if err != nil {
		return err
	}
	_, err = w.WriteString(FmtBoardLayerEnd)
	return err
}

func renderTitle(w io.StringWriter, n int, title string) error {
	_, err := w.WriteString(fmt.Sprintf(FmtTitle, n, n, title))
	fmt.Println("Exporting title => 100%")
	return err
}

func renderCellNum(w io.StringWriter, gnome []int) error {
	var err error
	defer fmt.Println()
	size := len(gnome)
	for i, v := range gnome {
		fmt.Printf("\033[LRendering cells numbers => %.2f%%", float32(i)*100.0/float32(size))
		_, err = w.WriteString(fmt.Sprintf(FmtColNum, CellSize*i+5, v+1))
		if err != nil {
			return err
		}
	}
	fmt.Printf("\033[LRendering cells numbers => 100%%")
	return nil
}

func Export(n int, title string, genome []int) (string, error) {
	fmt.Println("Exporting svg...")

	var out strings.Builder

	imgSize := n*CellSize + 10
	_, err := out.WriteString(fmt.Sprintf(
		FmtSvg,
		imgSize,
		imgSize+15,
		imgSize,
		imgSize+15,
	))
	if err != nil {
		return "", err
	}

	err = renderBoardLayer(&out, n, genome)
	if err != nil {
		panic(err)
	}
	err = renderTitle(&out, n, title)
	if err != nil {
		panic(err)
	}

	err = renderCellNum(&out, genome)
	if err != nil {
		panic(err)
	}

	_, err = out.WriteString(FmtSvgEnd)
	if err != nil {
		panic(err)
	}
	fmt.Println("Exporting svg completed.")

	return out.String(), nil
}
