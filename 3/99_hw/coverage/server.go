package main

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
	"time"
)

var Read ReadXmL = ReadXmL{}
var UserData []User
var ValidToken = "228"
var XMLLocation = "dataset.xml"
var IsXMLRead = false

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "asdasdaaffg", http.StatusSeeOther)
}

func badJSONStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, "NotParsableJSON")
}

func badJSONHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "NotParsableJSON")
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	handler(w, r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	req := SearchRequest{}
	req.Limit, _ = strconv.Atoi(keys["limit"][0])
	req.Offset, _ = strconv.Atoi(keys["offset"][0])
	req.OrderBy, _ = strconv.Atoi(keys["order_by"][0])
	req.Query = keys["query"][0]
	req.OrderField = keys["order_field"][0]
	// handling errors
	if !IsXMLRead {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if r.Header.Get("AccessToken") != ValidToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if req.OrderField != "" && req.OrderField != "Name" && req.OrderField != "Id" && req.OrderField != "Age" {
		w.WriteHeader(http.StatusBadRequest)
		curError := SearchErrorResponse{
			Error: "ErrorBadOrderField",
		}
		curErrorBytes, _ := json.Marshal(curError)
		io.WriteString(w, string(curErrorBytes))
		return
	}
	if req.OrderBy < -1 || req.OrderBy > 1 {
		w.WriteHeader(http.StatusBadRequest)
		curError := SearchErrorResponse{
			Error: "ErrorBadOrderByField",
		}
		curErrorBytes, _ := json.Marshal(curError)
		io.WriteString(w, string(curErrorBytes))
		return
	}
	filteredUsers := make([]User, 0)
	for _, curUser := range UserData {
		if strings.Contains(curUser.Name, req.Query) || strings.Contains(curUser.About, req.Query) {
			filteredUsers = append(filteredUsers, curUser)
		}
	}

	sort.Slice(filteredUsers, func(i, j int) bool {
		if req.OrderBy == -1 {
			// Descending
			if req.OrderField == "Name" || req.OrderField == "" {
				return filteredUsers[i].Name > filteredUsers[j].Name
			} else if req.OrderField == "Id" {
				return filteredUsers[i].Id > filteredUsers[j].Id
			} else {
				return filteredUsers[i].Age > filteredUsers[j].Age
			}
		} else if req.OrderBy == 1 {
			// Ascending
			if req.OrderField == "Name" || req.OrderField == "" {
				return filteredUsers[i].Name < filteredUsers[j].Name
			} else if req.OrderField == "Id" {
				return filteredUsers[i].Id < filteredUsers[j].Id
			} else {
				return filteredUsers[i].Age < filteredUsers[j].Age
			}
		} else {
			return i < j
		}
	})

	result := []User{}
	for i := req.Offset; i < len(filteredUsers) && i < req.Offset+req.Limit; i++ {
		result = append(result, filteredUsers[i])
	}
	byteRes, _ := json.Marshal(result)
	io.WriteString(w, string(byteRes))
}

type BigUser struct {
	Id        int    `xml:"id"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

type ReadXmL struct {
	List []BigUser `xml:"row"`
}

var tmps *httptest.Server

func startServer(f func(http.ResponseWriter, *http.Request)) string {
	// Parsing xml
	data, err := ioutil.ReadFile(XMLLocation)
	if err != nil {
		IsXMLRead = false
	} else {
		IsXMLRead = true
	}
	xml.Unmarshal(data, &Read)
	for _, cur := range Read.List {
		UserData = append(UserData, User{
			Id:     cur.Id,
			Name:   cur.FirstName + cur.LastName,
			Age:    cur.Age,
			Gender: cur.Gender,
			About:  cur.About,
		})
	}

	tmps = httptest.NewServer(http.HandlerFunc(f))
	return tmps.URL
}

func closeServer() {
	tmps.Close()
}
