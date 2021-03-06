package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/hiroakis/tor-status-proxy/tor"
)

var status *tor.Status

func allNodesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Last-Modified", status.AllNodeLastModified().Format(http.TimeFormat))
	w.WriteHeader(200)
	w.Write(status.RawAllNodes())
}

func isTorNodesHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.FormValue("ip")
	if status.IsTorNode(ip) {
		w.Header().Set("Last-Modified", status.ExitNodeLastModified().Format(http.TimeFormat))
		w.WriteHeader(200)
	} else {
		w.WriteHeader(404)
	}
}

func exitNodesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Last-Modified", status.ExitNodeLastModified().Format(http.TimeFormat))
	w.WriteHeader(200)
	w.Write(status.RawExitNodes())
}

func isExitNodeHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.FormValue("ip")
	if status.IsExitNode(ip) {
		w.Header().Set("Last-Modified", status.AllNodeLastModified().Format(http.TimeFormat))
		w.WriteHeader(200)
	} else {
		w.WriteHeader(404)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var ip string

	xff := r.Header.Get("X-FORWARDED-FOR")

	if xff == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	} else {
		if strings.Contains(xff, ",") {
			ip = strings.Split(xff, ",")[0]
		} else {
			ip = xff
		}
	}
	isTor := status.IsTorNode(ip)
	isExit := status.IsExitNode(ip)

	body := fmt.Sprintf(`<html>
<head>
<title>Tor Status Proxy</title>
<meta charset="utf-8"/>
</head>
<body>
<h2>Your Machine</h2>
<pre>
	Your IP  : %s
	TorNode? : %v
	ExitNode?: %v
</pre>
<h2>UpdatedAt</h2>
<pre>
	AllNodes : %v
	ExitNodea: %v
</pre>
<h2>Example</h2>
<pre>
	$ curl https://tor-status-proxy.herokuapp.com/all
	 => all nodes list

	$ curl https://tor-status-proxy.herokuapp.com/exit
	 => exit nodes list

	$ curl -XPOST -i -d "ip=%s" https://tor-status-proxy.herokuapp.com/istor
	 => If your IP is tor node, the response code will be 200.

	$ curl -XPOST -i -d "ip=%s" https://tor-status-proxy.herokuapp.com/isexit
	 => If your IP is exit node, the response code will be 200.
</pre>
<h2>Source code</h2>
<pre>
	<a href="https://github.com/hiroakis/tor-status-proxy" target="_blank">https://github.com/hiroakis/tor-status-proxy</a>
</pre>
</body>
</html>`, ip, isTor, isExit, status.AllNodeLastModified(), status.ExitNodeLastModified(), ip, ip)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	w.Write([]byte(body))
}

func runServer(listen string) {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/all", allNodesHandler)
	http.HandleFunc("/istor", isTorNodesHandler)
	http.HandleFunc("/exit", exitNodesHandler)
	http.HandleFunc("/isexit", isExitNodeHandler)
	http.ListenAndServe(listen, nil)
}

func main() {

	var (
		p   int
		err error
	)

	envPort := os.Getenv("PORT") // Heroku
	if envPort == "" {
		p = 9000
	} else {
		p, err = strconv.Atoi(envPort)
		if err != nil {
			fmt.Println("Couldn't parse $PORT")
			return
		}
	}

	var (
		host     string
		port     int
		interval int
	)

	flag.StringVar(&host, "h", "0.0.0.0", "The listen IP.")
	flag.IntVar(&port, "p", p, "The listen port.")
	flag.IntVar(&interval, "i", 3600, "The polling interval in sec.")
	flag.Parse()

	status = tor.NewStatus(interval)

	runServer(fmt.Sprintf("%s:%d", host, port))
}
