package purchase

import (
	"net/http"
	"github.com/julienschmidt/httprouter"

	"endpoint"
    "model/user"
    "appengine"
    "issue"
    "encoding/json"
)

func init() {
	router := httprouter.New()

	v := endpoint.Vehicle{Router: router}
	v.Router.GET("/mobility/bus/:id", v.GetBus)
	v.Router.GET("/mobility/vehicles/available", v.GetAvailable)
	v.Router.GET("/mobility/vehicles/booked", v.GetBooked)
	v.Router.GET("/mobility/vehicles/my", v.GetMy)
	v.Router.GET("/mobility/vehicles", v.GetAll)
	v.Router.GET("/mobility/vehicle/:vehicle_id", v.GetOne)
	v.Router.POST("/mobility/vehicle", v.New)
	v.Router.PUT("/mobility/vehicles/:vehicle_id", v.Update)
	v.Router.DELETE("/mobility/vehicle/:vehicle_id", v.Delete)

	o := endpoint.Order{Router: router}
	o.Router.POST("/mobility/order/:vehicle_id", o.New)
	o.Router.DELETE("/mobility/order/:order_id", o.Delete)

	u := endpoint.User{Router: router}
	u.Router.GET("/mobility/user", u.Get)
	u.Router.POST("/mobility/user", u.Add)
	u.Router.PUT("/mobility/user", u.Update)
	u.Router.DELETE("/mobility/user", u.Delete)

	u.Router.GET("/mobility/points", u.GetPoints)
	u.Router.POST("/mobility/points/:points", u.AddPoints)

	u.Router.PATCH("/mobility/token/:token", u.SetToken)

	r := endpoint.Rating{Router: router}
	r.Router.GET("/mobility/ratings", r.GetRatings)
	r.Router.GET("/mobility/ratings/:rating_id", r.GetOne)
	r.Router.GET("/mobility/rated", r.GetRated)
	r.Router.POST("/mobility/ratings", r.Add)
	r.Router.PUT("/mobility/ratings/:rating_id", r.Update)
	r.Router.DELETE("/mobility/ratings/:rating_id", r.Delete)

	t := endpoint.Trip{Router: router}
	t.Router.POST("/mobility/trip", t.New)

	p := endpoint.Payment{Router: router}
	p.Router.POST("/mobility/payment/scan/:order_id", p.Scan)
	p.Router.POST("/mobility/payment/accept/:order_id", p.Accept)
	p.Router.POST("/mobility/payment/notify/:order_id", p.Notify)

	router.GET("/mobility/leaderboard", leaderboard)

	http.Handle("/", router)
}

func leaderboard(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    users, err := user.GetAll(appengine.NewContext(r))
    issue.Handle(w, err, http.StatusBadRequest)

    data, err := json.Marshal(&users)
    issue.Handle(w, err, http.StatusInternalServerError)

    w.Write(data)
}