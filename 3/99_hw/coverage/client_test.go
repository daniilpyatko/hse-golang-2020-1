package main

import (
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
	for _, curTest := range Tests {
		curClient.FindUsers(curTest)
	}
	closeServer()
}

func TestUnknownError(*testing.T) {
	Address := startServer(redirectHandler)
	var curClient = SearchClient{
		AccessToken: ValidToken,
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
	for _, curTest := range Tests {
		curClient.FindUsers(curTest)
	}
	closeServer()
}

func TestBadStatusJSONError(*testing.T) {
	Address := startServer(badJSONStatusHandler)
	var curClient = SearchClient{
		AccessToken: ValidToken,
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
	for _, curTest := range Tests {
		curClient.FindUsers(curTest)
	}
	closeServer()
}
func TestXMLLocation(*testing.T) {
	XMLLocation = "NotValidXMLLocation"
	Address := startServer(handler)
	var curClient = SearchClient{
		AccessToken: ValidToken,
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
	for _, curTest := range Tests {
		curClient.FindUsers(curTest)
	}
	closeServer()
	XMLLocation = "dataset.xml"
}
func TestNonValidXML(*testing.T) {
	XMLLocation = "NotValidDataset.xml"
	Address := startServer(handler)
	var curClient = SearchClient{
		AccessToken: ValidToken,
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
	for _, curTest := range Tests {
		curClient.FindUsers(curTest)
	}
	closeServer()
	XMLLocation = "dataset.xml"
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
	for _, curTest := range Tests {
		curClient.FindUsers(curTest)
	}
	closeServer()
}

func TestBadJSONHandler(*testing.T) {
	Address := startServer(badJSONHandler)
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
	for _, curTest := range Tests {
		curClient.FindUsers(curTest)
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
	for _, curTest := range Tests {
		curClient.FindUsers(curTest)
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
	for _, curTest := range Tests {
		curClient.FindUsers(curTest)
	}
	closeServer()
}
