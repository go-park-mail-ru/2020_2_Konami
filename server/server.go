package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/bcrypt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func WriteJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, resp *ErrResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.ResponseCode)
	json.NewEncoder(w).Encode(resp)
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

func GetMeetingsList(w http.ResponseWriter, r *http.Request) {
	var meetings []*Meeting
	todayOnly := r.URL.Query().Get("today") == "true"
	tomorrowOnly := r.URL.Query().Get("tomorrow") == "true"
	myOnly := r.URL.Query().Get("mymeetings") == "true"
	favOnly := r.URL.Query().Get("favorites") == "true"
	currentTime := time.Now()
	today := currentTime.Format("2006-01-02")
	tomorrow := currentTime.Add(24 * time.Hour).Format("2006-01-02")

	var userId int
	if myOnly || favOnly {
		session, err := r.Cookie("authToken")
		var ok bool
		userId, ok = Sessions[session.Value]
		if err != nil || !ok {
			WriteError(w, &ErrResponse{http.StatusUnauthorized, "client unauthorized"})
			return
		}
	}

	for _, value := range MeetingStorage {
		mDate := value.StartDate[:strings.Index(value.StartDate, " ")]
		if todayOnly && mDate == today ||
			tomorrowOnly && mDate == tomorrow ||
			myOnly && UserRegistered(userId, value.Id) ||
			favOnly && UserLikes(userId, value.Id) ||
			!todayOnly && !tomorrowOnly && !myOnly && !favOnly {
			meetings = append(meetings, value)
		}
	}
	if len(meetings) == 0 {
		WriteError(w, &ErrResponse{http.StatusNotFound, "no meetings found"})
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
		WriteError(w, &ErrResponse{http.StatusNotFound, "no users found"})
		return
	}
	sort.Sort(UsersByName(users))
	WriteJson(w, users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil {
		WriteError(w, &ErrResponse{http.StatusNotFound, "user id not found"})
		return
	}
	profile, ok := UserStorage[userId]
	if !ok {
		WriteError(w, &ErrResponse{http.StatusNotFound, "profile not found"})
		return
	}
	WriteJson(w, profile)
}

func EditUser(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("authToken")
	if err != nil {
		WriteError(w, &ErrResponse{http.StatusUnauthorized, "client unauthorized"})
		return
	}
	userId, ok := Sessions[session.Value]
	if !ok {
		WriteError(w, &ErrResponse{http.StatusUnauthorized, "client unauthorized"})
		return
	}
	buf := &UserUpdate{}
	err = json.NewDecoder(r.Body).Decode(&buf)
	if err != nil {
		log.Println(err)
		WriteError(w, &ErrResponse{http.StatusBadRequest, "invalid request body"})
		return
	}
	usr, exists := UserStorage[userId]
	if !exists {
		WriteError(w, &ErrResponse{http.StatusNotFound, "profile not found"})
		return
	}
	ok = CommitUserUpdate(buf, usr)
	if !ok {
		WriteError(w, &ErrResponse{http.StatusBadRequest, "unable to update profile"})
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetUserId(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("authToken")
	if err != nil {
		WriteError(w, &ErrResponse{http.StatusUnauthorized, "client unauthorized"})
		return
	}
	uId, ok := Sessions[session.Value]
	if !ok {
		WriteError(w, &ErrResponse{http.StatusUnauthorized, "client unauthorized"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	WriteJson(w, UserId{uId})
}

func LogIn(w http.ResponseWriter, r *http.Request) {
	var userData Credentials
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil || userData.Login == "" || userData.Password == "" {
		WriteError(w, &ErrResponse{http.StatusBadRequest, "invalid request body"})
		return
	}
	credData, ok := CredStorage[userData.Login]
	if !ok {
		WriteError(w, &ErrResponse{http.StatusUnauthorized, "invalid credentials"})
		return
	}
	cmpRes := bcrypt.CompareHashAndPassword([]byte(credData.Password), []byte(userData.Password))
	if cmpRes != nil {
		WriteError(w, &ErrResponse{http.StatusUnauthorized, "invalid credentials"})
		return
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
		WriteError(w, &ErrResponse{http.StatusBadRequest, "invalid request body"})
		return
	}
	_, exists := CredStorage[creds.Login]
	if exists {
		WriteError(w, &ErrResponse{http.StatusConflict, "login has already been taken"})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.MinCost)
	if err != nil {
		WriteError(w, &ErrResponse{http.StatusInternalServerError, "internal error"})
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
		ImgSrc:       "assets/empty-avatar.jpeg",
		InterestTags: []string{},
		SkillTags:    []string{},
		Meetings:     []*UserMeeting{},
	}
	CredStorage[creds.Login] = &creds
	CreateSession(w, newInd)
	w.WriteHeader(http.StatusOK)
}

func UploadUserPic(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("authToken")
	if err != nil {
		WriteError(w, &ErrResponse{http.StatusUnauthorized, "client unauthorized"})
		return
	}
	userId, ok := Sessions[session.Value]
	if !ok {
		WriteError(w, &ErrResponse{http.StatusUnauthorized, "client unauthorized"})
		return
	}
	profile, exists := UserStorage[userId]
	if !exists {
		WriteError(w, &ErrResponse{http.StatusNotFound, "profile not found"})
		return
	}
	err = r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		WriteError(w, &ErrResponse{http.StatusBadRequest, "invalid multipart form"})
		return
	}
	file, handler, err := r.FormFile("fileToUpload")
	if err != nil {
		WriteError(w, &ErrResponse{http.StatusBadRequest, "invalid form file"})
		return
	}
	defer file.Close()
	fname := strings.Split(handler.Filename, ".")
	ext := fname[len(fname)-1]
	if ext != "jpg" && ext != "jpeg" && ext != "png" {
		WriteError(w, &ErrResponse{http.StatusBadRequest, "invalid file format"})
		return
	}
	imgPath := "uploads/userpics/" + strconv.Itoa(userId) + "." + ext

	f, err := os.OpenFile(imgPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		WriteError(w, &ErrResponse{http.StatusInternalServerError, "unable to create file"})
		return
	}
	defer f.Close()
	var written int64 = 0
	written, err = io.Copy(f, file)
	if err != nil || written == 0 {
		WriteError(w, &ErrResponse{http.StatusInternalServerError, "unable to save file"})
		return
	}
	profile.ImgSrc = imgPath
	w.WriteHeader(http.StatusOK)
}

func CreateMeeting(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("authToken")
	if err != nil {
		WriteError(w, &ErrResponse{http.StatusUnauthorized, "client unauthorized"})
		return
	}
	userId, sOk := Sessions[session.Value]
	author, uOk := UserStorage[userId]
	if !sOk || !uOk {
		WriteError(w, &ErrResponse{http.StatusUnauthorized, "client unauthorized"})
		return
	}
	mData := &MeetingUpload{}
	err = json.NewDecoder(http.MaxBytesReader(w, r.Body, 10*1024*1024)).Decode(&mData)
	if err != nil {
		log.Println(err)
		WriteError(w, &ErrResponse{http.StatusBadRequest, "invalid request body"})
		return
	}
	newInd := rand.Intn(1 << 30)
	for ; ; newInd = rand.Int() {
		_, existsMeeting := MeetingStorage[newInd]
		if !existsMeeting {
			break
		}
	}
	imgPath := "assets/paris.jpg"
	if len(mData.Photo) != 0 {
		// TODO: separate image handling
		jpegPrefix := "data:image/jpeg;base64,"
		pngPrefix := "data:image/png;base64,"
		var img image.Image
		err = nil
		if strings.HasPrefix(mData.Photo, jpegPrefix) {
			rawImage := mData.Photo[len(jpegPrefix):]
			decoded, _ := base64.StdEncoding.DecodeString(rawImage)
			img, err = jpeg.Decode(bytes.NewReader(decoded))
		} else if strings.HasPrefix(mData.Photo, pngPrefix) {
			rawImage := mData.Photo[len(pngPrefix):]
			decoded, _ := base64.StdEncoding.DecodeString(rawImage)
			img, err = png.Decode(bytes.NewReader(decoded))
		} else {
			WriteError(w, &ErrResponse{http.StatusBadRequest, "invalid file format"})
			return
		}
		imgPath = "uploads/meetingpics/" + strconv.Itoa(newInd) + ".png"

		f, err := os.OpenFile(imgPath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			WriteError(w, &ErrResponse{http.StatusInternalServerError, "unable to create file"})
			return
		}
		defer f.Close()
		err = png.Encode(f, img)
		if err != nil {
			WriteError(w, &ErrResponse{http.StatusInternalServerError, "unable to save file"})
			return
		}
	}

	stDateTrimmed := strings.Replace(mData.Start, "T", " ", -1)
	endDateTrimmed := strings.Replace(mData.End, "T", " ", -1)
	meeting := &Meeting{
		Id:        newInd,
		AuthorId:  userId,
		Title:     mData.Name,
		Text:      mData.Description,
		ImgSrc:    imgPath,
		Tags:      mData.Tags,
		Place:     mData.City + ", " + mData.Address,
		StartDate: stDateTrimmed,
		EndDate:   endDateTrimmed,
		Seats:     100500,
		SeatsLeft: 100500,
	}
	if meeting.Tags == nil {
		meeting.Tags = []string{}
	}
	MeetingStorage[newInd] = meeting
	author.Meetings = append(author.Meetings, &UserMeeting{
		Title:  mData.Name,
		ImgSrc: imgPath,
		Link:   fmt.Sprintf("/meet?meetId=%d", newInd),
	})
	w.WriteHeader(http.StatusCreated)
}

func GetMeeting(w http.ResponseWriter, r *http.Request) {
	meetId, err := strconv.Atoi(r.URL.Query().Get("meetId"))
	if err != nil {
		WriteError(w, &ErrResponse{http.StatusNotFound, "user id not found"})
		return
	}
	meeting, ok := MeetingStorage[meetId]
	if !ok {
		WriteError(w, &ErrResponse{http.StatusNotFound, "profile not found"})
		return
	}
	userId := -1
	session, err := r.Cookie("authToken")
	if err == nil {
		userId, ok = Sessions[session.Value]
		if !ok {
			userId = -1
		}
	}
	if userId != -1 {
		meeting.Like = UserLikes(userId, meetId)
		meeting.Reg = UserRegistered(userId, meetId)
	}
	WriteJson(w, meeting)
}

func UpdateMeeting(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("authToken")
	if err != nil {
		WriteError(w, &ErrResponse{http.StatusUnauthorized, "client unauthorized"})
		return
	}
	userId, ok := Sessions[session.Value]
	if !ok {
		WriteError(w, &ErrResponse{http.StatusUnauthorized, "client unauthorized"})
		return
	}
	mData := &MeetingUpdate{}
	buf := new(strings.Builder)
	_, err = io.Copy(buf, r.Body)
	if err != nil {
		WriteError(w, &ErrResponse{http.StatusBadRequest, "invalid request body"})
		return
	}
	mStr := buf.String()
	err = json.Unmarshal([]byte(mStr), &mData)
	if err != nil {
		WriteError(w, &ErrResponse{http.StatusBadRequest, "invalid request body"})
		return
	}
	meeting, exists := MeetingStorage[mData.MeetId]
	if !exists {
		WriteError(w, &ErrResponse{http.StatusBadRequest, "invalid request body"})
		return
	}
	if meeting.StartDate < time.Now().Format("2006-01-02 15:04:05") {
		WriteError(w, &ErrResponse{http.StatusConflict, "meeting has already started"})
		return
	}
	if strings.Contains(mStr, "isLiked") {
		if mData.Fields.Like {
			SetEl(userId, mData.MeetId, Likes)
		} else {
			RemoveEl(userId, mData.MeetId, Likes)
		}
	}
	if strings.Contains(mStr, "isRegistered") {
		if mData.Fields.Reg {
			if meeting.SeatsLeft == 0 {
				WriteError(w, &ErrResponse{http.StatusConflict, "no more vacant seats left"})
				return
			}
			SetEl(userId, mData.MeetId, Registrations)
		} else {
			RemoveEl(userId, mData.MeetId, Registrations)
			meeting.SeatsLeft += 1
		}
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/meetings", GetMeetingsList).Methods("GET")
	r.HandleFunc("/api/meeting", CreateMeeting).Methods("POST")
	r.HandleFunc("/api/meet", GetMeeting).Methods("GET")
	r.HandleFunc("/api/meet", UpdateMeeting).Methods("POST")
	r.HandleFunc("/api/people", GetPeople).Methods("GET")
	r.HandleFunc("/api/user", GetUser).Methods("GET")
	r.HandleFunc("/api/user", EditUser).Methods("POST")
	r.HandleFunc("/api/me", GetUserId).Methods("GET")
	r.HandleFunc("/api/login", LogIn).Methods("POST")
	r.HandleFunc("/api/logout", LogOut).Methods("POST")
	r.HandleFunc("/api/signup", SignUp).Methods("POST")
	r.HandleFunc("/api/images", UploadUserPic).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
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
