package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	db "Barberia/db/gen"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	port := ":8080"

	http.HandleFunc("/", handler) //Estas en la raiz --> Ejecuta la funcion handler root}
	http.HandleFunc("/about", handlerAbout)
	http.HandleFunc("/formsPost", handlerFormsPost)
	http.HandleFunc("/usuarios", handlerListUsers)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("El servidor esta corriendo en el puesto 8080")
	fmt.Printf("Run server: http://localhost%s\n", port)
	http.ListenAndServe(port, nil)
}

func handler(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(rw, r)
		return
	} //No encuentra la direccion ejecuta el error 404

	http.ServeFile(rw, r, "templates/index.html") //Se sirve

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")

}

func handlerAbout(rw http.ResponseWriter, r *http.Request) {
	http.ServeFile(rw, r, "templates/about.html")
}

type Turno struct {
	Nombre   string
	Telefono string
	Email    string
	Servicio string
	Barbero  string
	Fecha    string
	Hora     string
	Notas    string
	Acepta   string
}

type User struct {
	ID       int32
	Nombre   string
	Telefono string
	Email    string
}

func handlerFormsPost(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	turno := Turno{
		Nombre:   r.FormValue("nombre"),
		Telefono: r.FormValue("telefono"),
		Email:    r.FormValue("email"),
		Servicio: r.FormValue("servicio"),
		Barbero:  r.FormValue("barbero"),
		Fecha:    r.FormValue("fecha"),
		Hora:     r.FormValue("hora"),
		Notas:    r.FormValue("notas"),
		Acepta:   r.FormValue("acepta_politicas"),
	}

	conn := conexionBD()
	queries := db.New(conn)

	err := queries.CreateUser(r.Context(), db.CreateUserParams{
		Nombre:   turno.Nombre,
		Telefono: turno.Telefono,
		Email:    turno.Email,
	})
	if err != nil {
		http.Error(rw, "Error insertando en BD", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/confirmacion.html")
	if err != nil {
		http.Error(rw, "Error cargando plantilla", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(rw, turno)
}

func conexionBD() *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:34460)/%s", "root", "", "Barberia")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	return db
}

func handlerListUsers(rw http.ResponseWriter, r *http.Request) {
	conn := conexionBD()
	queries := db.New(conn)

	users, err := queries.ListUsers(r.Context())
	if err != nil {
		http.Error(rw, "Error obteniendo usuarios: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/lista.html")
	if err != nil {
		http.Error(rw, "Error cargando plantilla: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(rw, users)
}
