package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T)  {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	has := form.Has("whatever")

	if has {
		t.Error("form shows has field when it doesn not ")
	}

	postedData := url.Values{}

	postedData.Add("a","a")
	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error(" shws form does not have field  when it should ")
	}
}


func TestForm_MinLength(t *testing.T)  {
	
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.MinLenght("x",10)

	if form.Valid() {
		t.Error(" forms shows minlenght for non existence field")
	}

	isError := form.Errors.Get("x") 
	if isError == ""{
		t.Error("should have an error but pass ")
	}





	postedValues := url.Values{}

	postedValues.Add("some_field","some_value")
	form = New(postedValues)

	form.MinLenght("some_field",100)
	if form.Valid(){
		t.Error("shows min lenght of 100 when data is shorter")
	}

	postedValues = url.Values{}

	postedValues.Add("another_field","abc666")
	form = New(postedValues)

	form.MinLenght("another_field",1)

	if !form.Valid(){
		t.Error("Shows minleght of 1 is not met when it is")
	}

	isError = form.Errors.Get("another_field") 
	if isError != ""{
		t.Error("should not get an error but get one ")
	}

}


func TestForm_IsEmail(t *testing.T)  {
	postedValues := url.Values{}
	form := New(postedValues)

	form.IsEmail("x")
	if form.Valid(){
		t.Error("form shows valid email for non existence field")
	}

	postedValues = url.Values{}

	postedValues.Add("email","me@me.com")
	form = New(postedValues)

	form.IsEmail("email")
	if !form.Valid(){
		t.Error("return invalid email when it should not ")
	}

	postedValues = url.Values{}

	postedValues.Add("email","meme.com")
	form = New(postedValues)

	form.IsEmail("email")
	if form.Valid(){
		t.Error("return a valid email when it should not ")
	}
}


