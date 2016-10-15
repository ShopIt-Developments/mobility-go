package sasa

import (
    "math"
    "net/http"
    "appengine/urlfetch"
    "appengine"
)

func ReadJsonFromUrl(w http.ResponseWriter, r *http.Request, url string) *http.Response {
    context := appengine.NewContext(r)
    client := urlfetch.Client(context)
    response, err := client.Get(url)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return nil
    }

    return response;
}

func ToWgs84(coordinates []float64, zone int) (float64, float64) {
    d := 0.99960000000000004
    d1 := 6378137.0
    d2 := 0.0066943799999999998
    d4 := (1 - math.Sqrt(1 - d2)) / (1 + math.Sqrt(1 - d2))
    d3 := d2 / (1 - d2)
    d12 := coordinates[1] / d / (d1 * (1 - d2 / 4 - (3 * d2 * d2) / 64 - (5 * math.Pow(d2, 3)) / 256))
    d14 := d12 + ((3 * d4) / 2 - (27 * math.Pow(d4, 3)) / 32) * math.Sin(2 * d12) + ((21 * d4 * d4) / 16 - (55 * math.Pow(d4, 4)) / 32) * math.Sin(4 * d12) + ((151 * math.Pow(d4, 3)) / 96) * math.Sin(6 * d12)
    d5 := d1 / math.Sqrt(1 - d2 * math.Sin(d14) * math.Sin(d14))
    d6 := math.Tan(d14) * math.Tan(d14)
    d7 := d3 * math.Cos(d14) * math.Cos(d14)
    d8 := (d1 * (1 - d2)) / math.Pow(1 - d2 * math.Sin(d14) * math.Sin(d14), 1.5)
    d9 := (coordinates[0] - 500000) / (d5 * d)
    d17 := rad2deg(d14 - ((d5 * math.Tan(d14)) / d8) * (((d9 * d9) / 2 - (((5 + 3 * d6 + 10 * d7) - 4 * d7 * d7 - 9 * d3) * math.Pow(d9, 4)) / 24) + (((61 + 90 * d6 + 298 * d7 + 45 * d6 * d6) - 252 * d3 - 3 * d7 * d7) * math.Pow(d9, 6)) / 720))
    d18 := float64(((zone - 1) * 6 - 180) + 3) + rad2deg(((d9 - ((1 + 2 * d6 + d7) * math.Pow(d9, 3)) / 6) + (((((5 - 2 * d7) + 28 * d6) - 3 * d7 * d7) + 8 * d3 + 24 * d6 * d6) * math.Pow(d9, 5)) / 120) / math.Cos(d14))

    return toFixed(d17, 5), toFixed(d18, 5)
}

func rad2deg(radians float64) float64 {
    return radians * 180 / math.Pi
}

func toFixed(num float64, precision int) float64 {
    output := math.Pow(10, float64(precision))
    return float64(round(num * output)) / output
}

func round(num float64) int {
    return int(num + math.Copysign(0.5, num))
}