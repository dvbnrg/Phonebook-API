package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type customer struct {
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
}

var customers []customer

func main() {
	router := mux.NewRouter()
	customers = append(customers, customer{Email: "asdf", Firstname: "as", Lastname: "df", Phone: "1800"})
	customers = append(customers, customer{Email: "qwer", Firstname: "qw", Lastname: "er", Phone: "1900"})
	customers = append(customers, customer{Email: "zxcv", Firstname: "zx", Lastname: "cv", Phone: "2000"})
	router.HandleFunc("/customer", readAll).Methods("GET")
	router.HandleFunc("/customer/{phone}", read).Methods("GET")
	router.HandleFunc("/customer/{phone}", create).Methods("PUT")
	router.HandleFunc("/customer/{phone}", delete).Methods("DELETE")
	router.HandleFunc("/export", dumpcsv).Methods("GET")
	router.HandleFunc("/import", grabcsv).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func create(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var customer customer
	_ = json.NewDecoder(req.Body).Decode(&customer)
	customer.Phone = params["phone"]
	customers = append(customers, customer)
	json.NewEncoder(w).Encode(customers)
}

func read(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range customers {
		if item.Phone == params["phone"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&customer{})
}

func readAll(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(customers)
}

func delete(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range customers {
		if item.Phone == params["phone"] {
			customers = append(customers[:index], customers[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(customers)
}

func dumpcsv(w http.ResponseWriter, req *http.Request) {
	file, err := os.OpenFile("dump.csv", os.O_CREATE|os.O_WRONLY, 0777)
	defer file.Close()

	if err != nil {
		os.Exit(1)
	}
	csvWriter := csv.NewWriter(file)
	strWrite := parseobject(customers)
	csvWriter.WriteAll(strWrite)
	csvWriter.Flush()
	fmt.Println("Data has been dumped")
}

func grabcsv(w http.ResponseWriter, req *http.Request) {
	file, err := os.Open("grab.csv")
	reader := csv.NewReader(bufio.NewReader(file))

	if err != nil {
		os.Exit(1)
	}

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		customers = append(customers, customer{
			Firstname: line[0],
			Lastname:  line[1],
			Phone:     line[2],
			Email:     line[3],
		},
		)
	}
	parseobject(customers)
	fmt.Println("Data has been grabbed")
}

func parseobject(c []customer) [][]string {
	result := make([][]string, 15)
	for i := range c {
		result[i] = make([]string, 4)
		result[i] = append(result[i], c[i].Email)
		result[i] = append(result[i], c[i].Firstname)
		result[i] = append(result[i], c[i].Lastname)
		result[i] = append(result[i], c[i].Phone)
	}
	return result
}
