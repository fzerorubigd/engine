package userpb

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/fzerorubigd/engine/pkg/assert"
	"github.com/fzerorubigd/engine/pkg/kv"
	"github.com/fzerorubigd/engine/pkg/log"
	"github.com/fzerorubigd/engine/pkg/random"
)

// From the bcrypt package
const (
	minHashSize  = 59
	noPassString = "NO" // Size must be less than 6 character
)

//  TODO: NEEDS COMMENT INFO
var (
	isBcrypt = regexp.MustCompile(`^\$[^$]+\$[0-9]+\$`)
)

func (m *User) cryptPassword() {
	// TODO : Watch it if this creepy code is dangerous :)
	if (len(m.Password) < minHashSize || !isBcrypt.MatchString(m.Password)) && m.Password != noPassString {
		p, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
		assert.Nil(err)
		m.Password = string(p)
	}
}

// PreInsert the user on create
func (m *User) PreInsert() {
	m.cryptPassword()
}

// PreUpdate the user on update
func (m *User) PreUpdate() {
	m.cryptPassword()
}

// VerifyPassword try to verify password for given hash
func (m *User) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password)) == nil
}

// FindUserByEmailPassword try to login user with username and password
func (m *Manager) FindUserByEmailPassword(ctx context.Context, email, password string) (*User, error) {
	u, err := m.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if u.Status == UserStatus_USER_STATUS_BANNED {
		return nil, errors.New("sorry, but you are banned")
	}

	if u.VerifyPassword(password) {
		return u, nil
	}

	return nil, errors.New("user not found or wrong password")
}

// FindUserByEmail is a function to find user based on app
func (m *Manager) FindUserByEmail(ctx context.Context, e string) (*User, error) {
	q := fmt.Sprintf(
		"SELECT %s FROM %s WHERE email = $1 ",
		strings.Join(m.getUserFields(), ","),
		UserTableFull,
	)

	r := m.GetDbMap().QueryRowxContext(ctx, q, e)

	return m.scanUser(r)
}

// RegisterUser is to register new user
func (m *Manager) RegisterUser(ctx context.Context, email, name, pass string) (*User, error) {
	u := User{
		Email:       email,
		DisplayName: name,
		Password:    pass,
		Status:      UserStatus_USER_STATUS_REGISTERED,
	}

	if err := m.CreateUser(ctx, &u); err != nil {
		return nil, errors.Wrap(err, "already registered")
	}

	return &u, nil
}

// CreateToken TODO: NEEDS COMMENT INFO
func (m *Manager) CreateToken(ctx context.Context, u *User, d time.Duration) string {
	t := <-random.ID
	m.UpdateToken(ctx, u, d, t)
	return t
}

// UpdateToken try to update/create token
func (m *Manager) UpdateToken(_ context.Context, u *User, d time.Duration, token string) {
	v, err := proto.Marshal(u)
	assert.Nil(err)
	kv.MustStoreKey(token, string(v), d)
}

// FindUserByIndirectToken TODO: NEEDS COMMENT INFO
func (m *Manager) FindUserByIndirectToken(ctx context.Context, token string) (*User, error) {
	t, err := kv.FetchKey(token)
	if err != nil {
		return nil, err
	}
	var u User
	// Invalid data is a bug
	assert.Nil(proto.Unmarshal([]byte(t), &u))

	return &u, nil
}

// DeleteToken TODO: NEEDS COMMENT INFO
func (m *Manager) DeleteToken(_ context.Context, token string) {
	kv.MustDeleteKey(token)
}

// ChangePassword TODO: NEEDS COMMENT INFO
func (m *Manager) ChangePassword(ctx context.Context, u *User, newPassword string) error {
	u.Password = newPassword
	return m.UpdateUser(ctx, u)
}

// CreateForgottenToken return a forgotten token, also return the age of already generated token
// TODO: rate limit
func (m *Manager) CreateForgottenToken(ctx context.Context, u *User) (string, time.Duration, error) {
	key := fmt.Sprintf("forgotten_%d", u.Id)
	v, err := kv.FetchKey(key)
	if err != nil {
		v = <-random.ID
		assert.Nil(kv.StoreKey(key, v, 24*time.Hour))
	}

	ttl := kv.MustTTLKey(key)
	return v, ttl, nil
}

// VerifyForgottenToken try to verify token and remove it after successful verify
func (m *Manager) VerifyForgottenToken(ctx context.Context, u *User, token string) error {
	key := fmt.Sprintf("forgotten_%d", u.Id)
	v, err := kv.FetchKey(key)
	if err != nil {
		return errors.Wrap(err, "not found")
	}

	if v != token {
		return errors.New("invalid token")
	}

	if err := kv.DeleteKey(key); err != nil {
		log.Error("Could not delete the key", log.Err(err))
	}

	return nil
}
