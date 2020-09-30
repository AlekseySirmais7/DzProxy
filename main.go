package main

import (
	myMongo "DzProxy/mongo"
	myProxy "DzProxy/proxyServer"
	myRequests "DzProxy/requests"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"path"
	"runtime"
)

const folderForCerts = "/home/alex/go/src/DzProxy/certs"

func main() {

	runtime.GOMAXPROCS(10)

	errMongo := myMongo.ConnectMongo()
	if errMongo != nil {
		log.Fatal(errMongo)
	}
	log.Println("mongo connect is ok")

	myRequests.LoadParams(folderForCerts + "/params")

	certSet := myProxy.СertSettings{
		Folder:       folderForCerts,
		CertName:     "alexCAcentre",
		RootKeyFile:  path.Join(folderForCerts, "ca-key.pem"),
		RootCertFile: path.Join(folderForCerts, "ca-cert.pem"),
	}

	ca, errLoadCa := certSet.LoadCA()
	if errLoadCa != nil {
		log.Fatal(errLoadCa)
	}
	log.Println("root certs is OK")
	/*
		proxy := &myProxy.Proxy{
			RootCA: &ca,
			TLSServerConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
			FlushInterval: time.Second * 10,
		}
	*/
	ch := make(chan error, 0)

	// proxy
	log.Println("start serve proxy port:8080:")
	go func(chan error) {
		//http.ListenAndServe(":8080", proxy)

		certCa := myProxy.CertCa{CertCaPtr: &ca}
		proxyHand := http.HandlerFunc(certCa.ProxyHandler)
		http.ListenAndServe(":8080", proxyHand)

	}(ch)

	// requests interface
	log.Println("start requests interface  port :3000")
	go func(chan error) {
		router := mux.NewRouter()
		router.HandleFunc("/request/{id:[0-9]+}", myRequests.RepeatRequests)
		router.HandleFunc("/requests", myRequests.ShowRequests)
		router.HandleFunc("/attack/{id:[0-9]+}", myRequests.RunAttack)
		http.ListenAndServe(":3000", router)
	}(ch)

	log.Fatal(<-ch)
}
