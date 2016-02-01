package auth

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type Keys struct {
	Authkey string `json:"authkey"`
}

func createNewAuth(key string) *Keys {
	return &Keys{key}
}

func TestToken(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, GetToken("holabolaNotLikeaBull"), http.StatusInternalServerError)
	}
	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		log.Fatal(err)
	}
	res := httptest.NewRecorder()
	handler(res, req)
	got := res.Body.String()
	want := "\n"
	if got == want {
		t.Errorf("want: %#v\n got: %#v ", want, got)
	}
}
func TestAuthWithoutToken(t *testing.T) {
	key := createNewAuth("abc123456abc")
	jsondata, _ := json.Marshal(key)
	post_data := strings.NewReader(string(jsondata))
	req, _ := http.NewRequest("POST", "http://localhost:8000/auth", post_data)
	client := &http.Client{}
	res, _ := client.Do(req)
	got, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	want := []byte{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %#v\n got: %#v ", want, got)
	}
}
func TestAuthWithToken(t *testing.T) {
	req, _ := http.NewRequest("POST", "http://localhost:8000/auth", nil)
	token := returnToken()
	req.Header.Set("Authorization", "bearer "+token)
	client := &http.Client{}
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != 401 {
		t.Errorf("want: %d\n got: %d ", 401, res.StatusCode)
	}
	want := "\"No auth key\"\n"
	if string(body) != want {
		t.Errorf("want: %#v\n got: %#v ", want, string(body))
	}
}
func returnToken() string {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, GetToken("holabolaNotLikeaBull"), http.StatusInternalServerError)
	}
	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		log.Fatal(err)
	}
	res := httptest.NewRecorder()
	handler(res, req)
	return res.Body.String()
}
