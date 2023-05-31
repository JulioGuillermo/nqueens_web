package nqueens

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/julioguillermo/nqueens_web/export"
	"github.com/julioguillermo/nqueens_web/tool"
)

type GA struct {
	Population []*Individual
	MutRate    float64
	Survivors  int
	Generation uint64
	GenSize    int
	MaxErrors  int
	Threads    int
	MC         bool
	start      time.Time
	dur        time.Duration
}

func NewGA(N int, popsize int, survivors int, mutrate float64, mc bool, threads int) *GA {
	pop := make([]*Individual, popsize)
	for i := 0; i < popsize; i++ {
		pop[i] = NewIndividual(N)
	}
	ga := &GA{
		MutRate:    mutrate,
		Population: pop,
		Survivors:  survivors,
		GenSize:    N,
		Generation: 0,
		Threads:    threads,
		MC:         mc,
		start:      time.Now(),
	}
	ga.CalMaxErrors()
	ga.Evaualte()
	return ga
}

func (ga *GA) CalMaxErrors() {
	errs := 0
	for i := ga.GenSize - 1; i > 0; i-- {
		errs += i
	}
	ga.MaxErrors = errs
}

func (ga *GA) Evaualte() {
	if ga.Threads > 1 {
		ga.AsyncEvaluate()
	} else {
		for _, p := range ga.Population {
			p.CalFitness()
		}
	}

	sort.Slice(ga.Population, func(i, j int) bool {
		return ga.Population[i].Fitness < ga.Population[j].Fitness
	})
}

func (ga *GA) EvalInd(ind chan int, wg *sync.WaitGroup) {
	for i := range ind {
		ga.Population[i].CalFitness()
		wg.Done()
	}
}

func (ga *GA) AsyncEvaluate() {
	ind := make(chan int, ga.Threads)
	var wg sync.WaitGroup
	for i := 0; i < ga.Threads; i++ {
		go ga.EvalInd(ind, &wg)
	}

	for i := range ga.Population {
		wg.Add(1)
		ind <- i
	}

	wg.Wait()
}

func (ga *GA) CrossInd(ind int) {
	father := rand.Intn(ga.Survivors)
	mother := rand.Intn(ga.Survivors)
	for mother == father {
		mother = rand.Intn(ga.Survivors)
	}

	end := rand.Intn(ga.GenSize)
	start := rand.Intn(ga.GenSize)
	for end == start {
		start = rand.Intn(ga.GenSize)
	}
	if end < start {
		end, start = start, end
	}

	others := make([]int, ga.GenSize)
	for i := 0; i < ga.GenSize; i++ {
		if i >= start && i < end {
			ga.Population[ind].Genome[i] = ga.Population[father].Genome[i]
			others[i] = -1
		} else {
			ga.Population[ind].Genome[i] = -1
			others[i] = ga.Population[father].Genome[i]
		}
	}
	motherIndex := ga.Population[mother].GenIndexs()
	sort.SliceStable(others, func(i, j int) bool {
		if others[i] == -1 {
			return false
		}
		if others[j] == -1 {
			return true
		}
		return motherIndex[others[i]] < motherIndex[others[j]]
		// return ga.Population[mother].FindGen(others[i]) < ga.Population[mother].FindGen(others[j])
	})
	oi := 0
	for i := 0; i < ga.GenSize; i++ {
		if ga.Population[ind].Genome[i] == -1 {
			ga.Population[ind].Genome[i] = others[oi]
			oi++
		}
	}
}

func (ga *GA) MapCrossInd(ind int) {
	father := rand.Intn(ga.Survivors)
	mother := rand.Intn(ga.Survivors)
	for mother == father {
		mother = rand.Intn(ga.Survivors)
	}
	end := rand.Intn(ga.GenSize)
	start := rand.Intn(ga.GenSize)
	for end == start {
		start = rand.Intn(ga.GenSize)
	}
	if end < start {
		end, start = start, end
	}
	for i := 0; i < ga.GenSize; i++ {
		if i >= start && i < end {
			ga.Population[ind].Genome[i] = ga.Population[father].Genome[i]
		} else {
			ga.Population[ind].Genome[i] = -1
		}
	}
	motherIndexs := ga.Population[mother].GenIndexs()
	for i := start; i < end && ga.Population[ind].HasInvalidGen(); i++ {
		gen := ga.Population[mother].Genome[i]
		if ga.Population[ind].FindGen(gen) != -1 {
			continue
		}
		index := motherIndexs[ga.Population[ind].Genome[i]]
		for ga.Population[ind].Genome[index] != -1 {
			index = motherIndexs[ga.Population[ind].Genome[index]]
			// index = ga.Population[mother].FindGen(ga.Population[ind].Genome[index])
		}
		ga.Population[ind].Genome[index] = gen
	}
	for i, v := range ga.Population[ind].Genome {
		if v != -1 {
			continue
		}
		ga.Population[ind].Genome[i] = ga.Population[mother].Genome[i]
	}
}

func (ga *GA) AsyncCrossInd(ind chan int, wg *sync.WaitGroup) {
	for i := range ind {
		if ga.MC {
			ga.MapCrossInd(i)
		} else {
			ga.CrossInd(i)
		}
		wg.Done()
	}
}

func (ga *GA) AsyncCross() {
	var wg sync.WaitGroup
	ind := make(chan int, ga.Threads)
	for i := 0; i < ga.Threads; i++ {
		go ga.AsyncCrossInd(ind, &wg)
	}
	for i := ga.Survivors; i < len(ga.Population); i++ {
		wg.Add(1)
		ind <- i
	}
	wg.Wait()
}

func (ga *GA) Cross() {
	if ga.Threads > 1 {
		ga.AsyncCross()
	} else {
		for i := ga.Survivors; i < len(ga.Population); i++ {
			if ga.MC {
				ga.MapCrossInd(i)
			} else {
				ga.CrossInd(i)
			}
		}
	}
}

func (ga *GA) AsyncMutInd(ind chan int, wg *sync.WaitGroup) {
	for i := range ind {
		ga.Population[i].Mutate()
		wg.Done()
	}
}

func (ga *GA) AsyncMut() {
	var wg sync.WaitGroup
	ind := make(chan int, ga.Threads)
	for i := 0; i < ga.Threads; i++ {
		go ga.AsyncMutInd(ind, &wg)
	}
	for i := ga.Survivors; i < len(ga.Population); i++ {
		wg.Add(1)
		ind <- i
	}
	wg.Wait()
}

func (ga *GA) Mutate() {
	if ga.Threads > 1 {
		ga.AsyncMut()
	} else {
		for i := ga.Survivors; i < len(ga.Population); i++ {
			if rand.Float64() < ga.MutRate {
				ga.Population[i].Mutate()
			}
		}
	}
}

func (ga *GA) NextGen() {
	ga.Cross()
	ga.Mutate()
	ga.Evaualte()
	ga.Generation++
}

func (ga *GA) GetBestError() int {
	return ga.Population[0].Fitness
}

func (ga *GA) Info() (uint64, float64, time.Duration, int, int) {
	dur := time.Since(ga.start)
	ga.dur = dur
	bf := ga.GetBestError()
	pro := float64(bf) / float64(ga.MaxErrors)
	pro = 100 - pro*100
	pro = (math.Exp(pro)) / (math.Exp(100))
	pro = pro * 100
	return ga.Generation, pro, dur, bf, ga.MaxErrors
}

func (ga *GA) Save() (string, error) {
	return export.Export(
		ga.GenSize,
		fmt.Sprintf("%d - %s", ga.Generation, tool.GetDurStr(ga.dur)),
		ga.Population[0].Genome,
	)
}
