package vehicle

import (
    "model/sasa"
)

type Vehicles struct {
    Vehicles []Vehicle `json:"vehicles"`
    Buses []sasa.RealtimeBus `json:"buses"`
}
