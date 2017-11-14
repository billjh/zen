package main

import "net/http"
import "io/ioutil"
import "fmt"

func zen() (string, error) {
	client := &http.Client{}
	res, err := client.Get("https://api.github.com/zen")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return string(body), nil
	}
	return "", fmt.Errorf("Status is not 200")
}

func main() {
	words := make(map[string]string)

	for len(words) < 3 {
		word, err := zen()
		if err != nil {
			continue
		}
		if _, ok := words[word]; ok {
			continue
		}
		words[word] = word
		fmt.Println(word)
	}
}
