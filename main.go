package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

/////cookies
var cookieHandler = securecookie.New(securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

//настроить сессию с пользователем
func SetsSession(userName string,response http.ResponseWriter){
	value := map[string]string{"username":userName}
	encoded, err := cookieHandler.Encode("session",value)
	if err == nil {
		cookie := &http.Cookie{
			Name: "session",
			Value: encoded,
			Path: "/",
		}
		http.SetCookie(response,cookie)
	}
}
//получить имя пользователя из сессии
func GetUserName(request *http.Request)(userName string){
	cookie,err := request.Cookie("session")
	if err == nil {
		cookieValue := make(map[string]string)
		err = cookieHandler.Decode("session",cookie.Value,&cookieValue)
		if err == nil {
			userName = cookieValue["username"]
		}
	}
	return userName
}


func ClearSession(response http.ResponseWriter){
	cookie := &http.Cookie{
		Name: "session",
		Value: "",
		Path: "/",
		MaxAge: -1,
	}
	http.SetCookie(response,cookie)
}

/////cookies

//handlers

// "/login"
var LoginPageHandler = http.HandlerFunc(
	func(w http.ResponseWriter,r *http.Request){
		if r.Method == "GET" {
			parsedTemplate,_ := template.ParseFiles("templates/loginPage.html")
			parsedTemplate.Execute(w,nil)
		} else {
			username := r.FormValue("username")
			password := r.FormValue("password")
			target   := "/login"
			if username != "" && password != ""{
				SetsSession(username,w)
				target = "/books"
			}
			http.Redirect(w,r,target,302)
		}


	})
// "/books"
var BooksPageHandler = http.HandlerFunc(
	func(w http.ResponseWriter,r *http.Request){
		Username := GetUserName(r)
		if Username != ""{

			parsedTemplate,_ := template.ParseFiles("templates/books.html")
			parsedTemplate.Execute(w,Data)
		} else {
			http.Redirect(w,r,"/login",302)
		}

	})


var CreateBookHandler = http.HandlerFunc(
	func(w http.ResponseWriter,r *http.Request) {
		Username := GetUserName(r)
		if Username == ""{
			http.Redirect(w,r,"/login",302)
			return
		}

		if r.Method == "GET" {
			parsedTemplate,_ := template.ParseFiles("templates/createBook.html")
			parsedTemplate.Execute(w,nil)
		} else {
			title := r.FormValue("title")
			pagecount := r.FormValue("pagecount")
			author := r.FormValue("author")

			count, err := strconv.Atoi(pagecount)
			if err != nil {
				count = 0
			}

			Data.Books = append(Data.Books,Book{
				Title:  title,
				Pages:  count,
				Author: author,
			})

			http.Redirect(w,r,"/books/create",302)
		}
	})



var JournalPageHandler = http.HandlerFunc(
	func(w http.ResponseWriter,r *http.Request){
		Username := GetUserName(r)
		if Username != ""{
			parsedTemplate,_ := template.ParseFiles("templates/journal.html")
			parsedTemplate.Execute(w,DataJ)
		} else {
			http.Redirect(w,r,"/login",302)
		}
	})

var CreateJournalHandler = http.HandlerFunc(
	func(w http.ResponseWriter,r *http.Request){
		Username := GetUserName(r)
		if Username == ""{
			http.Redirect(w,r,"/login",302)
			return
		}

		if r.Method == "GET" {
			parsedTemplate,_ := template.ParseFiles("templates/createJournal.html")
			parsedTemplate.Execute(w,nil)
		} else {
			redactor := r.FormValue("redactor")
			serialnumber := r.FormValue("serialnumber")
			edition := r.FormValue("edition")
			pagecount := r.FormValue("pagecount")

			serNum, err := strconv.Atoi(serialnumber)
			if err != nil {
				serNum = 0
			}
			count, err := strconv.Atoi(pagecount)
			if err != nil {
				count = 0
			}

			DataJ.Journals = append(DataJ.Journals,Journal{
				Redactor:     redactor,
				SerialNumber: serNum,
				Edition:      edition,
				PageCount:    count,
			})

			http.Redirect(w,r,"/journals/create",302)
		}
	})

var LogoutFormPageHandler = func(w http.ResponseWriter,r *http.Request){
	ClearSession(w)
	http.Redirect(w,r,"/login",302)
}

 func Hello(w http.ResponseWriter,r *http.Request){
	fmt.Fprintf(w,"Hello")
}


//handlers


type User struct {
	Username string
	Password string
}

type Book struct {
	Title string
	Pages int
	Author string
}

type Journal struct {
	Redactor string //Редактор
	SerialNumber int //Серийный номер
	Edition string //Издание
	PageCount int //Количество страниц
}

type BookCollection struct {
	Books []Book
}

type JournalCollection struct {
	Journals []Journal
}

var Data BookCollection
var DataJ JournalCollection

func init(){
	Data.Books = append(Data.Books,Book{
		Title:     "Book 1",
		Pages: 712,
		Author:    "user",
	})

	DataJ.Journals = append(DataJ.Journals,Journal{
		Redactor:     "New User",
		SerialNumber: 554,
		Edition:      "Special",
		PageCount:    77,
	})
}


const (
	connPort = "8080"
	connHost = "0.0.0.0"
)



func main(){
	router := mux.NewRouter()

	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	router.HandleFunc("/",Hello)

	router.Handle("/login",handlers.LoggingHandler(logFile,
		LoginPageHandler)).Methods("GET","POST")

	router.Handle("/books",handlers.LoggingHandler(logFile,
		BooksPageHandler)).Methods("GET")

	router.Handle("/books/create",handlers.LoggingHandler(logFile,
		CreateBookHandler)).Methods("GET","POST")

	router.Handle("/journals/create",handlers.LoggingHandler(logFile,
		CreateJournalHandler)).Methods("GET","POST")

	router.Handle("/journals",handlers.LoggingHandler(logFile,
		JournalPageHandler)).Methods("GET")

	router.Handle("/logout",handlers.LoggingHandler(logFile,
		http.HandlerFunc(LogoutFormPageHandler))).Methods("POST")
	fmt.Println("listening:"+connPort)
	err = http.ListenAndServe(connHost+":"+connPort,router)
	if err != nil {
		log.Fatal("error starting server: ",err)
		return
	}

}
