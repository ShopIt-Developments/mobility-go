package rating

import (
	"time"
	"io"
	"encoding/json"
	"appengine"
	"appengine/datastore"
	"errors"
)

type Rating struct {
	RatingId int64 `json:"rating_id" datastore:"-"`
	RaterId  string `json:"rater_id"`
	RatedId  string `json:"rated_id"`
	Date     time.Time `json:"date"`
	Rating   int8 `json:"rating"`
	Comment  string `json:"comment"`
}

func (rating *Rating) key(c appengine.Context) *datastore.Key {
	if rating.RaterId == "" {
		return datastore.NewIncompleteKey(c, "Rating", nil)
	}

	return datastore.NewKey(c, "Rating", "", rating.RatingId, nil)
}

func (rating *Rating) save(c appengine.Context) error {
	k, err := datastore.Put(c, rating.key(c), rating)

	if err != nil {
		return err
	}

	rating.RatingId = k.IntID()

	return nil
}

func GetOne(c appengine.Context, ratingId int64) (*Rating, error) {
	rating := Rating{RatingId: ratingId}

	k := rating.key(c)
	err := datastore.Get(c, k, &rating)

	if err != nil {
		return nil, err
	}

	rating.RatingId = k.IntID()

	return &rating, nil
}

func GetRatings(c appengine.Context, userId string) ([]Rating, error) {
	q := datastore.NewQuery("Rating").Filter("RatedId =", userId)

	var ratings []Rating

	keys, err := q.GetAll(c, &ratings)

	if err != nil {
		return nil, err
	}

	for i := 0; i < len(ratings); i++ {
		ratings[i].RatingId = keys[i].IntID()
	}

	return ratings, nil
}

func GetRated(c appengine.Context, userId string) ([]Rating, error) {
	q := datastore.NewQuery("Rating").Filter("RaterId =", userId)

	var ratings []Rating

	keys, err := q.GetAll(c, &ratings)

	if err != nil {
		return nil, err
	}

	for i := 0; i < len(ratings); i++ {
		ratings[i].RatingId = keys[i].IntID()
	}

	return ratings, nil
}

func New(c appengine.Context, r io.ReadCloser, raterId string) (*Rating, error) {
	rating := new(Rating)

	if err := json.NewDecoder(r).Decode(&rating); err != nil {
		return nil, err
	}

	if rating.Rating < 1 || rating.Rating > 5 {
		return nil, errors.New("Rating must be between 1 and 5")
	}

	rating.RaterId = raterId

	if err := rating.save(c); err != nil {
		return nil, err
	}

	return rating, nil
}

func Update(c appengine.Context, ratingId int64, r io.ReadCloser) (*Rating, error) {
	rating, err := GetOne(c, ratingId)

	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(r).Decode(&rating); err != nil {
		return nil, err
	}

	if err := rating.save(c); err != nil {
		return nil, err
	}

	return rating, nil
}

func Delete(c appengine.Context, ratingId int64) (*Rating, error) {
	rating, err := GetOne(c, ratingId)

	if err != nil {
		return nil, err
	}

	err = datastore.Delete(c, rating.key(c))

	if err != nil {
		return nil, err
	}

	return rating, nil
}