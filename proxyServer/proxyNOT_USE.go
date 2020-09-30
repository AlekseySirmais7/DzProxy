package proxyServer

/*
type Proxy struct {

	// корневой для генерации сертиф для входящих
	RootCA *tls.Certificate

	TLSServerConfig *tls.Config

	// таймаут на чтение body
	FlushInterval time.Duration

	// только для дебага
	PrintMutex sync.Mutex
}



var dumpAndServe = func(upstream http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		dumpFunc := func(bodyBytes *[]byte) {
			// send to mongo
			//parseAndSendMongo(r, bodyBytes)
			//log.Println("\n\nDUMP\n" + string(dump) + "------------------------\n")
		}
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		go dumpFunc(&bodyBytes)
		r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))


		upstream.ServeHTTP(w, r)
	})
}




func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "CONNECT" {
		p.decryptAndServe(w, r)
		return
	}
	rp := &httputil.ReverseProxy{
		Director:      httpDirector,
		FlushInterval: p.FlushInterval,
	}

	dumpAndServe(rp).ServeHTTP(w, r)
}

func (p *Proxy) decryptAndServe(w http.ResponseWriter, r *http.Request) {

	name, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		name = ""
	}

	newDomainCert, err := p.getCertForDomain(name)
	if err != nil {
		log.Println("getCertForDomain", err)
		http.Error(w, "no upstream", 503)
		return
	}

	sConfig := new(tls.Config)
	if p.TLSServerConfig == nil {
		panic("need set TLSServerConfig")
	}

	*sConfig = *p.TLSServerConfig

	var secureConn *tls.Conn

	sConfig.Certificates = []tls.Certificate{*newDomainCert}
	sConfig.GetCertificate = func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		cConfig := new(tls.Config)
		cConfig.ServerName = hello.ServerName
		secureConn, err = tls.Dial("tcp", r.Host, cConfig)
		if err != nil {
			log.Println("dial", r.Host, err)
			return nil, err
		}
		return p.getCertForDomain(hello.ServerName)
	}

	conAfterHandshake, err := handshake(w, sConfig)
	if err != nil {
		log.Println("handshake ERR", r.Host, err)
		return
	}
	defer conAfterHandshake.Close()
	if secureConn == nil {
		log.Println("could not determine getCertForDomain name for " + r.Host)
		return
	}

	//111
	defer secureConn.Close()


	od := &oneShotDialer{c: secureConn,
		//timer: *time.NewTimer(time.Second * 5),
	}

	//od.timer.Reset(time.Second * 5)

	rp := &httputil.ReverseProxy {
		Director:      httpsDirector,
		Transport:     &http.Transport{DialTLS: od.Dial},
		FlushInterval: p.FlushInterval,
	}


	ch := make(chan int)
	waitClose := &callFuncOnClose{conAfterHandshake, func() { ch <- 0 }}
	http.Serve(&oneShotListener{waitClose}, dumpAndServe(rp))
	<-ch


}




func (p *Proxy) getCertForDomain(names ...string) (*tls.Certificate, error) {
	return genCert(p.RootCA, names)
}

var okHeader = []byte("HTTP/1.1 200 OK\r\n\r\n")

func handshake(w http.ResponseWriter, config *tls.Config) (net.Conn, error) {
	raw, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		http.Error(w, "no upstream", 503)
		return nil, err
	}
	if _, err = raw.Write(okHeader); err != nil {
		raw.Close()
		return nil, err
	}
	conn := tls.Server(raw, config)
	err = conn.Handshake()
	if err != nil {
		conn.Close()
		raw.Close()
		return nil, err
	}
	return conn, nil
}

func httpDirector(r *http.Request) {
	r.URL.Host = r.Host
	r.URL.Scheme = "http"
}

func httpsDirector(r *http.Request) {
	r.URL.Host = r.Host
	r.URL.Scheme = "https"
}







// реализует net.Dialer; отдает соединение пока оно не закрыто
// изменил
type oneShotDialer struct {
	c  net.Conn
	mu sync.Mutex
	timer time.Timer

}

func (d *oneShotDialer) Dial(network, addr string) (net.Conn, error) {


	//if !d.timer.Stop() {
	//	d.c = nil
	//	return nil, errors.New("Dial()::closed")
//	}
	//return d.c, nil



	d.mu.Lock()
	defer d.mu.Unlock()
	if d.c == nil {
		return nil, errors.New("Dial()::closed")
	}
	c := d.c
	d.c = nil
	return c, nil
}


type oneShotListener struct {
	connect net.Conn
}

func (l *oneShotListener) Accept() (net.Conn, error) {
	if l.connect == nil {
		return nil, errors.New(" Accept()::closed")
	}
	c := l.connect
	l.connect = nil
	return c, nil
}

func (l *oneShotListener) Close() error {
	return nil
}

func (l *oneShotListener) Addr() net.Addr {
	return l.connect.LocalAddr()
}

// A callFuncOnClose implements net.Conn and calls its f on Close.
type callFuncOnClose struct {
	net.Conn
	f func()
}

func (c *callFuncOnClose) Close() error {
	if c.f != nil {
		c.f()
		c.f = nil
	}
	return c.Conn.Close()
}
*/
