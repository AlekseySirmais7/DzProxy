package requests

import (
	myMongo "DzProxy/mongo"
	"bufio"
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var paramSlice []string

const paramValue = "shefuisehfuishe123"

func RunAttack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, errGetInt := strconv.Atoi(idStr)
	if errGetInt != nil {
		panic("need use integer id!")
	}

	request, errGetRequest := myMongo.GetOneRequest(id)
	if errGetRequest != nil {
		log.Fatal(errGetRequest)
	}

	response := fmt.Sprintf("Request: method:%s\t%s\t\tID:%d\n\n\n", request.Method, request.FullPath, request.Id)
	fmt.Fprint(w, response)

	for _, element := range paramSlice {

		addStr := element + "=" + paramValue

		if strings.Contains(request.FullPath, "?") {
			addStr = "&" + addStr
		} else {
			addStr = "?" + addStr
		}

		pathRequest := request.FullPath + addStr

		if !strings.Contains(pathRequest, "http:") {
			pathRequest = "https://" + pathRequest
		}

		var newRequest, err = http.NewRequest(request.Method, pathRequest,
			bytes.NewBuffer(request.Body))
		if err != nil {
			panic(err)
		}

		for key, element := range request.Headers {
			newRequest.Header.Add(key, element[0])
		}

		client := &http.Client{Transport: tr}
		resp, err := client.Do(newRequest)
		if err != nil {
			panic(err)
		}

		log.Println("Check:" + request.FullPath + addStr + "   STATUS:" + resp.Status)

		respHeaders := ""
		for name, values := range r.Header {
			respHeaders += name + ":   " + values[0] + "\n"
		}

		outStr := "URL:" + request.FullPath + "\n-------------------------------------\n"
		if strings.Contains(respHeaders, paramValue) {

			outStr += "Headers contain input GET param:" + element + " (value:" + paramValue + ")\n"
			outStr += "Checked url:" + request.FullPath + addStr + "\n"
			fmt.Fprint(w, outStr)
			resp.Body.Close()
			return
		}

		body, _ := ioutil.ReadAll(resp.Body)

		if strings.Contains(string(body), paramValue) {

			outStr += "Body contains input GET param:" + element + " (value:" + paramValue + ")\n"
			outStr += "Checked url:" + request.FullPath + addStr + "\n"
			fmt.Fprint(w, outStr)
			resp.Body.Close()
			return
		}
		resp.Body.Close()
	}

	outStr := "Not find get value!"
	fmt.Fprint(w, outStr)
}

func LoadParams(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		paramSlice = append(paramSlice, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
