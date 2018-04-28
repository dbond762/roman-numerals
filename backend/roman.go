package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

var (
	table = [][]string{
		{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX"},
		{"", "X", "XX", "XXX", "XL", "L", "LX", "LXX", "LXXX", "XC"},
		{"", "C", "CC", "CCC", "CD", "D", "DC", "DCC", "DCCC", "CM"},
		{"", "M", "MM", "MMM"},
	}

	errBadNumber = errors.New("bad number")
)

// Конвертирует арбские числа в римские.
// Если число <= 0 или > 3999 - возвращает ошибку "bad number".
func arab2roman(arab int) (string, error) {
	const maxNumber = 3999
	if arab <= 0 || arab > maxNumber {
		return "", errBadNumber
	}

	var (
		roman = ""
		digit = 1000
	)
	for i := 3; i >= 0; i-- {
		d := arab / digit
		roman += table[i][d]
		arab %= digit
		digit /= 10
	}

	return roman, nil
}

// Конвертирует римские числа в арабские.
// Если введено неправильное число, вернёт ошибку "bad number".
func roman2arab(roman string) (int, error) {
	re := regexp.MustCompile(`^(M{0,3})(D?C{0,3}|C[DM])(L?X{0,3}|X[LC])(V?I{0,3}|I[VX])$`)

	if !re.MatchString(roman) {
		return 0, errBadNumber
	}

	digits := re.FindAllStringSubmatch(roman, 1)

	arab := 0
	digit := 1000
	for i := 1; i < len(digits[0]); i++ {
		for k, v := range table[4-i] {
			if digits[0][i] == v {
				arab += k * digit
				break
			}
		}
		digit /= 10
	}

	return arab, nil
}

func convert(w http.ResponseWriter, r *http.Request) {
	number := chi.URLParam(r, "number")
	number = strings.ToUpper(number)

	result := make(map[string]interface{})

	if arab, err := strconv.Atoi(number); err == nil {
		// Convert arab to roman.
		roman, err := arab2roman(arab)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result["res"] = roman

	} else if arab, err := roman2arab(number); err == nil {
		// Convert roman to arab.
		result["res"] = arab

	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Println(err.Error())

		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	w.Write(jsonData)
}

func main() {
	r := chi.NewRouter()

	CORS := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(CORS.Handler)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Use(middleware.Logger)

	r.Get("/convert/{number}", convert)

	log.Printf("Server run on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
