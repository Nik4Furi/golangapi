package main

import (
	"encoding/json"
	"log"
	"net/http"
	"../../gorilla/mux"
	"fmt"
	//Mongo db setup
	"context"
         "go.mongodb.org/mongo-driver/mongo"
         "go.mongodb.org/mongo-driver/mongo/options"
		 "go.mongodb.org/mongo-driver/bson"
)



//----------- Create book modal
type Book struct {
	ID string  `json:"id"`
	Title string `json:"title"`
}

var books []Book


//----------------- Creaet the controllers to handle routes
func getBooks(w http.ResponseWriter,r *http.Request){
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println("Failed to connect to MongoDB:", err)
		return
	}

	collection := client.Database("booksapi").Collection("books")

	cursor, err := collection.Find(context.TODO(),bson.M{})
	if err != nil {
		fmt.Println("Failed to fetch data:", err)
		return
	}
	defer cursor.Close(context.TODO())
	
	for cursor.Next(context.TODO()) {
		var book Book
		err := cursor.Decode(&book)
		if err != nil {
			fmt.Println("Failed to decode data:", err)
			return
		}
		books = append(books, book)
	}

	if err := cursor.Err(); err != nil {
		fmt.Println("Cursor error:", err)
		return
	}

	// Process the fetched books
	for _, book := range books {
		fmt.Println("ID:", book.ID)
		fmt.Println("Title:", book.Title)
		fmt.Println("--------------------")
	}

}

// func getBook(w http.ResponseWriter,r *http.Request){
// 	params := mux.Vars(r)
// 	bookID := params["id"]

// 	filter := bson.M{"_id": bookID}

// 	result, err := collection.FindOne(context.TODO(), filter)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// fmt.Fprintf(w, "Deleted %d document(s)", result)
// 	// Process the fetched books
// 	// for _, book := range books {
// 	// 	fmt.Println("ID:", book.ID)
// 	// 	fmt.Println("Title:", book.Title)
// 	// 	fmt.Println("--------------------")
// 	// }
// 	fmt.Fprint(w,result);
// }

func createBook(w http.ResponseWriter,r *http.Request){
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB", http.StatusInternalServerError)
		return
	}
	collection := client.Database("booksapi").Collection("books")

	_, err = collection.InsertOne(context.TODO(), book)
	if err != nil {
		http.Error(w, "Failed to insert document", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Book added successfully!")
}

// var collection *mongo.Collection

func deleteBook(w http.ResponseWriter,r *http.Request){
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println("Failed to connect to MongoDB:", err)
		return
	}

	collection := client.Database("booksapi").Collection("books")

	params := mux.Vars(r)
	bookID := params["id"]

	filter := bson.M{"_id": bookID}

	fmt.Println("value of filter ",filter);

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("the data is ",result)

	fmt.Fprintf(w, "Deleted %v document(s)", result.DeletedCount)
}

func updateBook(w http.ResponseWriter,r *http.Request){
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println("Failed to connect to MongoDB:", err)
		return
	}

	collection := client.Database("booksapi").Collection("books")

	params := mux.Vars(r)
	bookID := params["id"]

	fmt.Println("check book id is ",bookID);

	
	var updatedBook Book
	_ = json.NewDecoder(r.Body).Decode(&updatedBook)
	
	
	filter := bson.M{"_id": bookID}
	update := bson.M{"$set": bson.M{
		"id":       updatedBook.ID,
		"title":       updatedBook.Title,
	}}
	var book Book
	content := collection.FindOne(context.TODO(),filter).Decode(&book)
	json.NewEncoder(w).Encode(book)
	fmt.Println("check content ",content);

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "Modified %d document(s)", result.ModifiedCount)
}


func main(){
	



	// fmt.Println("hello world")

	//------------- initialzing the routers
	r := mux.NewRouter()

	//---------- Handle the routers
	r.HandleFunc("/api/books",getBooks).Methods("GET")
	// r.HandleFunc("/api/book/{id}",getBook).Methods("GET")
	r.HandleFunc("/api/book",createBook).Methods("POST")
	r.HandleFunc("/api/book/{id}",updateBook).Methods("PUT")
	r.HandleFunc("/api/book/{id}",deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000",r))


}