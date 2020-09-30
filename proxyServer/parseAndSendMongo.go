package proxyServer

import (
	myMongo "DzProxy/mongo"
	"net/http"
)

func parseAndSendMongo(r *http.Request, body *[]byte, urlForDump *string) (errParse error) {

	request := myMongo.MyRequest{
		Method: r.Method}

	request.FullPath = *urlForDump

	request.Headers = make(map[string][]string)
	for name, values := range r.Header {
		request.Headers[name] = values
	}

	request.Body = *body
	errParse = myMongo.AddRequest(request)
	return errParse
}
