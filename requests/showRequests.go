package requests

import (
	myMongo "DzProxy/mongo"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func ShowRequests(w http.ResponseWriter, r *http.Request) {

	requests, errGetR := myMongo.GetRequests()
	if errGetR != nil {
		log.Fatal(errGetR)
	}

	outStr := ""
	strSeparator := "\n------------------------------------------------------------------------\n"
	for _, element := range requests {
		outStr += "Id:" + strconv.Itoa(element.Id) + "  ## Method:" + element.Method + "  ## " + "Url:" + element.FullPath + strSeparator
	}
	//body True???

	fmt.Fprint(w, outStr)
}
