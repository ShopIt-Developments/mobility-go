package user

import (
	"appengine"
	"appengine/datastore"
	"io"
	"encoding/json"
)

type User struct {
	UserID        string `json:"user_id" datastore:"-"`
	Name          string `json:"name"`
	Address       string `json:"address"`
	Telephone     string `json:"telephone"`
	AverageRating float32 `json:"average_rating"`
	RatingsCount  int64`json:"ratings_count"`
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

func New(c appengine.Context, r io.ReadCloser, userID string) (*User, error) {
	user := new(User)

	if err := json.NewDecoder(r).Decode(&user); err != nil {
		return nil, err
	}

	user.UserID = userID

	if err := user.save(c); err != nil {
		return nil, err
	}

	return user, nil
}

func Get(c appengine.Context, userId string) (*User, error) {
	user := User{UserID: userId}

	k := user.key(c)
	err := datastore.Get(c, k, &user)

	if err != nil {
		return nil, err
	}

	user.UserID = k.StringID()

	return &user, nil
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

