package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/miekg/dns"
)

var (
	port     *int
	wwwport  *int
	dnsMap   map[string]string
	dnsMutex sync.Mutex
)

func saveRecord() error {
	f, ferr := os.OpenFile(".\\ddns.dat", os.O_CREATE|os.O_WRONLY, 0644)

	if ferr != nil {
		return ferr
	}

	var rcnt int

	defer func() {
		f.Close()

		if rcnt == 0 {
			os.Remove(".\\ddns.dat")
		}
	}()

	for _, v := range dnsMap {
		f.Write([]byte(v + "\n"))
		rcnt++
	}
	return nil
}

func loadRecord() error {
	file, err := os.Open(".\\ddns.dat")

	if err != nil {
		return err
	}

	defer file.Close()

	scan := bufio.NewScanner(file)

	for scan.Scan() {

		rd := scan.Text()

		if rd != "" {
			rr, rrerr := dns.NewRR(rd)

			if rrerr == nil {
				key, kerr := getKey(rr.(dns.RR).Header().Name, 1)
				if kerr == nil {
					dnsMap[key] = rd
				}
			}
		}
	}

	err = scan.Err()

	if err != nil {
		fmt.Println(err.Error())
	}
	return nil
}

func getKey(domain string, rtype uint16) (r string, e error) {
	if n, ok := dns.IsDomainName(domain); ok {
		labels := dns.SplitDomainName(domain)

		// Reverse domain, starting from top-level domain
		// eg.  ".com.mkaczanowski.test "
		var tmp string
		for i := 0; i < int(math.Floor(float64(n/2))); i++ {
			tmp = labels[i]
			labels[i] = labels[n-1]
			labels[n-1] = tmp
		}

		reverse_domain := strings.Join(labels, ".")
		r = strings.Join([]string{reverse_domain, strconv.Itoa(int(rtype))}, "_")
	} else {
		e = errors.New("Invailid domain: " + domain)
		fmt.Println(e.Error())
	}

	return r, e
}

func deleteRecord(domain string, rtype uint16) (err error) {

	dnsMutex.Lock()
	defer dnsMutex.Unlock()

	key, kerr := getKey(domain, rtype)

	if kerr != nil {
		return kerr
	}

	_, exists := dnsMap[key]

	if exists {
		delete(dnsMap, key)
	} else {
		e := errors.New("Delete record failed for domain:  " + domain)
		fmt.Println(e.Error())
		return e
	}

	fmt.Println("Delete Record", "-", domain)

	saveRecord()

	return nil
}

func updateRecord(domain string, ipaddr string) (err error) {

	rr := new(dns.A)

	rr.A = net.ParseIP(ipaddr)
	rr.Hdr.Name = domain
	rr.Hdr.Class = dns.ClassINET
	rr.Hdr.Rrtype = 1 // A
	rr.Hdr.Ttl = 30

	err = storeRecord(rr)

	fmt.Println("Update Record", "-", domain, ipaddr)

	return err
}

func storeRecord(rr dns.RR) (err error) {
	dnsMutex.Lock()
	defer dnsMutex.Unlock()

	key, kerr := getKey(rr.Header().Name, rr.Header().Rrtype)

	if kerr != nil {
		return kerr
	}

	dnsMap[key] = rr.String()

	saveRecord()

	return nil
}

func getRecord(domain string, rtype uint16) (rr dns.RR, err error) {

	key, kerr := getKey(domain, rtype)

	if kerr != nil {
		return nil, kerr
	}

	v, exists := dnsMap[key]

	if exists {
		if v == "" {
			e := errors.New("Record not found, key:  " + key)
			fmt.Println(e.Error())

			return nil, e
		}

		rr, err = dns.NewRR(v)

		if err != nil {
			return nil, err
		}

		return rr, nil
	} else {
		e := errors.New("Record not found, key:  " + key)
		fmt.Println(e.Error())

		return nil, e
	}
}

func newRecordA(domain, ipaddr string) {

	rr := new(dns.A)

	rr.A = net.ParseIP(ipaddr)
	rr.Hdr.Name = domain
	rr.Hdr.Class = dns.ClassINET
	rr.Hdr.Rrtype = 1 // A
	rr.Hdr.Ttl = 30

	storeRecord(rr)
}

func parseQuery(m *dns.Msg) {
	var rr dns.RR

	for _, q := range m.Question {
		if read_rr, e := getRecord(q.Name, q.Qtype); e == nil {
			rr = read_rr.(dns.RR)
			if rr.Header().Name == q.Name {
				m.Answer = append(m.Answer, rr)
			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}
	w.WriteMsg(m)
}

func serve(port int) {
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}

	err := server.ListenAndServe()
	defer server.Shutdown()

	if err != nil {
		fmt.Println("Failed to setup the udp server:", err.Error())
	}
}

func main() {
	dnsMap = make(map[string]string)
	loadRecord()

	// Parse flags
	port = flag.Int("port", 53, "server port (dns server)")
	wwwport = flag.Int("cport", 8080, "control port (httpd)")

	flag.Parse()

	// Attach request handler func
	dns.HandleFunc(".", handleDnsRequest)

	go wwwServ(*wwwport)

	// Start server
	serve(*port)
}
