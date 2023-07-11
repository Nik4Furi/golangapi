package main

import (
	"encoding/json"
	"log"
	"net/http"
	"../../gorilla/mux"
)


//----------- Create book modal
type Book struct {
	ID string  `json:"id"`
	Title string `json:"title"`
}

var books []Book


//----------------- Creaet the controllers to handle routes
func getBooks(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type" , "application/json")

	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type" , "application/json")

	params := mux.Vars(r) // params data

	for _,item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}  
	}

	json.NewEncoder(w).Encode(&Book{})

}

func createBook(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type" , "application/json")

	var book Book = json.NewDecoder(r.body).Decode(&book)

	params := mux.Vars(r) // params data

	for _,item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}  
	}

	json.NewEncoder(w).Encode(&Book{})

}

func main(){
	// fmt.Println("hello world")

	//------------- initialzing the routers
	r := mux.NewRouter()

	books = append(books,Book{ID : "1",Title:"Book -1"})
	books = append(books,Book{ID : "2",Title:"Book -2"})

	//---------- Handle the routers
	r.HandleFunc("/api/books",getBooks).Methods("GET")
	r.HandleFunc("/api/book/{id}",getBook).Methods("GET")
	// r.HandleFunc("/api/book",createBook).Methods("POST")
	// r.HandleFunc("/api/book/{id}",updateBook).Methods("PUT")
	// r.HandleFunc("/api/book/{id}",deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000",r))


}