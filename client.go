package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"context"
	"os"
    "time"
	"google.golang.org/grpc"
	pb "github.com/BernyCR96/sopes1p2"
)

//Caso struct
type Caso struct {
	Nombre       string `json:"Nombre"`
	Departamento string `json:"Departamento"`
	Edad         int    `json:"Edad"`
	Contagio     string `json:"Forma de Contagio"`
	Estado       string `json:"Estado"`
}

const (
	address     = "0.0.0.0:50051"
	defaultName = "world"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/data", EnviarData).Methods("POST")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

//EnviarData envia la data
func EnviarData(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	var p Caso

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Do something with the Person struct...
	fmt.Println(p.Nombre)

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
			log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCalculatorClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
			name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SquareRoot(ctx, &pb.Number{data: p.Nombre})
	if err != nil {
			log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Dato1 string `json:"mensaje"`
	}{Dato1: "Llego super bien"})
}

