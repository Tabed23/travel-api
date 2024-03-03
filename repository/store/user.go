package store

import (
	"context"
	"time"

	"github.com/tabed23/travel-api/models"
	"github.com/tabed23/travel-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserStore struct {
	coll mongo.Collection
}

func NewUserStore(c mongo.Collection) *UserStore {
	return &UserStore{coll: c}
}
func (u *UserStore) CreaterUser(usr models.User) (models.User, error) {
	usr.ID = primitive.NewObjectID()
	usr.CreatedAt = time.Now().UTC()
	usr.UpdatedAt = time.Now().UTC()
	hash, err := utils.EnscryptPassword(usr.Password)
	if err != nil {
		return models.User{}, err
	}
	usr.Password = hash
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	_, err = u.coll.InsertOne(ctx, &usr)
	if err != nil {
		return models.User{}, err
	}
	return usr, nil
}

func (u *UserStore) GetAll(page, limit int) ([]models.User, int, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	users := []models.User{}
	skip := (page - 1) * limit

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip))
	cur, err := u.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	for cur.Next(ctx) {
		var user models.User
		if err := cur.Decode(&user); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}
	count, err := u.coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}
	return users, int(count), nil
}

func (u UserStore) Get(Email string) (models.User, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	var usr models.User
	if err := u.coll.FindOne(ctx, bson.M{"email": Email}).Decode(&usr); err != nil {
		return models.User{}, err
	}
	return usr, nil

}
func (u UserStore) Delete(Email string) (bool, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	_, err := u.coll.DeleteOne(ctx, bson.M{"email": Email})
	if err != nil {
		return false, err
	}
	return true, nil
}
func (u UserStore) UpdateUser(email string, usr models.User) (models.User, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	update := bson.M{
		"first_name": usr.FirstName,
		"last_name":  usr.Lastname,
		"Password":   usr.Password,
		"role":       usr.Role,
		"updatedAt":  time.Now().UTC(),
	}
	_, err := u.coll.UpdateOne(ctx, bson.M{"email": email}, bson.M{"$set": update})
	if err != nil {
		return models.User{}, err
	}

	return u.Get(email)

}

func (u UserStore) CountUser() (int, error) {
	ctx, cancle := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancle()
	opts := options.Count().SetHint("_id_")
	count, err := u.coll.CountDocuments(ctx, bson.D{}, opts)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (u *UserStore) GetUserByValue(ctx context.Context, value, data string) (models.User, error) {
	usr := models.User{}
	u.coll.FindOne(ctx, bson.M{value: data}).Decode(&u)
	return usr, nil
}
