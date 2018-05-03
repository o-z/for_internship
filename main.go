package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var price float32

func set_get(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "templates/index.html")

	case "POST":

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		number1, err := strconv.Atoi(r.FormValue("number1")) // convert
		number2, err := strconv.Atoi(r.FormValue("number2")) // convert
		number3, err := strconv.Atoi(r.FormValue("number3")) // convert
		number4, err := strconv.Atoi(r.FormValue("number4")) // convert
		if err != nil {
			log.Fatal(err)
		}
		price = prices(number1, number2, number3, number4)
		fltB, _ := json.Marshal(price)
		fmt.Fprintf(w, string(fltB)) // print

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func prices(number1 int, number2 int, number3 int, number4 int) float32 {
	price = 0
	var full int = 1000*number1 + 100*number2 + 10*number3 + number4
	var prefix int = 100*number2 + 10*number3 + number4
	var two_same_1 int = 10*number1 + number2
	var two_same_2 int = 10*number3 + number4
	important_date := []int{1905, 1903, 1923}
	if number1 == number2 && number2 == number3 && number3 == number4 { //Son 4 rakamı aynı olan numaralardan
		price = price + kdv(1000)
	} else if prefix == 850 || prefix == 885 { // Ön prefix ile simetrik numaralardan
		price = price + kdv(1000)
	} else if number2 == number3 && number3 == number4 { // Son 3 rakamı aynı olan numaralardan
		price = price + kdv(500)
	} else if two_same_1 == two_same_2 { // Son 4 rakamı 2’şer simetrik olanlardan
		price = price + kdv(500)
	} else if number1 == number2 && number3 == number4 { // Son 4 rakamı 2’şer olarak aynı olanlardan
		price = price + kdv(250)
	} else if number2 == 0 && number3 != 0 && number4 == 0 { // Son 3 rakamı 010, 020, 030, 040 gibi artanlardan
		price = price + kdv(100)
	} else if (number2 == number1+1) && (number3 == number2+1) && (number4 == number3+1) { //Son dört rakamı ardışık olanlar
		price = price + kdv(100)
	}
	for _, i := range important_date {
		if i == full {
			price = price + kdv(50)
		}

	}
	return price

}
func kdv(value float32) float32 {
	value = value + value*18/100
	return value
}
func main() {
	http.HandleFunc("/", set_get)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
