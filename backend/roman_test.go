package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func buildURL(URL, value string) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("number", value)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func getResponse(URL string) (*Response, error) {
	got, err := http.Get(URL)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(got.Body)
	if err != nil {
		return nil, err
	}

	var res Response
	json.Unmarshal(body, &res)
	return &res, nil
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
		{40, "XL"},
		{50, "L"},
		{90, "XC"},
		{100, "C"},
		{400, "CD"},
		{500, "D"},
		{900, "CM"},
		{1000, "M"},
		{3999, "MMMCMXCIX"},
	}
	for _, c := range cases {
		URL, err := buildURL(server.URL, c.roman)
		if err != nil {
			t.Fatal(err)
		}
		got, err := getResponse(URL)
		if err != nil {
			t.Fatal(err)
		}
		i, err := strconv.Atoi(got.Res)
		if err != nil {
			t.Fatal(err)
		}
		if i != c.arab {
			t.Errorf("Query: %s,\tgot %d,\twant %d", c.roman, i, c.arab)
		}

		URL, err = buildURL(server.URL, string(c.arab))
		if err != nil {
			t.Fatal(err)
		}
		got, err = getResponse(URL)
		if err != nil {
			t.Fatal(err)
		}
		if got.Res != c.roman {
			t.Errorf("Query: %d,\tgot %s,\twant %s", c.arab, got.Res, c.roman)
		}
	}
}
