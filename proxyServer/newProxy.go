package proxyServer

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

type CertCa struct {
	CertCaPtr *tls.Certificate
}

func (c *CertCa) ProxyHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "CONNECT" {

		name, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			name = ""
		}

		cert, errGen := c.genCertForSite(name)
		if errGen != nil {
			panic(errGen)
		}

		config := tls.Config{
			Certificates:       []tls.Certificate{*cert},
			InsecureSkipVerify: true,
		}

		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "ERR hijack", http.StatusInternalServerError)
			return
		}
		con, buff, err := hj.Hijack()
		buff.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		buff.Flush()


		timer2 := time.NewTimer(time.Second)

		timer2.Reset(time.Second * 10)

		go (func() {
			if err != nil {
				return
			}
			defer con.Close()
			tlsCon := tls.Server(con, &config)
			clientTlsReader := bufio.NewReader(tlsCon)
			clientTlsWriter := bufio.NewWriter(tlsCon)
			tlsCon.Handshake()

			for {






				r, err := http.ReadRequest(clientTlsReader)
				if err != nil {
					//fmt.Println("EOF")
					return
				}

				dumpFunc := func(bodyBytes *[]byte, urlForDump *string) {
					// send to mongo
					parseAndSendMongo(r, bodyBytes, urlForDump)
					//log.Println("\n\nDUMP\n" + string(dump) + "------------------------\n")
				}
				urlForDump := r.Host + r.RequestURI
				bodyBytes, _ := ioutil.ReadAll(r.Body)
				go dumpFunc(&bodyBytes, &urlForDump)
				r.Body.Close()
				r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

				r.RequestURI = ""
				r.URL, _ = url.Parse("https://" + r.Host + r.URL.String())

				response, err := http.DefaultTransport.RoundTrip(r)

				if err != nil {
					//break
					continue
				}
				//fmt.Println(response)

				//EncodeResponse(response)
				response.Write(clientTlsWriter)
				clientTlsWriter.Flush()


				if timer2.Stop() {
					return
				}
			}
		})()
		return
	}

	httpClient := &http.Client{}

	dumpFunc := func(bodyBytes *[]byte, urlForDump *string) {
		// send to mongo
		parseAndSendMongo(r, bodyBytes, urlForDump)
		//log.Println("\n\nDUMP\n" + string(dump) + "------------------------\n")
	}
	urlForDump := r.Host + r.RequestURI
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	go dumpFunc(&bodyBytes, &urlForDump)
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	r.RequestURI = ""
	response, _ := httpClient.Do(r)

	copyHeaders(w.Header(), response.Header)
	defer response.Body.Close()
	io.Copy(w, response.Body)
}

func (c *CertCa) genCertForSite(names ...string) (*tls.Certificate, error) {
	return genCert(c.CertCaPtr, names)
}

func copyHeaders(dest http.Header, source http.Header) {
	for header := range source {
		dest.Add(header, source.Get(header))
	}
}

func EncodeResponse(r *http.Response) {

	raw, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	buf1 := bytes.NewBuffer(raw)
	bufReader := ioutil.NopCloser(buf1)
	r.Body = bufReader

	buf2 := bytes.NewBuffer(raw)

	var reader io.ReadCloser
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(buf2)
		defer reader.Close()
	default:
		reader = ioutil.NopCloser(buf2)
	}

}
