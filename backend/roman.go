package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Response struct {
	Res string
}

func Convert(w http.ResponseWriter, r *http.Request) {
	numbers, ok := r.URL.Query()["number"]
	if !ok || len(numbers) != 1 {
		http.Error(w, "Bad number", 500)
		return
	}
	number := numbers[0]

	table := [][]string{
		{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX"},
		{"", "X", "XX", "XXX", "XL", "L", "LX", "LXX", "LXXX", "XC"},
		{"", "C", "CC", "CCC", "CD", "D", "DC", "DCC", "DCCC", "CM"},
		{"", "M", "MM", "MMM"},
	}

	romanRegexp := regexp.MustCompile(`^(M{0,3})(D?C{0,3}|C[DM])(L?X{0,3}|X[LC])(V?I{0,3}|I[VX])$`)

	var result []byte
	if arab, err := strconv.Atoi(number); err == nil {
		if arab <= 0 || arab >= 4000 {
			http.Error(w, "Bad number", 500)
			return
		}
		digit := 1000
		for i := 3; i >= 0; i-- {
			d := arab / digit
			result = append(result, []byte(table[i][d])...)
			arab %= digit
			digit /= 10
		}
	} else if romanRegexp.MatchString(number) {
		result = []byte("1")
	} else {
		http.Error(w, "Bad number", 500)
		return
	}

	res := Response{string(result)}
	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	http.HandleFunc("/convert", Convert)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
