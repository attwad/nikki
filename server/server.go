package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/attwad/nikki/model"
	"github.com/attwad/nikki/store"
	"github.com/attwad/nikki/store/sqlite"
)

func InitMux(mux *http.ServeMux, prefix string, store *sqlite.SQLiteStore) {
	mux.Handle(fmt.Sprintf("GET %s/recent", prefix), &ListThingsHandler{Store: store})
	mux.Handle(fmt.Sprintf("GET %s/things/oftendone", prefix), &ListOftenDone{Store: store})
	mux.Handle(fmt.Sprintf("POST %s/things", prefix), &NewThingHandler{Store: store})
	mux.Handle(fmt.Sprintf("POST %s/things/{id}/delete", prefix), &DeleteThingHandler{Store: store})
}

func writeJSON(w http.ResponseWriter, v interface{}) error {
	j, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(j)
	return nil
}

// ListThingsHandler returns all things done this year.
type ListThingsHandler struct {
	Store *sqlite.SQLiteStore
}

func StoreThingToModelThing(thing store.Thing) model.Thing {
	return model.Thing{
		ID:   thing.ID,
		What: thing.What,
	}
}

func (h *ListThingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Query things.
	things, err := h.Store.Queries.ListThings(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Index things by date for easy lookup.
	thingsByDate := map[time.Time][]model.Thing{}
	for _, thing := range things {
		date := time.Date(int(thing.Year), time.Month(thing.Month), int(thing.Day), 0, 0, 0, 0, time.UTC)
		thingsByDate[date] = append(thingsByDate[date], StoreThingToModelThing(thing))
	}

	// Generate dates from Jan 1st to now and add things done in each of those.
	var dates []*model.Date
	now := time.Now()
	movingDate := time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	for {
		date := &model.Date{
			Year:  int64(movingDate.Year()),
			Month: int64(movingDate.Month()),
			Day:   int64(movingDate.Day()),
		}
		if things, ok := thingsByDate[movingDate]; ok {
			date.Things = append(date.Things, things...)
		}
		dates = append(dates, date)

		movingDate = movingDate.Add(24 * time.Hour)
		if movingDate.After(now) {
			break
		}
	}
	slices.Reverse(dates)
	if err := writeJSON(w, dates); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type ListOftenDone struct {
	Store *sqlite.SQLiteStore
}

func (h *ListOftenDone) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rows, err := h.Store.Queries.ListTopThings(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var things []string
	for _, thing := range rows {
		things = append(things, thing.What)
	}
	if err := writeJSON(w, things); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type NewThingHandler struct {
	Store *sqlite.SQLiteStore
}

func (h *NewThingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var thing model.NewThing
	if err := decoder.Decode(&thing); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	lastID, err := h.Store.Queries.InsertThing(r.Context(), store.InsertThingParams{
		Year:  thing.Year,
		Month: thing.Month,
		Day:   thing.Day,
		What:  thing.What,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/api/things/%d", lastID))
	w.WriteHeader(http.StatusCreated)
	if err := writeJSON(w, lastID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type DeleteThingHandler struct {
	Store *sqlite.SQLiteStore
}

func (h *DeleteThingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.Store.Queries.DeleteThing(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
