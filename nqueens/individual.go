package nqueens

import (
	"encoding/json"
	"math/rand"
	"os"
	"sort"
)

type Individual struct {
	Fitness int
	Genome  []int
}

func NewIndividual(Size int) *Individual {
	genome := make([]int, Size)
	for i := 0; i < Size; i++ {
		genome[i] = i
	}
	sort.Slice(genome, func(i, j int) bool {
		if rand.Float64() < 0.5 {
			return true
		}
		return false
	})

	ind := &Individual{
		Genome: genome,
	}
	return ind
}

func (ind *Individual) Mutate() {
	g1 := rand.Intn(len(ind.Genome))
	g2 := rand.Intn(len(ind.Genome))
	for g1 == g2 {
		g2 = rand.Intn(len(ind.Genome))
	}

	ind.Genome[g1], ind.Genome[g2] = ind.Genome[g2], ind.Genome[g1]
}

func (ind *Individual) CalFitness() {
	errors := 0
	var df int
	for i := 0; i < len(ind.Genome)-1; i++ {
		for j := i + 1; j < len(ind.Genome); j++ {
			df = (ind.Genome)[i] - (ind.Genome)[j]
			if df < 0 {
				df = -df
			}
			if df == 0 || df == j-i {
				errors++
			}
		}
	}
	ind.Fitness = errors
}

func (ind *Individual) HasInvalidGen() bool {
	for _, e := range ind.Genome {
		if e == -1 {
			return true
		}
	}
	return false
}

func (ind *Individual) FindGen(g int) int {
	for i, e := range ind.Genome {
		if e == g {
			return i
		}
	}
	return -1
}

func (ind *Individual) GenIndexs() []int {
	index := make([]int, len(ind.Genome))
	for i, v := range ind.Genome {
		index[v] = i
	}
	return index
}

func (ind *Individual) Save() {
	bytes, _ := json.Marshal(map[string]any{
		"size":   len(ind.Genome),
		"result": ind.Genome,
	})
	os.WriteFile("nqueens_result.json", bytes, 0777)
}
