package main

import (
    "database/sql"
    "log"
    "net/http"
    "text/template"
    "os"

    _ "github.com/go-sql-driver/mysql"
)

type Tool struct {
    Id       int
    Name     string
    Price    string
    Quantity string
    Status   string
}

func dbConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := ""
    dbName := "demo_go"
    db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    return db
}

var tmpl = template.Must(template.ParseGlob("templates/*"))
var err error

//Index handler
func Index(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    selDB, err := db.Query("SELECT * FROM product ORDER BY id DESC")
    if err != nil {
        panic(err.Error())
    }

    tool := Tool{}
    res := []Tool{}

    for selDB.Next() {
        var id int
        var name string
        var price,quantity string
        var status string
        err := selDB.Scan(&id, &name, &price, &quantity, &status)
        if err != nil {
            panic(err.Error())
        }
        log.Println("Listing Row: Id " + string(id) + " | name " + name + " | price " + string(price) + " | quantity " + string(quantity) + " | status " + string(status))

        tool.Id = id
        tool.Name = name
        tool.Price = price
        tool.Quantity = quantity
        tool.Status = status
        res = append(res, tool)
    }
    tmpl.ExecuteTemplate(w, "Index", res)
    defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
    tmpl.ExecuteTemplate(w, "New", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        name := r.FormValue("name")
        price := r.FormValue("price")
        quantity := r.FormValue("quantity")
        status := 1
        insForm, err := db.Prepare("INSERT INTO product (name, price, quantity, status) VALUES (?, ?, ?, ?)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(name, price, quantity, status)
        //log.Println("Insert Data: name " + name + " | price " + price + " | quantity " + quantity + " | status " + status)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
       name := r.FormValue("name")
        price := r.FormValue("price")
        quantity := r.FormValue("quantity")
        id := r.FormValue("id")
        insForm, err := db.Prepare("UPDATE product SET name=?, price=?, quantity=? WHERE id=?")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(name, price, quantity, id)
        //log.Println("UPDATE Data: name " + name + " | category " + category + " | url " + url + " | rating " + rating + " | notes " + notes)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    tool := r.URL.Query().Get("id")
    delForm, err := db.Prepare("DELETE FROM product WHERE id=?")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(tool)
    log.Println("DELETE " + tool)
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Edit(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT * FROM product WHERE id=?", nId)
    if err != nil {
        panic(err.Error())
    }

    tool := Tool{}

    for selDB.Next() {
        var id int
        var name string
        var price,quantity string
        var status string
        err := selDB.Scan(&id, &name, &price, &quantity, &status)
        if err != nil {
            panic(err.Error())
        }

        tool.Id = id
        tool.Name = name
        tool.Price = price
        tool.Quantity = quantity
        tool.Status = status
    }

    tmpl.ExecuteTemplate(w, "Edit", tool)
    defer db.Close()
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

    mux := http.NewServeMux()

    // Add the following two lines
    fs := http.FileServer(http.Dir("assets"))
    mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

    fa := http.FileServer(http.Dir("files"))
    mux.Handle("/files/", http.StripPrefix("/files/", fa))

    mux.HandleFunc("/", Index)
    mux.HandleFunc("/Show", Index)
    mux.HandleFunc("/New", New)
    mux.HandleFunc("/edit", Edit)
    mux.HandleFunc("/update", Update)
    mux.HandleFunc("/delete", Delete)
    mux.HandleFunc("/insert", Insert)

    http.ListenAndServe(":"+port, mux)
}
