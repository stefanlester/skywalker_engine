package data

import (
	"errors"
	"time"

	up "github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	// omitempty means that if the field is empty, it will not be inserted into the database
	ID        int       `db:"id, omitempty"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Active    int       `db:"user_active"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Token     Token     `db:"-"`
}

func (u *User) Table() string {
	return "users"
}

func (u *User) GetAll() ([]*User, error) {
	collection := upper.Collection(u.Table())

	var all []*User

	res := collection.Find().OrderBy("last_name")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// GetByEmail gets one user, email
func (u *User) GetByEmail(email string) (*User, error) {
	var theUser User
	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"email": email})
	err := res.One(&theUser)
	if err != nil {
		return nil, err
	}

	var token Token
	collection = upper.Collection(token.Table())

	res = collection.Find(up.Cond{"user_id": theUser.ID, "expires >": time.Now()}).OrderBy("created_at desc")
	err = res.One(&token)

	if err != nil {
		if err != up.ErrNilRecord && err != up.ErrNoMoreRows {
			return nil, err
		}
	}

	return &theUser, nil
}

// gets a user by ID
func (u *User) GetByID(id int) (*User, error) {
	var theUser User
	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"id =": id})

	err := res.One(&theUser)
	if err != nil {
		return nil, err
	}
	var token Token
	collection = upper.Collection(token.Table())

	res = collection.Find(up.Cond{"user_id": theUser.ID, "expires >": time.Now()}).OrderBy("created_at desc")
	err = res.One(&token)

	if err != nil {
		if err != up.ErrNilRecord && err != up.ErrNoMoreRows {
			return nil, err
		}
	}

	theUser.Token = token

	return &theUser, nil
}

func (u *User) Update(theUser User) error {
	theUser.UpdatedAt = time.Now()

	collection := upper.Collection(u.Table())

	res := collection.Find(theUser.ID)
	err := res.Update(&theUser)
	if err != nil {
		return err

	}

	return nil
}

func (u *User) Delete(id int) error {
	collection := upper.Collection(u.Table())

	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err

	}

	return nil
}

func (u *User) Insert(theUser User) (int, error) {
	newHash, err := bcrypt.GenerateFromPassword([]byte(theUser.Password), 12)
	if err != nil {
		return 0, err
	}

	theUser.CreatedAt = time.Now()
	theUser.UpdatedAt = time.Now()
	theUser.Password = string(newHash)

	collection := upper.Collection(u.Table())
	res, err := collection.Insert(theUser)

	if err != nil {
		return 0, err
	}

	id := getInsertID(res.ID())

	return id, nil
}

func (u *User) ResetPassword(id int, password string) error {
	newHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	theUser, err := u.GetByID(id)

	if err != nil {
		return err
	}

	u.Password = string(newHash)

	err = theUser.Update(*u)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) PasswordMatches(plaintText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plaintText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}

	}

	return true, nil

}
