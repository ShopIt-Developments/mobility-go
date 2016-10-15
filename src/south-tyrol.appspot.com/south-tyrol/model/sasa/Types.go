package sasa

type RealtimeBuses struct {
    Buses []RealtimeBus `json:"buses"`
}

type RealtimeBus struct {
    LineName    string `json:"name"`
    BusStop     string `json:"description"`
    HydrogenBus bool `json:"hydrogen_bus"`
    TripId      int `json:"id"`
    Latitude    float64 `json:"lat"`
    Longitude   float64 `json:"lng"`
}

type Buses struct {
    Type     string `json:"type"`
    Features []Bus `json:"features"`
}

type Bus struct {
    Type       string `json:"type"`
    Geometry   Geometry `json:"geometry"`
    Properties Properties `json:"properties"`
}

type Geometry struct {
    Type        string `json:"type"`
    Coordinates []float64 `json:"coordinates"`
}

type Properties struct {
    TripId      int `json:"frt_fid"`
    DelaySec    int `json:"delay_sec"`
    VehicleId   string `json:"vehiclecode"`
    LineId      int `json:"li_nr"`
    LineName    string `json:"lidname"`
    BusStop     int `json:"ort_nr"`
    BusStopName string `json:"ort_name"`
}