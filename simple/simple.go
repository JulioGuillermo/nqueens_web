package simple

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/julioguillermo/nqueens_web/export"
	"github.com/julioguillermo/nqueens_web/tool"
)

type SimpleNQ struct {
	board []int
	size  int
	start time.Time
}

func NewSimpleNQ(size int) *SimpleNQ {
	board := make([]int, size)
	return &SimpleNQ{
		size:  size,
		board: board,
		start: time.Now(),
	}
}

func (s *SimpleNQ) Check(idx, f int) bool {
	for i := 0; i < idx; i++ {
		if s.board[i] == f {
			return false
		}
		d := s.board[i] - f
		if d < 0 {
			d = -d
		}
		if d == idx-i {
			return false
		}
	}
	return true
}

func (s *SimpleNQ) GetSolutionsFor(x, y, c int) int {
	taken := make([]bool, s.size)

	// For the board
	for i := 0; i < x; i++ {
		taken[s.board[i]] = true
		d := c - i
		dp := s.board[i] + d
		dm := s.board[i] - d
		if dm > 0 {
			taken[dm] = true
		}
		if dp < s.size {
			taken[dp] = true
		}
	}
	// For the new queen
	taken[y] = true
	d := c - x
	dp := y + d
	dm := y - d
	if dm > 0 {
		taken[dm] = true
	}
	if dp < s.size {
		taken[dp] = true
	}
	sol := 0
	for _, t := range taken {
		if !t {
			sol++
		}
	}
	return sol
}

func (s *SimpleNQ) GetSolutions(x, y int) int {
	sol := 0
	for i := x + 1; i < s.size; i++ {
		sol += s.GetSolutionsFor(x, y, i)
	}
	return sol
}

func (s *SimpleNQ) GetFSolutions(x, y int) int {
	sol := 0
	for i := x + 1; i < s.size && i < x+10; i++ {
		sol += s.GetSolutionsFor(x, y, i)
	}
	return sol
}

func (s *SimpleNQ) GetSolsForPos(c int, pos []int) []int {
	sols := make([]int, len(pos))
	for i, p := range pos {
		sols[i] = s.GetSolutions(c, p)
	}
	return sols
}

func (s *SimpleNQ) GetFSolsForPos(c int, pos []int) []int {
	sols := make([]int, len(pos))
	for i, p := range pos {
		sols[i] = s.GetFSolutions(c, p)
		// sols[i] = s.GetSolutionsFor(c, p, c+1)
	}
	return sols
}

func (s *SimpleNQ) GetValids(idx int) []int {
	taken := make([]bool, s.size)

	for i := 0; i < idx; i++ {
		taken[s.board[i]] = true
		d := idx - i
		dp := s.board[i] + d
		dm := s.board[i] - d
		if dm > 0 {
			taken[dm] = true
		}
		if dp < s.size {
			taken[dp] = true
		}
	}

	valid := make([]int, 0, s.size)
	for i, t := range taken {
		if !t {
			valid = append(valid, i)
		}
	}
	return valid
}

func (s *SimpleNQ) Run() {
	moves := make([][]int, s.size)
	valid := make([]int, s.size)
	for i := range valid {
		valid[i] = i
	}
	moves[0] = valid

	sel := 0
	for i := 0; i < s.size && i >= 0; i++ {
		if len(moves[i]) > 0 {
			s.board[i] = moves[i][sel]
			moves[i] = append(moves[i][:sel], moves[i][sel+1:]...)
			if i < s.size-1 {
				valid = s.GetValids(i + 1)
				moves[i+1] = valid
			}
		} else {
			i -= 2
		}
		fmt.Printf(
			"\033[LSimple %d: %s => %.2f%% -- %d/%d",
			s.size,
			tool.GetDurStr(time.Since(s.start)),
			float32(i)*100.0/float32(s.size),
			i,
			s.size,
		)
	}
	fmt.Printf("\033[LSimple %d: %s => 100%%\n", s.size, tool.GetDurStr(time.Since(s.start)))
}

func (s *SimpleNQ) RunRand() {
	moves := make([][]int, s.size)
	valid := make([]int, s.size)
	for i := range valid {
		valid[i] = i
	}
	moves[0] = valid

	var sel int
	for i := 0; i < s.size && i >= 0; i++ {
		if len(moves[i]) > 0 {
			sel = rand.Intn(len(moves[i]))
			s.board[i] = moves[i][sel]
			moves[i] = append(moves[i][:sel], moves[i][sel+1:]...)
			if i < s.size-1 {
				valid = s.GetValids(i + 1)
				moves[i+1] = valid
			}
		} else {
			i -= 2
		}
		fmt.Printf(
			"\033[LSimple %d: %s => %.2f%% -- %d/%d",
			s.size,
			tool.GetDurStr(time.Since(s.start)),
			float32(i)*100.0/float32(s.size),
			i,
			s.size,
		)
	}
	fmt.Printf("\033[LSimple %d: %s => 100%%\n", s.size, tool.GetDurStr(time.Since(s.start)))
}

func (s *SimpleNQ) GetMaxIndex(sols []int) int {
	if len(sols) < 1 {
		return -1
	}
	p := 0
	max := sols[0]
	for i, m := range sols {
		if m > max {
			p = i
		}
	}
	return p
}

func (s *SimpleNQ) RunHeur() {
	moves := make([][]int, s.size)
	sols := make([][]int, s.size)
	valid := make([]int, s.size)
	for i := range valid {
		valid[i] = i
	}
	moves[0] = valid
	sols[0] = s.GetSolsForPos(0, valid)

	var sel int
	for i := 0; i < s.size && i >= 0; i++ {
		if len(moves[i]) > 0 {
			sel = s.GetMaxIndex(sols[i])
			s.board[i] = moves[i][sel]
			moves[i] = append(moves[i][:sel], moves[i][sel+1:]...)
			sols[i] = append(sols[i][:sel], sols[i][sel+1:]...)

			if i < s.size-1 {
				valid = s.GetValids(i + 1)
				moves[i+1] = valid
				sols[i+1] = s.GetSolsForPos(i+1, valid)
			}
		} else {
			i -= 2
		}
		fmt.Printf(
			"\033[LSimple %d: %s => %.2f%% -- %d/%d",
			s.size,
			tool.GetDurStr(time.Since(s.start)),
			float32(i)*100.0/float32(s.size),
			i,
			s.size,
		)
	}
	fmt.Printf("\033[LSimple %d: %s => 100%%\n", s.size, tool.GetDurStr(time.Since(s.start)))
}

func (s *SimpleNQ) RunFHeur() {
	moves := make([][]int, s.size)
	sols := make([][]int, s.size)
	valid := make([]int, s.size)
	for i := range valid {
		valid[i] = i
	}
	moves[0] = valid
	sols[0] = s.GetFSolsForPos(0, valid)

	var sel int
	for i := 0; i < s.size && i >= 0; i++ {
		if len(moves[i]) > 0 {
			sel = s.GetMaxIndex(sols[i])
			s.board[i] = moves[i][sel]
			moves[i] = append(moves[i][:sel], moves[i][sel+1:]...)
			sols[i] = append(sols[i][:sel], sols[i][sel+1:]...)

			if i < s.size-1 {
				valid = s.GetValids(i + 1)
				moves[i+1] = valid
				sols[i+1] = s.GetFSolsForPos(i+1, valid)
			}
		} else {
			i -= 2
		}
		fmt.Printf(
			"\033[LSimple %d: %s => %.2f%% -- %d/%d",
			s.size,
			tool.GetDurStr(time.Since(s.start)),
			float32(i)*100.0/float32(s.size),
			i,
			s.size,
		)
	}
	fmt.Printf("\033[LSimple %d: %s => 100%%\n", s.size, tool.GetDurStr(time.Since(s.start)))
}

func RunSimple(size int) {
	snq := NewSimpleNQ(size)
	snq.Run()
	dur := time.Since(snq.start)
	export.Export(snq.size, fmt.Sprintf("Simple %s", tool.GetDurStr(dur)), snq.board)
}

func RunSimpleRand(size int) {
	snq := NewSimpleNQ(size)
	snq.RunRand()
	dur := time.Since(snq.start)
	export.Export(snq.size, fmt.Sprintf("RandSimple %s", tool.GetDurStr(dur)), snq.board)
}

func RunSimpleHeur(size int) {
	snq := NewSimpleNQ(size)
	snq.RunHeur()
	dur := time.Since(snq.start)
	export.Export(snq.size, fmt.Sprintf("HeurSimple %s", tool.GetDurStr(dur)), snq.board)
}

func RunSimpleFHeur(size int) {
	snq := NewSimpleNQ(size)
	snq.RunFHeur()
	dur := time.Since(snq.start)
	export.Export(snq.size, fmt.Sprintf("HeurSimple %s", tool.GetDurStr(dur)), snq.board)
}
