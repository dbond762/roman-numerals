package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/chi/middleware"
)

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

	var result interface{}
	if arab, err := strconv.Atoi(number); err == nil {
		// Convert arab to roman.
		if arab <= 0 || arab >= 4000 {
			http.Error(w, "Bad number", 500)
			return
		}

		roman := make([]byte, 0, 15)
		digit := 1000
		for i := 3; i >= 0; i-- {
			d := arab / digit
			roman = append(roman, []byte(table[i][d])...)
			arab %= digit
			digit /= 10
		}

		result = map[string]string {"Res": string(roman)}
	} else if romanRegexp.MatchString(number) {
		// Convert roman to arab.
		digits := romanRegexp.FindAllStringSubmatch(number, 1)

		arab := 0
		for i, digit := 1, 1000; i < len(digits[0]); i, digit = i+1, digit/10 {
			for k, v := range table[4-i] {
				if digits[0][i] == v {
					arab += k * digit
					break
				}
			}
		}

		result = map[string]int {"Res": arab}
	} else {
		http.Error(w, "Bad number", 500)
		return
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(jsonData)
}

func main() {
	r := chi.NewRouter()

	CORS := cors.New(cors.Options{
		AllowedOrigins:		[]string{"*"},
		AllowedMethods:		[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:		[]string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:		[]string{"Link"},
		AllowCredentials:	true,
		MaxAge:				300,
	})
	r.Use(CORS.Handler)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.Get("/convert", Convert)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
