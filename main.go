package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var regRoll = regexp.MustCompile(`(\d+)?[Dd](\d+)(([ab\+\-])(\d+))?`)

func main() {
	var roll string
	var history []string
	var curIndex int
	regRepeat := regexp.MustCompile(`[r]([\-\s]?)+\d+`)
	file, err := os.ReadFile("apikey")
	if err != nil {
		log.Fatal("no apikey file found")
	}

	apikey := strings.Trim(string(file), "\n")

	for {
		fmt.Println("input dice roll:")
		fmt.Scanln(&roll)

		if roll == "quit" ||
			roll == "exit" ||
			roll == "q" {
			quit()
		}

		switch {
		case regRoll.MatchString(roll):
			history = append(history, roll)

			printDiceRoll(apikey, roll, curIndex)
			curIndex++
		case regRepeat.MatchString(roll):
			regNum := regexp.MustCompile("[0-9]+")
			index, _ := strconv.Atoi(strings.Join(regNum.FindAllString(roll, -1), ""))

			if index < len(history) {
				history = append(history, history[index])

				fmt.Println("repeating:", history[index])
				printDiceRoll(apikey, history[index], curIndex)
				curIndex++
			} else {
				color.Red("index out of bounds")
			}
		case roll == "r":
			if len(history) != 0 {
				history = append(history, history[len(history)-1])

				fmt.Println("repeating last")
				printDiceRoll(apikey, history[len(history)-1], curIndex)
				curIndex++
			} else {
				fmt.Println("No history")
			}
		default:
			color.Blue("Format your dice rolls like this:")
			// Make it color coded and fancy
			fmt.Printf("%s%s%s%s%s, %s\n",
				color.HiRedString("n"),
				color.BlueString("d"),
				color.HiMagentaString("p"),
				color.HiYellowString("+"),
				color.HiCyanString("q"),
				color.BlueString("for example 'd10', '2d20', '2d12-5', or '10d6a3' etc"),
			)
			fmt.Println(color.HiRedString("n:"), " n is number of dice (optional)")
			fmt.Println(color.HiMagentaString("p:"), " is type of dice")
			fmt.Println(color.HiYellowString("+:"), " is the type of modifier")
			fmt.Println(color.HiCyanString("q:"), " is the modifier")

			color.Blue("\nThe types of modifier are:")
			fmt.Println(color.HiYellowString("+:"), "Add a constant number to roll")
			fmt.Println(color.HiYellowString("-:"), "Subtract a constant number from roll")
			fmt.Println(color.HiYellowString("a:"), "Show all dice which rolled the modifier number and above")
			fmt.Println(color.HiYellowString("b:"), "Show all dice which rolled the modifier number and below")

			color.Blue("\nOther commands:")
			fmt.Println(color.HiMagentaString("r or [enter]:"), " for repeat last")
			fmt.Println(color.HiMagentaString("r[n]:"), " where n is a specific roll to repeat")
			fmt.Println(color.HiMagentaString("q:"), " to quit")
		}
	}
}

func printDiceRoll(apikey string, roll string, curIndex int) {
	var num int
	var result []int
	input := regRoll.FindStringSubmatch(roll)[1:]

	if input[0] != "" {
		num, _ = strconv.Atoi(input[0])
	} else {
		num = 1
	}

	diceType, _ := strconv.Atoi(input[1])
	modifier := input[3]
	modValue, _ := strconv.Atoi(input[4])
	var modResult []int

	result, sum := fetchDiceRoll(apikey, num, diceType)

	switch modifier {
	case "-":
		modResult = append(modResult, sum-modValue)
	case "+":
		modResult = append(modResult, sum+modValue)
	case "a":
		for _, v := range result {
			if v >= modValue {
				modResult = append(modResult, v)
			}
		}
	case "b":
		for _, v := range result {
			if v <= modValue {
				modResult = append(modResult, v)
			}
		}
	}

	color.Set(color.FgMagenta)
	fmt.Printf("%d: ", curIndex)

	color.Set(color.FgBlue)
	fmt.Printf("%d ", result)

	color.Set(color.FgHiYellow)
	if len(modResult) > 0 {
		if len(modResult) > 1 {
			fmt.Printf("sum: %d, successes:%d, total: %d\n", sum, modResult, len(modResult))
		} else {
			fmt.Printf("sum: %d, result: %d\n", sum, modResult[0])
		}
	} else {
		fmt.Printf("sum: %d\n", sum)
	}
	defer color.Unset()
}

func quit() {
	fmt.Println("quitting...")
	os.Exit((0))
}

func fetchDiceRoll(apikey string, num int, diceType int) ([]int, int) {
	postBody, _ := json.Marshal(RandomOrgRequest{
		Jsonrpc: "2.0",
		Method:  "generateIntegers",
		ID:      1337,
		Params: Params{
			APIKey:      apikey,
			N:           num,
			Min:         1,
			Max:         diceType,
			Replacement: true,
		},
	})
	responseBody := bytes.NewBuffer(postBody)

	url := "https://api.random.org/json-rpc/4/invoke"
	response, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		log.Fatalf("Failed to create resource at: %s and the error is: %v\n", url, err)
	}

	defer response.Body.Close()

	var data RandomOrgResponse
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	sum := 0
	for _, value := range data.Result.Random.Data {
		sum += value
	}

	return data.Result.Random.Data, sum
}

type RandomOrgResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Random struct {
			Data           []int  `json:"data"`
			CompletionTime string `json:"completionTime"`
		} `json:"random"`
		BitsUsed      int `json:"bitsUsed"`
		BitsLeft      int `json:"bitsLeft"`
		RequestsLeft  int `json:"requestsLeft"`
		AdvisoryDelay int `json:"advisoryDelay"`
	} `json:"result"`
	ID int `json:"id"`
}

type RandomOrgRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
	ID      int    `json:"id"`
}

type Params struct {
	APIKey      string `json:"apiKey"`
	N           int    `json:"n"`
	Min         int    `json:"min"`
	Max         int    `json:"max"`
	Replacement bool   `json:"replacement"`
}
