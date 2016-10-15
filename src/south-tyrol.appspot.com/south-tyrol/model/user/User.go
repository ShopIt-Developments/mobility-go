package user

import (
	"appengine"
	"appengine/datastore"
	"io"
	"encoding/json"
)

type User struct {
	UserID    string `json:"user_id" datastore:"-"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Telephone string `json:"telephone"`
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

func Get(c appengine.Context, carId string) (*User, error) {
	user := User{UserID: carId}

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

func Delete(c appengine.Context, id string) (*User, error) {
	user, err := Get(c, id)

	if err != nil {
		return nil, err
	}

	err = datastore.Delete(c, user.key(c))

	if err != nil {
		return nil, err
	}

	return user, nil
}