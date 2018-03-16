package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getResponse(URL string) ([]byte, error) {
	got, err := http.Get(URL)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(got.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func TestConvert(t *testing.T) {
	handler := http.HandlerFunc(Convert)
	server := httptest.NewServer(handler)
	defer server.Close()

	cases := []struct {
		arab  int
		roman string
	}{
		{1, "I"},
		{2, "II"},
		{3, "III"},
		{4, "IV"},
		{5, "V"},
		{6, "VI"},
		{7, "VII"},
		{8, "VIII"},
		{9, "IX"},
		{10, "X"},
		{20, "XX"},
		{30, "XXX"},
		{40, "XL"},
		{50, "L"},
		{90, "XC"},
		{100, "C"},
		{400, "CD"},
		{500, "D"},
		{900, "CM"},
		{1000, "M"},
		{3999, "MMMCMXCIX"},
		{76, "LXXVI"},
		{94, "XCIV"},
		{99, "XCIX"},
		{283, "CCLXXXIII"},
		{499, "CDXCIX"},
		{999, "CMXCIX"},
		{1950, "MCML"},
	}
	for _, c := range cases {
		URL := fmt.Sprintf("%s?number=%s", server.URL, c.roman)
		got, err := getResponse(URL)
		if err != nil {
			t.Fatal(err)
		}
		var resArab map[string]int
		if err := json.Unmarshal(got, &resArab); err != nil {
			t.Fatal(err)
		}
		arab, ok := resArab["Res"]
		if !ok {
			t.Fatal("Don`t have number in response")
		}
		if arab != c.arab {
			t.Errorf("Query: %s,\tgot %d,\twant %d", URL, arab, c.arab)
		}

		URL = fmt.Sprintf("%s?number=%d", server.URL, c.arab)
		got, err = getResponse(URL)
		if err != nil {
			t.Fatal(err)
		}
		var resRoman map[string]string
		if err := json.Unmarshal(got, &resRoman); err != nil {
			t.Fatal(err)
		}
		roman, ok := resRoman["Res"]
		if !ok {
			t.Fatal("Don`t have number in response")
		}
		if roman != c.roman {
			t.Errorf("Query: %s,\tgot %s,\twant %s", URL, roman, c.roman)
		}
	}
}
