package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

type updateRecordData struct {
	Domain string `json:"Domain"`
	IP     string `json:"Ip"`
}

func getConnIP(conn net.Conn) string {
	rawIP := conn.RemoteAddr().String()

	return splitRemoteAddr(rawIP)
}

func splitRemoteAddr(rawIP string) string {
	if strings.Index(rawIP, ":") != -1 {
		sip := strings.Split(rawIP, ":")
		if len(sip) != 2 {
			return rawIP
		}
		return sip[0]
	}
	return rawIP
}

func forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(403)
	w.Write([]byte(http.StatusText(403)))
}

func notFoundResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte(http.StatusText(404)))
}

func updateRecordRequest(w http.ResponseWriter, r *http.Request) {
	raddr := splitRemoteAddr(r.RemoteAddr)

	od := updateRecordData{}

	err := json.NewDecoder(r.Body).Decode(&od)

	if err != nil {
		forbiddenResponse(w, r)
		return
	}

	if od.IP == "" {
		if err := updateRecord(od.Domain, raddr); err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	} else {
		if err := updateRecord(od.Domain, od.IP); err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	}
	w.Write([]byte("UPDATE RECORD SUCCESS"))
}

func deleteRecordRequest(w http.ResponseWriter, r *http.Request) {
	od := updateRecordData{}

	err := json.NewDecoder(r.Body).Decode(&od)

	if err != nil {
		forbiddenResponse(w, r)
		return
	}

	derr := deleteRecord(od.Domain, 1)

	if derr != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("DELETE RECORD SUCCESS"))
}

/**
*	webPageProc
**/
func webPageProc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-control", "no-cache")

	// Parsing From Data
	r.ParseForm()
	defer func() {
		r.Body.Close()
	}()
	//

	switch r.URL.Path {
	case "/UPDATE":
		updateRecordRequest(w, r)
	case "/DELETE":
		deleteRecordRequest(w, r)
	default:
		forbiddenResponse(w, r)
	}
}

/**
*	wwwServ
**/
func wwwServ(servPort int) {
	http.HandleFunc("/", webPageProc)

	//
	servAddr := fmt.Sprintf(":%d", servPort)

	err := http.ListenAndServe(servAddr, nil)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}
