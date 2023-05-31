package main

import (
	"log"
	"net/http"
)

func main() {
	a := GetApp()
	http.Handle("/", a)

	if err := http.ListenAndServe(":9988", nil); err != nil {
		log.Fatal(err)
	}

	//N := flag.Int("n", 100, "size of the board and number of queens")
	//popsize := flag.Int("ps", 1000, "population size for GA")
	//survivors := flag.Int("s", 10, "survivors for GA")
	//mutrate := flag.Float64("mr", 0.2, "mutation rate for GA")
	//threads := flag.Int("th", 10, "number of threads for GA")
	//mc := flag.Bool("mc", false, "use mapped crossover for GA")
	//
	//simpleAlg := flag.Bool("simple", false, "use Iterative Improvement Algorithm instead of GA")
	//simpleRandAlg := flag.Bool(
	//	"rsimple",
	//	false,
	//	"use Iterative Improvement Algorithm instead of GA with random selection",
	//)
	//simpleHeurAlg := flag.Bool(
	//	"hsimple",
	//	false,
	//	"use Iterative Improvement Algorithm instead of GA with heuristic selection",
	//)
	//simpleFHeurAlg := flag.Bool(
	//	"fhsimple",
	//	false,
	//	"use Iterative Improvement Algorithm instead of GA with a fast heuristic selection",
	//)
	//
	//flag.Parse()
	//
	//if *simpleAlg {
	//	simple.RunSimple(*N)
	//} else if *simpleRandAlg {
	//	simple.RunSimpleRand(*N)
	//} else if *simpleHeurAlg {
	//	simple.RunSimpleHeur(*N)
	//} else if *simpleFHeurAlg {
	//	simple.RunSimpleFHeur(*N)
	//} else {
	//	nqueens.RunGANQ(*N, *popsize, *survivors, *mutrate, *mc, *threads)
	//}
}
