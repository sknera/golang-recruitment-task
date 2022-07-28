package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type randInt struct {
	Data   []int   `json:"data"`
	Stddev float64 `json:"stddev"`
}

var wg = sync.WaitGroup{}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/random/mean/", getApi)

	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}

func getApi(w http.ResponseWriter, req *http.Request) {
	r, _ := strconv.Atoi(req.URL.Query().Get("requests"))
	l, _ := strconv.Atoi(req.URL.Query().Get("length"))

	// sum of all sets will be on last index
	var randInts = make([]randInt, r+1)

	for i := 0; i < r; i++ {
		wg.Add(1)
		go createRandInt(i, l, randInts)

	}

	wg.Wait()
	var sumSet []int
	for _, todo := range randInts[:len(randInts)-1] {
		sumSet = append(sumSet, todo.Data...)
	}

	randInts[len(randInts)-1].Data = sumSet
	randInts[len(randInts)-1].Stddev = calcStddev(randInts[len(randInts)-1].Data)

	fmt.Printf("all randints is %v", sumSet)
	json.NewEncoder(w).Encode(randInts)
	w.Header().Add("Content-Type", "application/json")
	fmt.Printf("params are %v and %v\n", r, l)
	w.WriteHeader(http.StatusOK)

}

func calcStddev(data []int) float64 {
	sum := 0
	for _, item := range data {
		sum += item
	}
	mean := float64(sum) / float64(len(data))
	var sumSquares = 0.0
	for _, item := range data {
		sumSquares += math.Pow(mean-float64(item), 2)
	}
	return math.Sqrt(sumSquares / float64(len(data)))
}

func getRandom(l int) (string, error) {
	url := fmt.Sprintf("https://www.random.org/integers/?num=%v&min=1&max=6&col=1&base=10&format=plain&rnd=new", l)
	response, err := http.Get(url)
	data, _ := ioutil.ReadAll(response.Body)
	return string(data), err
}

func createRandInt(i int, l int, randInts []randInt) {
	requestStr, _ := getRandom(l)
	requestInts := strings.Split(requestStr, "\n")

	for _, item := range requestInts[:len(requestInts)-1] {
		requestInt, _ := strconv.Atoi(item)
		randInts[i].Data = append(randInts[i].Data, requestInt)
	}

	randInts[i].Stddev = calcStddev(randInts[i].Data)
	// 3 seconds, just to show concurrency
	time.Sleep(3 * time.Second)
	wg.Done()
}
