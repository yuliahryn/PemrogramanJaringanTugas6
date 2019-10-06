package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var templates map[string]*template.Template

var ctx = context.Background()

type karyawan struct {
	Id     bson.ObjectId `bson:"_id"`
	Nama   string        `bson:"nama"`
	Email  string        `bson:"email"`
	Notelp string        `bson:"notelp"`
	Alamat string        `bson:"alamat"`
}

func connect() (*mongo.Database, error) {
	clientOptions := options.Client()
	clientOptions.ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client.Database("karyawan_db"), nil
}

func init() {
	loadTemplates()
}

func main() {

	router := mux.NewRouter()

	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/tambah", tambah).Methods("POST")
	router.HandleFunc("/update", update).Methods("POST")
	router.HandleFunc("/hapus", hapus).Methods("POST")

	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func index(res http.ResponseWriter, req *http.Request) {
	db, err := connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	csr, err := db.Collection("karyawan").Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer csr.Close(ctx)

	result := make([]karyawan, 0)
	for csr.Next(ctx) {
		var row karyawan
		err := csr.Decode(&row)
		if err != nil {
			log.Fatal(err.Error())
		}

		result = append(result, row)
	}

	var data = bson.M{"karyawan": result}

	if err := templates["index"].Execute(res, data); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func tambah(res http.ResponseWriter, req *http.Request) {
	var nama = req.FormValue("nama")
	var email = req.FormValue("email")
	var notelp = req.FormValue("notelp")
	var alamat = req.FormValue("alamat")

	db, err := connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = db.Collection("karyawan").InsertOne(ctx, karyawan{bson.NewObjectId(), nama, email, notelp, alamat})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Insert success!")

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func update(res http.ResponseWriter, req *http.Request) {
	var nama_before = req.FormValue("nama-before")
	var email_before = req.FormValue("email-before")
	var notelp_before = req.FormValue("notelp-before")
	var alamat_before = req.FormValue("alamat-before")

	var id = req.FormValue("id")
	var nama = req.FormValue("nama")
	var email = req.FormValue("email")
	var notelp = req.FormValue("notelp")
	var alamat = req.FormValue("alamat")

	db, err := connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	var selector = bson.M{"nama": nama_before, "email": email_before, "notelp": notelp_before, "alamat": alamat_before}
	var changes = karyawan{bson.ObjectIdHex(id), nama, email, notelp, alamat}

	_, err = db.Collection("karyawan").UpdateOne(ctx, selector, bson.M{"$set": changes})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Update success!")

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func hapus(res http.ResponseWriter, req *http.Request) {
	var nama_before = req.FormValue("nama-before")
	var email_before = req.FormValue("email-before")
	var notelp_before = req.FormValue("notelp-before")
	var alamat_before = req.FormValue("alamat-before")

	db, err := connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	var selector = bson.M{"nama": nama_before, "email": email_before, "notelp": notelp_before, "alamat": alamat_before}
	_, err = db.Collection("karyawan").DeleteOne(ctx, selector)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Remove success!")

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func loadTemplates() {
	templates = make(map[string]*template.Template)

	templates["index"] = template.Must(template.ParseFiles("views/index.html"))
}
