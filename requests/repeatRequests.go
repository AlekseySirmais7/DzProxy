package requests

import (
	myMongo "DzProxy/mongo"
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var tr = &http.Transport{
	MaxIdleConns:       10,
	IdleConnTimeout:    5 * time.Second,
	DisableCompression: true,
}

func RepeatRequests(w http.ResponseWriter, r *http.Request) {
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

	pathRepeat := "https://" + request.FullPath
	if strings.Contains(request.FullPath, "http/") || strings.Contains(request.FullPath, "http:") {
		pathRepeat = request.FullPath
	}

	var newRequest, err = http.NewRequest(request.Method, pathRepeat, bytes.NewBuffer(request.Body))
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
	defer resp.Body.Close()

	respHeaders := ""
	for name, values := range r.Header {
		respHeaders += name + ":   " + values[0] + "\n"
	}
	respHeaders += "-------------------------------------------------------"

	body, _ := ioutil.ReadAll(resp.Body)

	strSeparator := "====================================================================\n"
	outStr := "Id:" + strconv.Itoa(request.Id) + "  ## Method:" + request.Method + "  ## " + "Url:" + request.FullPath + strSeparator
	outStr += "Response status:" + resp.Status + " \nResponse Headers:" + respHeaders + "\n\nResponse Body:\n" + fmt.Sprintf("%s", body)

	fmt.Fprint(w, outStr)
}
