// Copyright 2025 The DGM Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Prompt is a llm prompt
type Prompt struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// Query submits a query to the llm
func Query(query string) string {
	prompt := Prompt{
		Model:  "llama3.2",
		Prompt: query,
	}
	data, err := json.Marshal(prompt)
	if err != nil {
		panic(err)
	}
	buffer := bytes.NewBuffer(data)
	response, err := http.Post("http://10.0.0.54:11434/api/generate", "application/json", buffer)
	if err != nil {
		panic(err)
	}
	reader, answer := bufio.NewReader(response.Body), ""
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		data := map[string]interface{}{}
		err = json.Unmarshal([]byte(line), &data)
		text := data["response"].(string)
		answer += text
	}
	return answer
}

func main() {
	//fmt.Println(Query("Hello World!"))
	const (
		begin = "```javascript"
		end   = "```"
	)
	const prompt = "Use the following fitness function for optimization to zero, place all code in a single code block"
	result, i := Query(prompt+`: 
	`+begin+`
	var bestSolution = 1;
	var bestFitness = 12345;
	function fitness(solution) {
		var fitness = 12345 % solution;
		if (fitness < bestFitness) {
			bestSolution = solution;
			bestFitness = fitness;
			llama.best(solution, fitness);
		}
		return fitness;
	}
`+end+`
`), 0
	for {
		previous := ""
		js := ""
		for {
			goja := NewGOJA()
			index := strings.Index(result, begin)
			if index == -1 {
				fmt.Print(result)
				break
			}
			fmt.Print(result[:index+len(begin)])
			result = result[index+len(begin):]
			index = strings.Index(result, end)
			fmt.Println(result[:index+len(end)])
			fmt.Println("```goja")
			js = result[:index]
			err := goja.Run(i, js)
			if err != nil {
				fmt.Print("<<<")
				fmt.Println(err)
				js = previous
				break
			}
			i++
			fmt.Println("```")
			result = result[index+len(end):]
			fmt.Println("best", goja.Answer, goja.Cost)
		}
		previous = js
		result, i = Query(`Improve the following integer factoring code. `+prompt+`: 
			`+begin+`
				var bestSolution = 1;
				var bestFitness = 12345;
				function fitness(solution) {
					var fitness = 12345 % solution;
					if (fitness < bestFitness) {
						bestSolution = solution;
						bestFitness = fitness;
						llama.best(solution, fitness);
					}
					return fitness;
				}
				`+js+"\n"), 0
	}
}
