package userpb

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gogo/protobuf/types"

	"github.com/fzerorubigd/balloon/pkg/assert"
	"github.com/fzerorubigd/balloon/pkg/random"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// From the bcrypt package
const (
	minHashSize  = 59
	noPassString = "NO" // Size must be less than 6 character
)

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

// Initialize the user on create
func (m *User) PreInsert() {
	m.LastLogin, _ = types.TimestampProto(time.Now())
	m.cryptPassword()
}

// Initialize the user on update
func (m *User) PreUpdate() {
	m.cryptPassword()
}

// VerifyPassword try to verify password for given hash
func (m *User) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password)) == nil
}

// HasPassword check if user set password or not
func (m *User) HasPassword() bool {
	return m.Password != noPassString
}

// LoginUserByPassword try to login user with username and password
func (m *Manager) LoginUserByPassword(ctx context.Context, email, password string) (*User, error) {
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
func (m *Manager) RegisterUser(ctx context.Context, email, pass string) (*User, error) {
	u := User{
		Email:    email,
		Password: pass,
		Token:    <-random.ID,
		Status:   UserStatus_USER_STATUS_REGISTERED,
	}

	if err := m.CreateUser(ctx, &u); err != nil {
		return nil, errors.Wrap(err, "Already registered")
	}

	return &u, nil
}

func (m *Manager) FindUserByToken(ctx context.Context, token string) (*User, error) {
	q := fmt.Sprintf(
		"SELECT %s FROM %s WHERE token = $1 ",
		strings.Join(m.getUserFields(), ","),
		UserTableFull,
	)

	r := m.GetDbMap().QueryRowxContext(ctx, q, token)

	return m.scanUser(r)
}
