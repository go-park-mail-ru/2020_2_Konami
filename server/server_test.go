package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type TestCase struct {
	Param      string
	Response   string
	StatusCode int
}

var URL = "http://localhost:8080"

func TestGetMeetings(t *testing.T) {
	waitRes := `[{"id":0,"title":"Забив с++","text":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco labori","imgSrc":"assets/paris.jpg","tags":["C++"],"place":"Москва, улица Колотушкина, дом Пушкина","date":"2020-11-10"},{"id":1,"title":"Python for Web","text":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco labori","imgSrc":"assets/paris.jpg","tags":["Python","Web"],"place":"СПБ, улица Вязов, д.1","date":"2020-11-12"}]` + "\n"
	urlTest := URL + "/meetings"
	req := httptest.NewRequest("GET", urlTest, nil)
	w := httptest.NewRecorder()
	GetMeetingsList(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Wrong StatusCode: got %d, expected %d",
			w.Code, http.StatusOK)
	}

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)
	/*
		fmt.Println(bodyStr)
		resArr := make([]string, 1)
		res, _ := json.Marshal(MeetingStorage[0])
		fmt.Println(string(res))
		resArr = append(resArr, string(res))
		fmt.Println(resArr)
	*/
	if bodyStr != waitRes {
		t.Errorf("Wrong Response: got %+v, expected %+v",
			bodyStr, bodyStr)
	}
}

func TestGetPeople(t *testing.T) {
	waitRes := `[{"id":0,"name":"Александр","gender":"M","birthday":"1990-09-12","city":"Нурсултан","email":"lucash@mail.ru","telegram":"","vk":"https://vk.com/id241926559","meetingTags":["RandomTag1","RandomTag5"],"education":"МГТУ им. Н. Э. Баумана до 2010","job":"MAIL RU GROUP","imgSrc":"assets/luckash.jpeg","aims":"Хочу от жизни всего","interestTags":["Шыпшына","Бульба"],"interests":"Люблю, когда встаешь утром, а на столе #Шыпшына и #Бульба","skillTags":["Мелиорация"],"skills":"#Мелиорация - это моя жизнь","meetings":[]},{"id":1,"name":"Роман","gender":"M","birthday":"2000-09-10","city":"Москва","email":"lucash2@mail.ru","telegram":"","vk":"https://vk.com/id420","meetingTags":["RandomTag1","RandomTag5"],"education":"","job":"HH.ru","imgSrc":"assets/luckash.jpg","aims":"","interestTags":["ДВП","ДСП"],"interests":"Люблю клеить #ДВП и #ДСП","skillTags":["Деревообработка"],"skills":"Моя жизнь - это #Деревообработка","meetings":[]}]` + "\n"
	urlTest := URL + "/people"
	req := httptest.NewRequest("GET", urlTest, nil)
	w := httptest.NewRecorder()
	GetPeople(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Wrong StatusCode: got %d, expected %d",
			w.Code, http.StatusOK)
	}

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)

	if bodyStr != waitRes {
		t.Errorf("Wrong Response: got %+v, expected %+v",
			bodyStr, bodyStr)
	}
}

func TestGetUser(t *testing.T) {
	urlTest := URL + "/user"

	cases := []TestCase{
		TestCase{
			Param:      `0`,
			Response:   "{\"id\":0,\"name\":\"Александр\",\"gender\":\"M\",\"birthday\":\"1990-09-12\",\"city\":\"Нурсултан\",\"email\":\"lucash@mail.ru\",\"telegram\":\"\",\"vk\":\"https://vk.com/id241926559\",\"meetingTags\":[\"RandomTag1\",\"RandomTag5\"],\"education\":\"МГТУ им. Н. Э. Баумана до 2010\",\"job\":\"MAIL RU GROUP\",\"imgSrc\":\"assets/luckash.jpeg\",\"aims\":\"Хочу от жизни всего\",\"interestTags\":[\"Шыпшына\",\"Бульба\"],\"interests\":\"Люблю, когда встаешь утром, а на столе #Шыпшына и #Бульба\",\"skillTags\":[\"Мелиорация\"],\"skills\":\"#Мелиорация - это моя жизнь\",\"meetings\":[]}\n",
			StatusCode: http.StatusOK,
		},
		TestCase{
			Param:      `2`,
			Response:   `{"error": "profile not found"}`,
			StatusCode: http.StatusNotFound,
		},
		TestCase{
			Param:      ``,
			Response:   `{"error": "user id not found"}`,
			StatusCode: http.StatusNotFound,
		},
	}

	for caseNum, item := range cases {
		req := httptest.NewRequest("GET", urlTest, nil)
		q := req.URL.Query()
		q.Add("userId", item.Param)
		req.URL.RawQuery = q.Encode()

		w := httptest.NewRecorder()
		GetUser(w, req)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] Wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, http.StatusOK)
		}

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)

		if bodyStr != item.Response {
			fmt.Println(bodyStr, item.Response)
			t.Errorf("[%d] Wrong Response: got %+v, expected %+v",
				caseNum, bodyStr, item.Response)
		}
	}
}

func TestSign(t *testing.T) {
	urlTest := URL + "/people"
	credit := Credentials{
		Login:    "12345@mail.ru",
		Password: "12345",
		uId:      100,
	}
	jsonData, _ := json.Marshal(credit)
	jsonReader := strings.NewReader(string(jsonData))
	req := httptest.NewRequest("POST", urlTest, jsonReader)
	w := httptest.NewRecorder()

	LogIn(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Wrong StatusCode: got %d, expected %d",
			w.Code, http.StatusUnauthorized)
	}

	jsonData, _ = json.Marshal(credit)
	jsonReader = strings.NewReader(string(jsonData))
	req = httptest.NewRequest("POST", urlTest, jsonReader)
	w = httptest.NewRecorder()
	SignUp(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Wrong StatusCode: got %d, expected %d",
			w.Code, http.StatusOK)
	}

	jsonData, _ = json.Marshal(credit)
	jsonReader = strings.NewReader(string(jsonData))
	req = httptest.NewRequest("POST", urlTest, jsonReader)
	w = httptest.NewRecorder()
	SignUp(w, req)
	if w.Code != http.StatusConflict {
		t.Errorf("Wrong StatusCode: got %d, expected %d",
			w.Code, http.StatusConflict)
	}

	jsonData, _ = json.Marshal(credit)
	jsonReader = strings.NewReader(string(jsonData))
	req = httptest.NewRequest("POST", urlTest, jsonReader)
	w = httptest.NewRecorder()
	LogIn(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Wrong StatusCode: got %d, expected %d",
			w.Code, http.StatusOK)
	}
}

func TestLogOut(t *testing.T) {
	urlTest := URL + "/logout"
	req := httptest.NewRequest("POST", urlTest, nil)
	w := httptest.NewRecorder()
	cookie := http.Cookie{
		Name:    "authToken",
		Value:   "testCookie",
		Expires: time.Now().Add(30 * 24 * time.Hour),
	}

	http.SetCookie(w, &cookie)
	req.AddCookie(&cookie)
	LogOut(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Wrong StatusCode: got %d, expected %d",
			w.Code, http.StatusOK)
	}
}

func TestGetUserId(t *testing.T) {
	Sessions["123"] = 0

	urlTest := URL + "/logout"
	req := httptest.NewRequest("POST", urlTest, nil)
	w := httptest.NewRecorder()
	cookie := http.Cookie{
		Name:    "authToken",
		Value:   "123",
		Expires: time.Now().Add(30 * 24 * time.Hour),
	}

	http.SetCookie(w, &cookie)
	req.AddCookie(&cookie)
	GetUserId(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Wrong StatusCode: got %d, expected %d",
			w.Code, http.StatusOK)
	}
}

func TestEditUser(t *testing.T) {
	urlTest := URL + "/user"
	testStr := "test"
	credit := UserUpdate{
		Name:      &testStr,
		City:      &testStr,
		Telegram:  &testStr,
		Vk:        &testStr,
		Education: &testStr,
		Job:       &testStr,
		Aims:      &testStr,
		Interests: &testStr,
		Skills:    &testStr,
	}

	jsonData, _ := json.Marshal(credit)
	jsonReader := strings.NewReader(string(jsonData))
	req := httptest.NewRequest("POST", urlTest, jsonReader)
	w := httptest.NewRecorder()
	cookie := http.Cookie{
		Name:    "authToken",
		Value:   "123",
		Expires: time.Now().Add(30 * 24 * time.Hour),
	}
	http.SetCookie(w, &cookie)
	req.AddCookie(&cookie)
	EditUser(w, req)
}
