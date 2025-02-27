package main

import (
	"angya-backend/internal/photoApplication"
	"angya-backend/internal/poiApplication"
	"angya-backend/internal/spotApplication"
	"angya-backend/pkg/databases"
	"angya-backend/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	_ "net/http/pprof"
)

type (
	SpotHandler        struct{}
	PhotoHandler       struct{}
	PhotoDetailHandler struct{}
	PoiHandler         struct{}
)

func (h *SpotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	utils.EnhanceResponseWriter(&w)
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body) //ReadAllでResponse Bodyを読み切る
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	db := databases.NewLocalPostgres()
	spotUsecase := spotApplication.NewUsecase(db)

	d, _ := db.DB()
	defer d.Close()

	switch r.Method {
	case "POST":
		spot, err := spotUsecase.Register(r.Context(), b)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		fmt.Fprint(w, string(func() (b []byte) { j, _ := json.Marshal(spot); return j }()))
	case "GET":
		spots, err := spotUsecase.List(r.Context())
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		fmt.Fprint(w, string(func() (b []byte) { j, _ := json.Marshal(spots); return j }()))
	}
}

func (h *PhotoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	utils.EnhanceResponseWriter(&w)
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body) //ReadAllでResponse Bodyを読み切る
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	db := databases.NewLocalPostgres()
	photoUsecase := photoApplication.NewUsecase(db)

	d, _ := db.DB()
	defer d.Close()

	switch r.Method {
	case "POST":
		photo, err := photoUsecase.Register(r.Context(), b)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		fmt.Fprint(w, string(func() (b []byte) { j, _ := json.Marshal(photo); return j }()))
	case "GET":
		photo, err := photoUsecase.List(r.Context())
		if err != nil {
			fmt.Fprint(w, err.Error())
		}
		fmt.Fprint(w, string(func() (b []byte) { j, _ := json.Marshal(photo); return j }()))
	}
}

func (h *PhotoDetailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	utils.EnhanceResponseWriter(&w)
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body) //ReadAllでResponse Bodyを読み切る
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	id := utils.GetIdFromPath(r.URL.Path)

	db := databases.NewLocalPostgres()
	photoUsecase := photoApplication.NewUsecase(db)

	d, _ := db.DB()
	defer d.Close()

	switch r.Method {
	case "PATCH":
		photo, err := photoUsecase.Update(r.Context(), id, b)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		fmt.Fprint(w, string(func() (b []byte) { j, _ := json.Marshal(photo); return j }()))
	}
}

func (h *PoiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	utils.EnhanceResponseWriter(&w)
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body) //ReadAllでResponse Bodyを読み切る
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	db := databases.NewLocalPostgres()
	poiUsecase := poiApplication.NewUsecase(db)

	d, _ := db.DB()
	defer d.Close()

	switch r.Method {
	case "POST":
		poi, err := poiUsecase.Migrate(r.Context(), b)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		fmt.Fprint(w, string(func() (b []byte) { j, _ := json.Marshal(poi); return j }()))
	case "GET":
		pois, err := poiUsecase.List(r.Context())
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		fmt.Fprint(w, string(func() (b []byte) { j, _ := json.Marshal(pois); return j }()))
	}
}

func main() {
	s, p, p2, o := SpotHandler{}, PhotoHandler{}, PhotoDetailHandler{}, PoiHandler{}

	server := http.Server{
		Addr:    ":8088",
		Handler: nil,
	}

	http.Handle("/api/spots", &s)
	http.Handle("/api/photos", &p)
	http.Handle("/api/photos/", &p2)
	http.Handle("/api/pois", &o)

	server.ListenAndServe()
}
