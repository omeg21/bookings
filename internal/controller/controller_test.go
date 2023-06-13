package controller

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct{
	Key string
	value string
}

var theTest = []struct {
	name string
	url string
	method string
	params []postData
	expectedStatusCode int 
}{
	{"home","/","GET",[]postData{},http.StatusOK},
	{"about","/about","GET",[]postData{},http.StatusOK},
	{"primo","/primo","GET",[]postData{},http.StatusOK},
	{"jojo","/jojo","GET",[]postData{},http.StatusOK},
	{"sa","/search-availability","GET",[]postData{},http.StatusOK},
	{"contact","/contact","GET",[]postData{},http.StatusOK},
	{"mr","/reservation","GET",[]postData{},http.StatusOK},
	{"post-search-availability","/search-availability","POST",[]postData{
		{Key: "start", value: "2023-05-29" },
		{Key: "end",value:"2023-05-29"},
	},http.StatusOK},
	{"post-search-availability-JSON","/search-availability-json","POST",[]postData{
		{Key: "start", value: "2023-05-29" },
		{Key: "end",value:"2023-05-29"},
	},http.StatusOK},
	{"make-reservation-post","/reservation","POST",[]postData{
		{Key: "first_name", value: "john" },
		{Key: "last_name",value:"dough"},
		{Key: "email",value:"john@mail.com"},
		{Key: "last_name",value:"9008"},
	},http.StatusOK},
}

func TestController(t *testing.T)  {

	routes := getRoutes()

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTest{
		if e.method == "GET" {
			resp,err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s excpeted %d but got %d",e.name,e.expectedStatusCode,resp.StatusCode)
			}
		}else{
			values := url.Values{}
			for _,x := range e.params{
				values.Add(x.Key,x.value)
			}
			resp,err := ts.Client().PostForm(ts.URL + e.url,values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s excpeted %d but got %d",e.name,e.expectedStatusCode,resp.StatusCode)
			}
		}
	} 
	
}