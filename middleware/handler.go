package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"strconv"

	"github.com/OFUZORCHUKWUEMEKE/go-postgres/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func CreateConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to postgres")
	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatal("unable to decode the request body. %v", err)
	}
	insertID := insertStock()

	res := response{
		ID:      insertID,
		Message: "stock created Successfully",
	}

	json.NewEncoder(w).Encode(res)
}

func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := stdconv.Atoi(params["id"])

	if err != nil {
		log.Fatal("unable to convert the string into int. %v", err)
	}
	stock, err := getStock(int64(id))
	if err != nil {
		log.Fatalf("unable to get stock. %v", err)

	}
	json.NewDecoder(w).Encode(stock)
}

func GetAllStock(w http.Response, r *http.Request) {
	stocks,err := getAllStocks()

	id, err = strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int ")
	}
	json.NewDecoder(w).Encode(stocks)
}

func UpdateStock(w http.Response , r *http.Request) {
	params := mux.Vars(r)

	id,err := strconv.Atoi(params["id"])

	if err != nil{
		log.Fatalf("unable to convert the string into int. %v",err)
	}

	var stock models.Stock
	err = json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("Unable to decode request body ")
	}
	updatedRows:=  updateStock(int64(id),stock)
	msg := fmt.Sprintf("stock updated successfully")
	res := response {
		ID:int64(id),
		Message:msg
	}
	json.NewDecoder(w).Encode(res)
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id , err := strconv.ParseInt(params["id"])
	if err !=nil {
		log.Fatalf("unable to convert")
	}
}
