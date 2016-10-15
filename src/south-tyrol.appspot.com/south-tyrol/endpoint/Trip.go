package endpoint

import (
    "github.com/julienschmidt/httprouter"
    "net/http"
    "issue"
    "appengine"
    "network"
    "model/trip"
    "encoding/json"
    "model/user"
)

const POINTS_PER_TRIP

type Trip struct {
    Router *httprouter.Router
}

func (*Trip) New(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    travel, err := trip.New(appengine.NewContext(r), r.Body, network.Authorization(w, r))
    issue.Handle(w, err, http.StatusBadRequest)

    duration := travel.Arrival.Sub(travel.Departure)
    points := int(duration.Minutes() * POINTS_PER_TRIP)

    data, err := json.Marshal(trip.TripResponse{
        Points: int(points),
        Duration: duration.String(),
    })

    issue.Handle(w, err, http.StatusInternalServerError)
    user.AddPoints(appengine.NewContext(r), travel.UserId, int64(points))

    w.Write(data)
}