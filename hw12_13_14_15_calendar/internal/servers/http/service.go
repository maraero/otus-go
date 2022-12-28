package server_http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/app"
	event "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/domain"
)

type CreatedEvent struct {
	ID int64 `json:"id"`
}

type EventList struct {
	List []event.Event `json:"list"`
}

func handleCreateEvent(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		evt := event.Event{}
		err = json.Unmarshal(body, &evt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := app.EventService.CreateEvent(r.Context(), evt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		writeJson(w, CreatedEvent{ID: id})
	}
}

func handleUpdateEvent(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		evt := event.Event{}
		err = json.Unmarshal(body, &evt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		params := mux.Vars(r)
		idParam := params["id"]
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = app.EventService.UpdateEvent(r.Context(), id, evt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func handleDeleteEvent(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		idParam := params["id"]
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = app.EventService.DeleteEvent(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func handleGetEventList(app *app.App, period string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		dateParam := params["date"]
		date, err := time.Parse("2006-01-02", dateParam)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		list := []event.Event{}

		switch period {
		case "date":
			list, err = app.EventService.GetEventListByDate(r.Context(), date)
		case "week":
			list, err = app.EventService.GetEventListByWeek(r.Context(), date)
		case "month":
			list, err = app.EventService.GetEventListByDate(r.Context(), date)
		default:
			list, err = app.EventService.GetEventListByDate(r.Context(), date)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		writeJson(w, EventList{List: list})
	}
}

func handleGetEventByID(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		idParam := params["id"]
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		evt, err := app.EventService.GetEventByID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		writeJson(w, evt)
	}
}
