package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vearne/gin-timeout"
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

	router := gin.Default()
	defaultMsg := `{"code": -1, "msg":"http: Handler timeout"}`
	router.Use(timeout.Timeout(
		timeout.WithTimeout(6*time.Second),
		timeout.WithErrorHttpCode(http.StatusRequestTimeout),
		timeout.WithDefaultMsg(defaultMsg),
		timeout.WithCallBack(func(r *http.Request) {
			fmt.Println("timeout happen, url:", r.URL.String())
		})))
	router.GET("/random/mean/", getApi)

	fmt.Println("Listening on port 8080")
	log.Fatal(router.Run(":8080"))

}

func getApi(c *gin.Context) {

	r, _ := strconv.Atoi(c.Query("requests"))
	l, _ := strconv.Atoi(c.Query("length"))

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

	c.IndentedJSON(http.StatusOK, randInts)

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
	requestStr, err := getRandom(l)
	if err != nil {
		fmt.Printf("Error using random.org api: %s\n", err.Error())
	}

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
