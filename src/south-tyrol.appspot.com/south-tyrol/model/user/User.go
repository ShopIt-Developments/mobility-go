package user

import (
    "appengine"
    "appengine/datastore"
    "io"
    "encoding/json"
    "time"
)

type User struct {
    UserID          string `json:"user_id" datastore:"-"`
    Name            string `json:"name"`
    Address         string `json:"address"`
    Telephone       string `json:"telephone"`
    Points          int64 `json:"points"`
    AverageRating   float32 `json:"average_rating"`
    RatingsCount    int64 `json:"ratings_count"`
    DrivenTime      time.Duration `json:"driven_time"`
    Emissions       int64 `json:"emissions" datastore:"-"`
    Token           string `json:"token"`
    OfferedVehicles int64 `json:"offered_vehicles"`
    UsedVehicles    int64 `json:"used_vehicles"`
}

func (user *User) key(c appengine.Context) *datastore.Key {
    return datastore.NewKey(c, "User", user.UserID, 0, nil)
}

func (user *User) save(c appengine.Context) error {
    k, err := datastore.Put(c, user.key(c), user)

    if err != nil {
        return err
    }

    user.UserID = k.StringID()

    return nil
}

func New(c appengine.Context, r io.ReadCloser, userId string) (*User, error) {
    user := new(User)

    if err := json.NewDecoder(r).Decode(&user); err != nil {
        return nil, err
    }

    user.UserID = userId
    user.Points = 0
    user.AverageRating = 0
    user.RatingsCount = 0

    if err := user.save(c); err != nil {
        return nil, err
    }

    return user, nil
}

func AddPoints(c appengine.Context, userId string, points int64) (*Points, error) {
    user, err := Get(c, userId)

    if err != nil {
        return nil, err
    }

    user.Points += points

    user.save(c)

    if err != nil {
        return nil, err
    }

    return &Points{Points: user.Points}, nil
}

func AddDuration(c appengine.Context, userId string, duration time.Duration) error {
    user, err := Get(c, userId)

    if err != nil {
        return err
    }

    user.DrivenTime = user.DrivenTime + duration

    user.save(c)

    return err
}

func GetPoints(c appengine.Context, userId string) (*Points, error) {
    user, err := Get(c, userId)

    if err != nil {
        return nil, err
    }

    return &Points{Points: user.Points}, nil
}

func Get(c appengine.Context, userId string) (*User, error) {
    user := User{UserID: userId}

    k := user.key(c)
    err := datastore.Get(c, k, &user)

    if err != nil {
        return nil, err
    }

    user.UserID = k.StringID()
    user.Emissions = int64(float64(user.DrivenTime.Minutes()) * float64(106))

    return &user, nil
}

func GetAll(c appengine.Context) ([]User, error) {
    q := datastore.NewQuery("User")

    users := []User{}
    keys, err := q.GetAll(c, &users)

    if err != nil {
        return nil, err
    }

    for i := 0; i < len(users); i++ {
        users[i].UserID = keys[i].StringID()
    }

    return users, nil
}

func Update(c appengine.Context, userID string, r io.ReadCloser) (*User, error) {
    user, err := Get(c, userID)

    if err != nil {
        return nil, err
    }

    if err := json.NewDecoder(r).Decode(&user); err != nil {
        return nil, err
    }

    if err := user.save(c); err != nil {
        return nil, err
    }

    return user, nil
}

func Delete(c appengine.Context, userID string) (*User, error) {
    user, err := Get(c, userID)

    if err != nil {
        return nil, err
    }

    err = datastore.Delete(c, user.key(c))

    if err != nil {
        return nil, err
    }

    return user, nil
}

func NewRating(c appengine.Context, userID string, rating int8) error {
    user, err := Get(c, userID)

    if err != nil {
        return err
    }

    user.AverageRating = (user.AverageRating * float32(user.RatingsCount) + float32(rating)) / float32((user.RatingsCount + 1))
    user.RatingsCount++

    err = user.save(c)

    return err
}

func UpdateRating(c appengine.Context, userID string, oldRating int8, rating int8) error {
    user, err := Get(c, userID)

    if err != nil {
        return err
    }

    user.AverageRating = (user.AverageRating * float32(user.RatingsCount) + float32(rating) - float32(oldRating)) / float32((user.RatingsCount))

    err = user.save(c)

    return err
}

func DeleteRating(c appengine.Context, userID string, rating int8) error {
    user, err := Get(c, userID)

    if err != nil {
        return err
    }

    user.AverageRating = (user.AverageRating * float32(user.RatingsCount) - float32(rating) / float32((user.RatingsCount - 1)))
    user.RatingsCount--

    err = user.save(c)

    return err
}

func SetToken(c appengine.Context, userId string, token string) error {
    user, err := Get(c, userId)

    if err != nil {
        return err
    }

    user.Token = token

    return user.save(c)

    return err
}

func AddOfferedVehicle(c appengine.Context, userId string) error {
    user, err := Get(c, userId)

    if err != nil {
        return err
    }

    user.OfferedVehicles++

    return user.save(c)

}

func AddUsedVehicle(c appengine.Context, userId string) error {
    user, err := Get(c, userId)

    if err != nil {
        return err
    }

    user.UsedVehicles++

    return user.save(c)

}
