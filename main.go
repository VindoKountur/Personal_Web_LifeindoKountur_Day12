package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personalweb/connection"
	utils "personalweb/lib"
	"personalweb/middleware"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	router := mux.NewRouter()

	connection.DatabaseConnect()

	router.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	router.PathPrefix("/uploads").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/contactme", contactMeHandler)
	router.HandleFunc("/myproject", myProjectHandler)
	router.HandleFunc("/detail/{id}", detailProjectHandler)
	router.HandleFunc("/editproject/{id}", editProjectHandler)
	router.HandleFunc("/addproject", middleware.UploadFile(addProjectHandler)).Methods("POST")
	router.HandleFunc("/updateproject", middleware.UploadFile(updateProjectHandler)).Methods("POST")
	router.HandleFunc("/deleteproject/{id}", deleteProjectHandler)
	router.HandleFunc("/login", formLoginHandler).Methods("GET")
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/register", formRegisterHandler).Methods("GET")
	router.HandleFunc("/register", registerHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler)

	var address = "localhost:5000"
	fmt.Printf("Server started at %s", address)
	err := http.ListenAndServe(address, router)
	if err != nil {
		fmt.Println(err.Error())
	}
}

var Data = map[string]interface{}{
	"Title":   "Personal Web",
	"IsLogin": false,
}

type Project struct {
	Id              int
	Name            string
	StartDate       time.Time
	EndDate         time.Time
	StartDateFormat string
	EndDateFormat   string
	Description     string
	Tech            []string
	Image           string
	Duration        string
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

var store = sessions.NewCookieStore([]byte("SESSION_ID"))

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message + " + err.Error()))
		return
	}

	session, _ := store.Get(r, "SESSION_ID")

	var rows pgx.Rows
	if session.Values["IsLogin"] != true {
		Data["IsLogin"] = false

		strQuery := "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects ORDER BY id DESC"
		rows, err = connection.Conn.Query(context.Background(), strQuery)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Message + " + err.Error()))
			return
		}
	} else {
		Data["IsLogin"] = true
		Data["Name"] = session.Values["Name"]
		user_id := session.Values["UserID"]

		strQuery := "SELECT tb_projects.id, tb_projects.name, start_date, end_date, description, technologies, image FROM tb_projects INNER JOIN tb_users ON tb_users.id = tb_projects.user_id WHERE user_id = $1 ORDER BY id DESC"

		rows, err = connection.Conn.Query(context.Background(), strQuery, user_id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Message + " + err.Error()))
			return
		}
	}

	var result []Project

	for rows.Next() {
		var row = Project{}

		err := rows.Scan(&row.Id, &row.Name, &row.StartDate, &row.EndDate, &row.Description, &row.Tech, &row.Image)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		row.Duration = utils.CountDuration(row.StartDate, row.EndDate)
		result = append(result, row)
	}
	resp := map[string]interface{}{
		"Data":     Data,
		"Projects": result,
	}

	w.WriteHeader(http.StatusOK)
	err = tmpl.Execute(w, resp)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func contactMeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contactme.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message + " + err.Error()))
		return
	}

	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data["IsLogin"] = false
	} else {
		Data["IsLogin"] = true
		Data["Name"] = session.Values["Name"]
	}

	tmpl.Execute(w, Data)
}

func myProjectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	session, _ := store.Get(r, "SESSION_ID")

	var tmpl, err = template.ParseFiles("views/myproject.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message + " + err.Error()))
		return
	}

	if session.Values["IsLogin"] != true {
		Data["IsLogin"] = false
	} else {
		Data["IsLogin"] = true
		Data["Name"] = session.Values["Name"]
	}

	err = tmpl.Execute(w, Data)
}
func detailProjectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	tmpl, err := template.ParseFiles("views/detailproject.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message + " + err.Error()))
		return
	}

	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data["IsLogin"] = false
	} else {
		Data["IsLogin"] = true
		Data["Name"] = session.Values["Name"]
	}

	ProjectDetail := Project{}

	sqlQuery := "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects WHERE id = $1"

	err = connection.Conn.QueryRow(context.Background(), sqlQuery, id).Scan(&ProjectDetail.Id, &ProjectDetail.Name, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Description, &ProjectDetail.Tech, &ProjectDetail.Image)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message + " + err.Error()))
		return
	}

	ProjectDetail.StartDateFormat = utils.GetDateFormat(ProjectDetail.StartDate)
	ProjectDetail.EndDateFormat = utils.GetDateFormat(ProjectDetail.EndDate)
	ProjectDetail.Duration = utils.CountDuration(ProjectDetail.StartDate, ProjectDetail.EndDate)

	resp := map[string]interface{}{
		"Data":    Data,
		"Project": ProjectDetail,
	}
	err = tmpl.Execute(w, resp)
}

func addProjectHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "SESSION_ID")

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	projectName := r.PostForm.Get("projectName")
	startDateStr := r.PostForm.Get("startDate") + " 00:00:00"
	endDateStr := r.PostForm.Get("endDate") + " 00:00:00"
	format := "2006-01-02 15:04:05"
	startDate, _ := time.Parse(format, startDateStr)
	endDate, _ := time.Parse(format, endDateStr)
	description := r.PostForm.Get("description")
	node := r.PostForm.Get("node")
	vuejs := r.PostForm.Get("vuejs")
	react := r.PostForm.Get("react")
	php := r.PostForm.Get("php")
	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)

	var techList = []string{}

	if node != "" {
		techList = append(techList, node)
	}
	if vuejs != "" {
		techList = append(techList, vuejs)
	}
	if react != "" {
		techList = append(techList, react)
	}
	if php != "" {
		techList = append(techList, php)
	}
	var newProject = Project{
		Name:        projectName,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: description,
		Tech:        techList,
		Image:       image,
	}

	if session.Values["IsLogin"] != true {
		Data["IsLogin"] = false
	} else {
		Data["IsLogin"] = true
		Data["Name"] = session.Values["Name"]
	}

	sqlQuery := "INSERT INTO tb_projects (name, start_date, end_date, description, technologies, image, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err = connection.Conn.Exec(context.Background(), sqlQuery, newProject.Name, newProject.StartDate, newProject.EndDate, newProject.Description, newProject.Tech, newProject.Image, session.Values["UserID"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message" + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func deleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	sqlQuery := "DELETE FROM tb_projects WHERE id = $1"
	_, err := connection.Conn.Exec(context.Background(), sqlQuery, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func editProjectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	editProject := Project{}

	type TechUsed struct {
		Node  bool
		Vuejs bool
		React bool
		Php   bool
	}

	var tmpl, err = template.ParseFiles("views/editproject.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message + " + err.Error()))
		return
	}

	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data["IsLogin"] = false
	} else {
		Data["IsLogin"] = true
		Data["Name"] = session.Values["Name"]
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	sqlQuery := "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects WHERE id = $1"
	err = connection.Conn.QueryRow(context.Background(), sqlQuery, id).Scan(&editProject.Id, &editProject.Name, &editProject.StartDate, &editProject.EndDate, &editProject.Description, &editProject.Tech, &editProject.Image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message + " + err.Error()))
		return
	}

	editProject.StartDateFormat = utils.InputHtmlDateFormat(editProject.StartDate)
	editProject.EndDateFormat = utils.InputHtmlDateFormat(editProject.EndDate)

	techUsed := TechUsed{
		// nodejs: stringInSlice(editProject.Tech, "node"),
		// vuejs: ,
		Node:  utils.StringInSlice(editProject.Tech, "node"),
		Vuejs: utils.StringInSlice(editProject.Tech, "vuejs"),
		React: utils.StringInSlice(editProject.Tech, "react"),
		Php:   utils.StringInSlice(editProject.Tech, "php"),
	}

	resp := map[string]interface{}{
		"Data":        Data,
		"EditProject": editProject,
		"TechUsed":    techUsed,
	}
	err = tmpl.Execute(w, resp)
}

func updateProjectHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	id := r.PostForm.Get("id")
	projectName := r.PostForm.Get("projectName")
	startDateStr := r.PostForm.Get("startDate") + " 00:00:00"
	endDateStr := r.PostForm.Get("endDate") + " 00:00:00"
	format := "2006-01-02 15:04:05"
	startDate, _ := time.Parse(format, startDateStr)
	endDate, _ := time.Parse(format, endDateStr)
	description := r.PostForm.Get("description")
	node := r.PostForm.Get("node")
	vuejs := r.PostForm.Get("vuejs")
	react := r.PostForm.Get("react")
	php := r.PostForm.Get("php")
	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)

	var techList = []string{}

	if node != "" {
		techList = append(techList, node)
	}
	if vuejs != "" {
		techList = append(techList, vuejs)
	}
	if react != "" {
		techList = append(techList, react)
	}
	if php != "" {
		techList = append(techList, php)
	}
	var newProjectUpdate = Project{
		Name:        projectName,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: description,
		Tech:        techList,
		Image:       image,
	}
	sqlQuery := "UPDATE tb_projects SET name = $2, start_date = $3, end_date = $4, description = $5, technologies = $6, image = $7 WHERE id = $1"
	_, err = connection.Conn.Exec(context.Background(), sqlQuery, id, newProjectUpdate.Name, newProjectUpdate.StartDate, newProjectUpdate.EndDate, newProjectUpdate.Description, newProjectUpdate.Tech, newProjectUpdate.Image)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message" + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func formLoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/login.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message + " + err.Error()))
		return
	}

	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data["IsLogin"] = false
	} else {
		Data["IsLogin"] = true
		Data["Name"] = session.Values["Name"]
	}

	tmpl.Execute(w, Data)
}

func formRegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/register.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message + " + err.Error()))
		return
	}

	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data["IsLogin"] = false
	} else {
		Data["IsLogin"] = true
		Data["Name"] = session.Values["Name"]
	}

	tmpl.Execute(w, Data)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message + " + err.Error()))
		return
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	sqlQuery := "INSERT INTO tb_users (name, email, password) VALUES ($1, $2, $3)"
	_, err = connection.Conn.Exec(context.Background(), sqlQuery, name, email, hashedPassword)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message + " + err.Error()))
		return
	}

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message + " + err.Error()))
		return
	}

	var user = User{}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	sqlQuery := "SELECT id, name, email, password FROM tb_users WHERE email = $1"
	err = connection.Conn.QueryRow(context.Background(), sqlQuery, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message + " + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message + " + err.Error()))
		return
	}

	session, _ := store.Get(r, "SESSION_ID")

	session.Values["IsLogin"] = true
	session.Values["Name"] = user.Name
	session.Values["UserID"] = user.Id
	session.Options.MaxAge = 10800

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) { // Masih ERROR
	w.Header().Set("Cache-Control", "no-chace, no-store, must-revalidate")

	session, _ := store.Get(r, "SESSION_ID")
	session.Values["IsLogin"] = false
	session.Options.MaxAge = -1
	err := session.Save(r, w)

	if err != nil {
		fmt.Println("failed to delete session", err)
	}

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}
