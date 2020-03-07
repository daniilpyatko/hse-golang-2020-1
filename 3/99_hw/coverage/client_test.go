package main

import (
	"fmt"
	"testing"
)

func TestToken(*testing.T) {
	Address := startServer(handler)
	var curClient = SearchClient{
		AccessToken: "NotValidToken",
		URL:         Address,
	}
	Tests := []SearchRequest{
		SearchRequest{
			Limit:      5,
			Offset:     0,
			Query:      "accccc",
			OrderField: "Id",
			OrderBy:    0,
		},
	}
	for i, curTest := range Tests {
		res, err := curClient.FindUsers(curTest)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(i, len(res.Users))
		}
	}
	closeServer()
}

func TestSlowHandler(*testing.T) {
	Address := startServer(slowHandler)
	var curClient = SearchClient{
		AccessToken: ValidToken,
		URL:         Address,
	}
	Tests := []SearchRequest{
		SearchRequest{
			Limit:      5,
			Offset:     0,
			Query:      "",
			OrderField: "Id",
			OrderBy:    0,
		},
	}
	for i, curTest := range Tests {
		res, err := curClient.FindUsers(curTest)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(i, len(res.Users))
		}
	}
	closeServer()
}

func TestServerSorting(*testing.T) {
	Address := startServer(handler)
	var curClient = SearchClient{
		AccessToken: ValidToken,
		URL:         Address,
	}
	Tests := []SearchRequest{
		SearchRequest{
			Limit:      5,
			Offset:     0,
			Query:      "",
			OrderField: "Name",
			OrderBy:    -1,
		},
		SearchRequest{
			Limit:      5,
			Offset:     0,
			Query:      "",
			OrderField: "Id",
			OrderBy:    -1,
		},
		SearchRequest{
			Limit:      5,
			Offset:     0,
			Query:      "",
			OrderField: "Age",
			OrderBy:    -1,
		},
		SearchRequest{
			Limit:      5,
			Offset:     0,
			Query:      "",
			OrderField: "Name",
			OrderBy:    1,
		},
		SearchRequest{
			Limit:      5,
			Offset:     0,
			Query:      "",
			OrderField: "Id",
			OrderBy:    1,
		},
		SearchRequest{
			Limit:      5,
			Offset:     0,
			Query:      "",
			OrderField: "Age",
			OrderBy:    1,
		},
		SearchRequest{
			Limit:      5,
			Offset:     0,
			Query:      "",
			OrderField: "Age",
			OrderBy:    0,
		},
		SearchRequest{
			Limit:      5,
			Offset:     0,
			Query:      "",
			OrderField: "NotProperOrderField",
			OrderBy:    1,
		},
		SearchRequest{
			Limit:      5,
			Offset:     0,
			Query:      "",
			OrderField: "NotProperOrderField",
			OrderBy:    -1,
		},

		SearchRequest{
			Limit:      5,
			Offset:     0,
			Query:      "",
			OrderField: "",
			OrderBy:    228,
		},
		SearchRequest{
			Limit:      5,
			Offset:     10000,
			Query:      "",
			OrderField: "",
			OrderBy:    -1,
		},
	}
	for i, curTest := range Tests {
		res, err := curClient.FindUsers(curTest)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(i, len(res.Users))
		}
	}
	closeServer()
}

func TestGeneral(*testing.T) {
	Address := startServer(handler)
	var curClient = SearchClient{
		AccessToken: ValidToken,
		URL:         Address,
	}
	Tests := []SearchRequest{
		SearchRequest{
			Limit:      25,
			Offset:     0,
			Query:      "Boyd",
			OrderField: "",
			OrderBy:    0,
		},
		SearchRequest{
			Limit:      5,
			Offset:     -1,
			Query:      "accccc",
			OrderField: "Id",
			OrderBy:    0,
		},
		SearchRequest{
			Limit:      5,
			Offset:     0,
			Query:      "accccc",
			OrderField: "Id",
			OrderBy:    0,
		},
		SearchRequest{
			Limit:      -1,
			Offset:     0,
			Query:      "accccc",
			OrderField: "Id",
			OrderBy:    0,
		},
		SearchRequest{
			Limit:      26,
			Offset:     0,
			Query:      "accccc",
			OrderField: "Id",
			OrderBy:    0,
		},
		SearchRequest{
			Limit:      26,
			Offset:     0,
			Query:      "accccc",
			OrderField: "Id",
			OrderBy:    -1,
		},
		SearchRequest{
			Limit:      26,
			Offset:     0,
			Query:      "accccc",
			OrderField: "Id",
			OrderBy:    1,
		},
	}
	for i, curTest := range Tests {
		res, err := curClient.FindUsers(curTest)
		fmt.Println(i)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(len(res.Users))
		}
	}
	closeServer()
}
