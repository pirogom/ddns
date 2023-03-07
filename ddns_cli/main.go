package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type updateRecordData struct {
	Domain string `json:"Domain"`
	IP     string `json:"Ip"`
}

var (
	ddnsHost *string
	domain   *string
	command  *string
	ipaddr   *string
)

func getHttpData(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36")

	client := &http.Client{}
	rep, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		rep.Body.Close()
	}()

	buf, berr := ioutil.ReadAll(rep.Body)

	if berr != nil {
		return nil, berr
	}

	return buf, nil
}

func postJsonData(url string, jbuf []byte) ([]byte, error) {
	buff := bytes.NewBuffer(jbuf)
	resp, err := http.Post(url, "application/json", buff)
	if err != nil {
		return nil, err
	}

	buf, berr := ioutil.ReadAll(resp.Body)

	if berr != nil {
		return nil, berr
	}
	return buf, nil
}

func main() {
	domain = flag.String("name", "", "domain name")
	ddnsHost = flag.String("server", "ddns.example.com:8080", "ddns control server address ( host:port )")
	command = flag.String("cmd", "", "command ( update = update record, delete = delete record )")
	ipaddr = flag.String("ip", "", "IP address for Update ( optional )")

	flag.Parse()

	if *domain == "" {
		fmt.Println("domain name is empty! please input -name argument")
		return
	}

	switch strings.ToUpper(*command) {
	case "UPDATE":
		ud := updateRecordData{}
		ud.Domain = *domain
		ud.IP = *ipaddr

		jbuf, jerr := json.Marshal(&ud)

		if jerr != nil {
			fmt.Println(jerr.Error())
			return
		}

		res, resErr := postJsonData("http://"+*ddnsHost+"/UPDATE", jbuf)

		if resErr != nil {
			fmt.Println(resErr.Error())
			return
		}
		fmt.Println(string(res))
	case "DELETE":
		ud := updateRecordData{}
		ud.Domain = *domain

		jbuf, jerr := json.Marshal(&ud)

		if jerr != nil {
			fmt.Println(jerr.Error())
			return
		}

		res, resErr := postJsonData("http://"+*ddnsHost+"/DELETE", jbuf)

		if resErr != nil {
			fmt.Println(resErr.Error())
			return
		}
		fmt.Println(string(res))
	}
}
