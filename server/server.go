package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"
)

func WriteJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, msg string, responseCode int) {
	errMsg := `{"error": "` + msg + `"}`
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode)
	w.Write([]byte(errMsg))
}

func GetMeetings(w http.ResponseWriter, r *http.Request) {
	pageNum := r.URL.Query().Get("pageNum")
	fmt.Println("Get meetings")
	fmt.Println(pageNum)

	meets := make([]MeetCard, len(MeetCards))
	i := 0
	for  _, value := range MeetCards {
		meets[i] = value
		i++
	}
	WriteJson(w, meets)
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	pageNum := r.URL.Query().Get("pageNum")

	fmt.Println("Get people")
	fmt.Println(pageNum)

	var users []UserCard
	for _, v := range UserCards {
		users = append(users, v)
	}
	WriteJson(w, users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseUint(r.URL.Query().Get("userId"), 10, 32)
	if err != nil {
		WriteError(w, "user id not found", http.StatusNotFound)
		return
	}
	fmt.Println("Get person")
	fmt.Println(userId)
	profile, ok := UserProfiles[uint(userId)]
	if !ok {
		WriteError(w, "profile not found", http.StatusNotFound)
		return
	}
	WriteJson(w, profile)
}

func EditUser(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("authToken")
	if err != nil {
		WriteError(w, "client unauthorized", http.StatusUnauthorized)
		return
	}
	userId, ok := Sessions[session.Value]
	if !ok {
		WriteError(w, "client unauthorized", http.StatusUnauthorized)
		return
	}
	buf, ok := UserProfiles[userId]
	if !ok {
		WriteError(w, "client profile not found", http.StatusNotFound)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&buf)
	if err != nil {
		WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	UserProfiles[userId] = buf
	w.WriteHeader(http.StatusOK)
}

func GetUserId(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("authToken")
	if err != nil {
		WriteError(w, "client unauthorized", http.StatusUnauthorized)
		return
	}
	uId, ok := Sessions[session.Value]
	if !ok {
		WriteError(w, "client unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	WriteJson(w, UserId{uId})
}

func LogIn(w http.ResponseWriter, r *http.Request)  {
	var userData LogPassPair
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	credData, ok := CredStorage[userData.Login]
	if !ok || credData.Password != userData.Password {
		WriteError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	token := uuid.New()
	cookie := http.Cookie{
		Name:       "authToken",
		Value:      token,
		Expires:    time.Now().Add(30 * 24 * time.Hour),
	}
	http.SetCookie(w, &cookie)
	Sessions[token] = credData.Id
	w.WriteHeader(http.StatusOK)
}

func SignOut(w http.ResponseWriter, r *http.Request)  {
	session, err := r.Cookie("authToken")
	if err == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	delete(Sessions, session.Value)
	expire := time.Now().AddDate(0, 0, -1)
	cookie := http.Cookie{
		Name:       "authToken",
		Value:      session.Value,
		Expires:    expire,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func SignUp(w http.ResponseWriter, r *http.Request)  {
	var userData LogPassPair
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	_, exists := CredStorage[userData.Login]
	if exists {
		WriteError(w, "login has already been taken", http.StatusBadRequest)
		return
	}
	newInd := uint(len(UserProfiles))
	for ;; newInd++ {
		_, existsProfile := UserProfiles[newInd]
		_, existsCard := UserCards[newInd]
		if !existsProfile && !existsCard {
			break
		}
	}
	UserProfiles[newInd] = UserProfile{}
	UserCards[newInd] = UserCard{CardId: newInd, ImgSrc: "assets/luckash.jpeg" }
	CredStorage[userData.Login] = Credentials{
		Id:       newInd,
		Login:    userData.Login,
		Password: userData.Password,
	}
	w.WriteHeader(http.StatusOK)
}

func EditOnSignUp(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/meetings", GetMeetings).Methods("GET")
	r.HandleFunc("/people", GetPeople).Methods("GET")
	r.HandleFunc("/user", GetUser).Methods("GET")
	r.HandleFunc("/user", EditUser).Methods("POST")
	r.HandleFunc("/me", GetUserId).Methods("GET")
	r.HandleFunc("/login", LogIn).Methods("POST")
	r.HandleFunc("/signout", SignOut).Methods("POST")
	r.HandleFunc("/signup", SignUp).Methods("POST")
	r.HandleFunc("/edit_on_signup", EditOnSignUp).Methods("POST")

	r.PathPrefix("/").HandlerFunc(serveStatic)

	fmt.Println("Launching at port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Unable to launch server: ", err)
	}
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	const staticPath = "static"
	fPath := path.Join(staticPath, "index.html")
	if r.URL.Path != "/" {
		fPath = path.Join(staticPath, r.URL.Path)
	}
	http.ServeFile(w, r, fPath)
}
