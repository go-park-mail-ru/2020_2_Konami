package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
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

func CreateSession(w http.ResponseWriter, uId int) {
	token := uuid.New()
	cookie := http.Cookie{
		Name:    "authToken",
		Value:   token,
		Expires: time.Now().Add(30 * 24 * time.Hour),
	}
	http.SetCookie(w, &cookie)
	Sessions[token] = uId
}

func GetMeetings(w http.ResponseWriter, r *http.Request) {
	meetings := make([]*Meeting, len(MeetingStorage))
	i := 0
	for _, value := range MeetingStorage {
		meetings[i] = value
		i++
	}
	if len(meetings) == 0 {
		WriteError(w, "no meetings found", http.StatusNotFound)
		return
	}
	sort.Sort(MeetingsByDate(meetings))
	WriteJson(w, meetings)
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	users := make([]*User, len(UserStorage))
	i := 0
	for _, value := range UserStorage {
		users[i] = value
		i++
	}
	if len(users) == 0 {
		WriteError(w, "no users found", http.StatusNotFound)
		return
	}
	sort.Sort(UsersByName(users))
	WriteJson(w, users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil {
		WriteError(w, "user id not found", http.StatusNotFound)
		return
	}
	profile, ok := UserStorage[userId]
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
	buf := &UserUpdate{}
	err = json.NewDecoder(r.Body).Decode(&buf)
	if err != nil {
		log.Println(err)
		WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	usr, exists := UserStorage[userId]
	if !exists {
		WriteError(w, "profile not found", http.StatusNotFound)
	}
	ok = CommitUserUpdate(buf, usr)
	if !ok {
		WriteError(w, "unable to update profile", http.StatusBadRequest)
	}
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

func LogIn(w http.ResponseWriter, r *http.Request) {
	var userData Credentials
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil || userData.Login == "" || userData.Password == "" {
		WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	credData, ok := CredStorage[userData.Login]
	if !ok {
		WriteError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	cmpRes := bcrypt.CompareHashAndPassword([]byte(credData.Password), []byte(userData.Password))
	if cmpRes != nil {
		WriteError(w, "invalid credentials", http.StatusUnauthorized)
	}
	CreateSession(w, credData.uId)
	w.WriteHeader(http.StatusOK)
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("authToken")
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	delete(Sessions, session.Value)
	expire := time.Now().AddDate(0, 0, -1)
	cookie := http.Cookie{
		Name:    "authToken",
		Value:   session.Value,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	_, exists := CredStorage[creds.Login]
	if exists {
		WriteError(w, "login has already been taken", http.StatusConflict)
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.MinCost)
	if err != nil {
		WriteError(w, "internal error", http.StatusInternalServerError)
	}
	creds.Password = string(hashed)
	newInd := rand.Intn(1 << 30)
	for ; ; newInd = rand.Int() {
		_, existsProfile := UserStorage[newInd]
		if !existsProfile {
			break
		}
	}
	creds.uId = newInd
	UserStorage[newInd] = &User{
		Id:           newInd,
		ImgSrc:       "assets/luckash.jpeg",
		InterestTags: []string{},
		SkillTags:    []string{},
		Meetings:     []*Meeting{},
	}
	CredStorage[creds.Login] = &creds
	CreateSession(w, newInd)
	w.WriteHeader(http.StatusOK)
}

func UploadUserPic(w http.ResponseWriter, r *http.Request) {
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
	profile, exists := UserStorage[userId]
	if !exists {
		WriteError(w, "profile not found", http.StatusNotFound)
		return
	}
	err = r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		WriteError(w, "invalid multipart form", http.StatusBadRequest)
		return
	}
	file, handler, err := r.FormFile("fileToUpload")
	if err != nil {
		WriteError(w, "invalid form file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	fname := strings.Split(handler.Filename, ".")
	ext := fname[len(fname)-1]
	if ext != "jpg" && ext != "jpeg" && ext != "png" && ext != "gif" {
		WriteError(w, "invalid file format", http.StatusBadRequest)
	}
	imgPath := "uploads/userpics/" + strconv.Itoa(userId) + "." + ext

	f, err := os.OpenFile(imgPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		WriteError(w, "unable to create file", http.StatusInternalServerError)
		return
	}
	defer f.Close()
	var written int64 = 0
	written, err = io.Copy(f, file)
	if err != nil || written == 0 {
		WriteError(w, "unable to save file", http.StatusInternalServerError)
		return
	}
	profile.ImgSrc = imgPath
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
	r.HandleFunc("/logout", LogOut).Methods("POST")
	r.HandleFunc("/signup", SignUp).Methods("POST")
	r.HandleFunc("/images", UploadUserPic).Methods("POST")

	r.PathPrefix("/uploads/").HandlerFunc(serveUploads)
	r.PathPrefix("/").HandlerFunc(serveStatic)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	certFile := os.Getenv("CERTFILE")
	if certFile == "" {
		certFile = "/etc/letsencrypt/live/okto.pw/fullchain.pem"
	}
	keyFile := os.Getenv("KEYFILE")
	if keyFile == "" {
		keyFile = "/etc/letsencrypt/live/okto.pw/privkey.pem"
	}
	tlsHost := os.Getenv("TLSHOST")
	if tlsHost == "" {
		tlsHost = "okto.pw"
	}
	tlsPort := os.Getenv("TLSPORT")

	var err error = nil

	if tlsPort == "" {
		log.Println("Launching at HTTP port " + port)
		err = http.ListenAndServe(":"+port, r)
	} else {
		go redirectToHTTPS(port, tlsHost, tlsPort)
		log.Println("Launching at HTTPS port " + tlsPort)
		err = http.ListenAndServeTLS(":"+tlsPort, certFile, keyFile, r)
	}

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

func serveUploads(w http.ResponseWriter, r *http.Request) {
	relPath := strings.TrimPrefix(r.URL.Path, "/")
	http.ServeFile(w, r, relPath)
}

func redirectToHTTPS(port, tlsHost, tlsPort string) {
	log.Println("Redirect from :" + port + " to " + tlsHost + ":" + tlsPort)
	httpSrv := http.Server{
		Addr: ":" + port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u := r.URL
			u.Host = net.JoinHostPort(tlsHost, tlsPort)
			u.Scheme = "https"
			http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
		}),
	}
	log.Println(httpSrv.ListenAndServe())
}
