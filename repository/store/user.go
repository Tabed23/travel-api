package store

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/tabed23/travel-api/models"
	"github.com/tabed23/travel-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserStore struct {
	coll   mongo.Collection
	logger *slog.Logger
}

func NewUserStore(c mongo.Collection, l *slog.Logger) *UserStore {
	return &UserStore{coll: c, logger: l}
}

// Create User Document
func (u *UserStore) CreaterUser(usr models.User) (models.User, error) {
	u.logger.Info("repository", "CreaterUser", "User")

	usr.ID = primitive.NewObjectID()
	usr.CreatedAt = time.Now().UTC()
	usr.UpdatedAt = time.Now().UTC()
	hash, err := utils.EnscryptPassword(usr.Password)
	if err != nil {
		u.logger.Error("repository", "CreaterUser", err.Error())
		return models.User{}, err
	}
	usr.Password = hash
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	_, err = u.coll.InsertOne(ctx, &usr)
	if err != nil {
		u.logger.Error("repository", "CreaterUser", err.Error())

		return models.User{}, err
	}
	u.logger.Info("repository", "CreaterUser", fmt.Sprintf("User created %v", usr))

	return usr, nil
}

// Get All Users Document
func (u *UserStore) GetAll(page, limit int) ([]models.User, int, error) {
	u.logger.Info("repository", "GetAll", "User")

	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	users := []models.User{}
	skip := (page - 1) * limit

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip))
	cur, err := u.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		u.logger.Error("repository", "GetAll", err.Error())

		return nil, 0, err
	}
	for cur.Next(ctx) {
		var user models.User
		if err := cur.Decode(&user); err != nil {
			u.logger.Error("repository", "GetAll", err.Error())

			return nil, 0, err
		}
		users = append(users, user)
	}
	count, err := u.coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		u.logger.Error("repository", "GetAll", err.Error())
		return nil, 0, err
	}

	return users, int(count), nil
}

// Get Document
func (u UserStore) Get(Email string) (models.User, error) {
	u.logger.Info("repository", "Get", "User")

	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	var usr models.User
	if err := u.coll.FindOne(ctx, bson.M{"email": Email}).Decode(&usr); err != nil {
		u.logger.Error("repository", "Get", err.Error())

		return models.User{}, err
	}
	u.logger.Info("repository", "Get", fmt.Sprintf("User Found %v", usr))
	return usr, nil

}

// Delete User Document
func (u UserStore) Delete(Email string) (bool, error) {
	u.logger.Info("repository", "Delete", "User")
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	_, err := u.coll.DeleteOne(ctx, bson.M{"email": Email})
	if err != nil {
		u.logger.Info("repository", "Delete", err.Error())

		return false, err
	}

	u.logger.Info("repository", "Delete", fmt.Sprintf("User with %v Delete", Email))

	return true, nil
}

// Update the User Document
func (u UserStore) UpdateUser(email string, usr models.UserUpdate) (models.User, error) {
	u.logger.Info("repository", "Update", "User")

	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	var hashPass string
	if usr.Password != "" {
		hashPass, _ = utils.EnscryptPassword(usr.Password)
	}

	defer cancle()
	update := bson.M{
		"first_name": usr.FirstName,
		"last_name":  usr.Lastname,
		"Password":   hashPass,
		"role":       usr.Role,
		"photo":      usr.Photo,
		"updatedAt":  time.Now().UTC(),
	}
	_, err := u.coll.UpdateOne(ctx, bson.M{"email": email}, bson.M{"$set": update})
	if err != nil {
		u.logger.Error("repository", "UpdateUser", err.Error())
		return models.User{}, err
	}
	u.logger.Info("repository", "UpdateUser", fmt.Sprintf("user %v has been updated", email))

	return u.Get(email)

}

// Count User Document
func (u UserStore) CountUser() (int, error) {
	u.logger.Info("repository", "CountUser", "User")

	ctx, cancle := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancle()
	opts := options.Count().SetHint("_id_")
	count, err := u.coll.CountDocuments(ctx, bson.D{}, opts)
	if err != nil {
		u.logger.Error("repository", "CountUser", err.Error())

		return 0, err
	}

	u.logger.Info("repository", "CountUser", fmt.Sprintf("Total User Count %d", count))

	return int(count), nil
}

// Get User by Value Document
func (u *UserStore) GetUserByValue(ctx context.Context, value, data string) (models.User, error) {
	usr := models.User{}
	u.coll.FindOne(ctx, bson.M{value: data}).Decode(&u)
	return usr, nil
}

// Get User By Email Document
func (u *UserStore) GetUserByEmail(email string) (models.User, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	var user models.User
	if err := u.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		u.logger.Error("repository", "GetUserByEmail", err.Error())
		return models.User{}, err
	}
	u.logger.Info("repository", "GetUserByEmail", fmt.Sprintf("Found ther user with email: %v", email))
	return user, nil
}
