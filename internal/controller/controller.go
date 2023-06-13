package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/omeg21/project-repo/internal/config"
	"github.com/omeg21/project-repo/internal/driver"
	"github.com/omeg21/project-repo/internal/forms"
	"github.com/omeg21/project-repo/internal/helpers"
	"github.com/omeg21/project-repo/internal/models"
	"github.com/omeg21/project-repo/internal/render"
	"github.com/omeg21/project-repo/internal/repository"
	"github.com/omeg21/project-repo/internal/repository/dbrepo"
)

// Repo the repository used by the controller
var Repo *Repository

// Repository is the repositry type
type Repository struct {
	App *config.AppConfig
	DB repository.DatabaseRepo
}

//NewRepo creates a new repository
func NewRepo(a *config.AppConfig,db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB: dbrepo.NewPostgresRepo(db.SQL,a),
	}
}

// NewController sets the repository for the handlers
func NewController(r *Repository) {
	Repo = r
}

// Home page controller
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r,  "home.page.html", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {



	//send the data to the template
	render.Template(w,r, "about.page.html", &models.TemplateData{
	})
}

// func New(w http.ResponseWriter, r *http.Request) {
// 	render.Template(w, "new.page.html")
// }

// Renders the Primo roompage
func (m *Repository) Primo(w http.ResponseWriter, r *http.Request) {
	render.Template(w,r,"generals.page.html",&models.TemplateData{})
	
}
// Renders the Jojo roompage
func (m *Repository) Jojo(w http.ResponseWriter, r *http.Request) {
	render.Template(w,r,"majors.page.html",&models.TemplateData{})
	
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w,r,"searchAvailability.page.html",&models.TemplateData{})
	
}
// Post Availability form
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate,err := time.Parse(layout , start )
	if err != nil {
		helpers.ServerError(w,err)
		return
	}
	endDate,err := time.Parse(layout , end )
	if err != nil {
		helpers.ServerError(w,err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate,endDate)
	if err != nil {
		helpers.ServerError(w,err)
		return
	}

	if len(rooms) == 0{
		//No availability
		m.App.Session.Put(r.Context(),"error","No availability")
		http.Redirect(w,r,"/Search-availability",http.StatusSeeOther)
		return
	}

	data :=make(map[string]interface{})
	data["rooms"] = rooms


	res := models.Reservation{
		StartDate: startDate,
		EndDate: endDate,
	}
	
	m.App.Session.Put(r.Context(),"reservation",res)



	render.Template(w, r, "choose-room.page.html", &models.TemplateData{
		Data: data,
	})
	
}
type jsonResponse struct{
	OK bool `json:"ok"`
	Message string `json:"message"`
}
// Post JSON response form
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK :false ,
		Message: "Available!",
	}

	out,err := json.MarshalIndent(resp,"","     ")
	if err != nil {
		helpers.ServerError(w,err)
		return
	}
	log.Println(string(out))
	w.Header().Set("Content-type","appllication/json")
	w.Write(out)
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w,r,"contact.page.html",&models.TemplateData{})
	
}

//ReservationSummary
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("can't get reservation from session")
		m.App.Session.Put(r.Context(),"error","cannot get reservation from session")
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(),"reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservationSummary.page.html", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomId)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find room!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "makeReservation.page.html", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

//PostReservation handles the posting of Reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w,err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	//01/02 03:04:05PM '06 -0700 to 2023-10-01

	layout := "2006-01-02"
	startDate,err := time.Parse(layout , sd )
	if err != nil {
		helpers.ServerError(w,err)
		return
	}
	endDate,err := time.Parse(layout , ed )
	if err != nil {
		helpers.ServerError(w,err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w,err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate: endDate,
		RoomId: roomID,
	}

	form := forms.New(r.PostForm)


	form.Required("first_name","last_name","email","phone")
	form.MinLenght("first_name",3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(w, r, "makeReservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationID,err := m.DB.InsertReservation(reservation)
	if err != nil{
		helpers.ServerError(w,err)
		return
	}

	m.App.Session.Put(r.Context(),"reservation",reservation)

	http.Redirect(w,r,"/ReservationSummary",http.StatusSeeOther)

	restriction := models.RoomRestriction{
		StartDate: startDate,
		EndDate: endDate,
		RoomId: roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,

	}

	err = m.DB.InsertRoomRestriction(restriction)


}

func (m *Repository) ChooseRoom(w http.ResponseWriter,r *http.Request)  {
	roomID,err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil{
		helpers.ServerError(w,err)
		return
	}


	res,ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w,err)
		return
	}

	res.RoomId = roomID
	
	m.App.Session.Put(r.Context(),"reservation",res)
	http.Redirect(w,r,"/reservation",http.StatusSeeOther)
}