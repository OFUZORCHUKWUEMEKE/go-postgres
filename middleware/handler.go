package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/OFUZORCHUKWUEMEKE/go-postgres/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {
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
	insertID := insertStock(stock)

	res := response{
		ID:      insertID,
		Message: "stock created Successfully",
	}

	json.NewEncoder(w).Encode(res)
}

func GetStock(w http.ResponseWriter, r *http.Request) {
	// get the stockid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the getStock function with stock id to retrieve a single stock
	stock, err := getStock(int64(id))

	if err != nil {
		log.Fatalf("Unable to get stock. %v", err)
	}

	// send the response
	json.NewEncoder(w).Encode(stock)
}

func GetAllStocks(w http.ResponseWriter, r *http.Request) {

	stocks, err := getAllStocks()

	if err != nil {
		log.Fatalf("Unable to get all stock. %v", err)
	}

	// send all the stocks as response
	json.NewEncoder(w).Encode(stocks)
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// create an empty stock of type models.Stock
	var stock models.Stock

	// decode the json request to stock
	err = json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call update stock to update the stock
	updatedRows := updateStock(int64(id), stock)

	// format the message string
	msg := fmt.Sprintf("Stock updated successfully. Total rows/record affected %v", updatedRows)

	// format the response message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("unable to convert string to int. %v", err)
	}
	deletedRows := deleteStock(int64(id))

	msg := fmt.Sprint("Stock deleted successfully. Total rows/records %v", deletedRows)
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

func insertStock(stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `INSERT into stocks(name,price,company) VALUES ($1, $2, $3) RETURNING stockid`
	var id int64

	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	fmt.Printf("Inserted a single record %v", id)
	return id
}

func getStock(id int64) (models.Stock, error) {
	db := createConnection()
	defer db.Close()

	var stock models.Stock
	sqlstatement := `SELECT * FROM stocks WHERE stockid=$1`
	row := db.QueryRow(sqlstatement, id)

	err := row.Scan(&stock.StockId, &stock.Name, &stock.Price, &stock.Company)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned")
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatalf("unable to scan the row. %v", err)
	}
	return stock, nil
}

func getAllStocks() ([]models.Stock, error) {
	db := createConnection()
	defer db.Close()

	var stocks []models.Stock
	sqlstatement := `SELECT * FROM stocks`

	row, err := db.Query(sqlstatement)
	if err != nil {
		log.Fatalf("Unable to execute the query %v", err)
	}
	defer row.Close()
	for row.Next() {
		var stock models.Stock
		err := row.Scan(&stock.Name, &stock.Company, &stock.Price, &stock.StockId)
		if err != nil {
			log.Fatalf("Unable to scan the row %v", err)
		}
		stocks = append(stocks, stock)
	}
	return stocks, err
}

func updateStock(id int64, stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()

	sqlstatement := `UPDATE stocks SET name=$2 , price=$3 , company=$4 WHERE stockid=$1`

	res, err := db.Exec(sqlstatement, id, stock.Name, stock.Price, stock.Company)

	if err != nil {
		log.Fatalf("Unable to execute the query %v", err)
	}
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows .%v", err)
	}
	fmt.Printf("Total rows/records affected %v", rowsAffected)
	return rowsAffected

}

func deleteStock(id int64) int64 {
	db := createConnection()

	defer db.Close()

	sqlstatement := `DELETE FROM stocks WHERE stockid=$1`
	res, err := db.Exec(sqlstatement, id)
	if err != nil {
		log.Fatalf("Unable to execute the query %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	fmt.Printf("Total rows/records affected %v", rowsAffected)
	return rowsAffected
}
