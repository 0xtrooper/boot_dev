package db

import (
	"encoding/json"
	"errors"
	"os"
	"server_course/entities"
	"sort"
	"strings"
	"sync"
)

var (
	ErrDoesNotExist = errors.New("does not exist")
)

type DBStructure struct {
	Chirps     map[int]entities.Chirp `json:"chirps"`
	Users      map[int]entities.User  `json:"users"`
	ChirpIndex int                    `json:"chirp_index"`
	UserIndex  int                    `json:"user_index"`
}

type DB struct {
	store DBStructure
	path  string
	mux   *sync.RWMutex
}

func NewDB(p string) (*DB, error) {
	db := &DB{
		store: DBStructure{
			Chirps:     make(map[int]entities.Chirp),
			Users:      make(map[int]entities.User),
			ChirpIndex: 1,
			UserIndex:  1,
		},
		path: p + "/database.json",
		mux:  &sync.RWMutex{},
	}

	return db, db.loadDB()
}

func (db *DB) StoreChirp(c entities.Chirp) (entities.Chirp, error) {
	db.mux.Lock()
	c.ID = db.store.ChirpIndex // idk
	db.store.Chirps[db.store.ChirpIndex] = c
	db.store.ChirpIndex++
	db.mux.Unlock() // unlock manual cause writeDB relocks

	return c, db.writeDB()
}

func (db *DB) StoreUser(u entities.User) (entities.User, error) {
	db.mux.Lock()
	u.ID = db.store.UserIndex // idk
	encryptedUser, err := u.EncryptPassword()
	if err != nil {
		return entities.User{}, err
	}
	db.store.Users[db.store.UserIndex] = encryptedUser
	db.store.UserIndex++
	db.mux.Unlock() // unlock manual cause writeDB relocks

	return u, db.writeDB()
}

func (db *DB) DeleteChirp(chirpID int) error {
	db.mux.Lock()
	delete(db.store.Chirps, chirpID)
	db.mux.Unlock() // unlock manual cause writeDB relocks

	return db.writeDB()
}

func (db *DB) GetChirp(chirpID int) (entities.Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	if c, exits := db.store.Chirps[chirpID]; exits {
		return c, nil
	} else {
		return entities.Chirp{}, ErrDoesNotExist
	}
}

func (db *DB) GetUser(userID int) (entities.User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	if c, exits := db.store.Users[userID]; exits {
		return c, nil
	} else {
		return entities.User{}, ErrDoesNotExist
	}
}

func (db *DB) GetUserByEmail(requestedEmail string) (entities.User, error) {
	allUsers, err := db.GetUsers()
	if err != nil {
		return entities.User{}, err
	}

	for _, storedUser := range allUsers {
		if strings.EqualFold(storedUser.Email, requestedEmail) {
			return storedUser, nil
		}
	}

	return entities.User{}, ErrDoesNotExist
}

func (db *DB) GetChirps() (map[int]entities.Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	return db.store.Chirps, nil
}

func (db *DB) GetUsers() (map[int]entities.User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	return db.store.Users, nil
}

func (db *DB) GetChirpsSlice() ([]entities.Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	chirps := make([]entities.Chirp, len(db.store.Chirps))
	i := 0
	for _, c := range db.store.Chirps {
		chirps[i] = c
		i++
	}

	sort.Slice(chirps, func(i, j int) bool { return chirps[i].ID < chirps[j].ID })

	return chirps, nil
}

func (db *DB) GetUsersSlice() ([]entities.User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	users := make([]entities.User, len(db.store.Users))
	i := 0
	for _, u := range db.store.Users {
		users[i] = u
		i++
	}

	sort.Slice(users, func(i, j int) bool { return users[i].ID < users[j].ID })

	return users, nil
}

func (db *DB) UpdateUser(newUser entities.User) (entities.User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	oldUser, exits := db.store.Users[newUser.ID]
	if !exits {
		return entities.User{}, ErrDoesNotExist
	}

	if newUser.Password != "" {
		encryptedUser, err := newUser.EncryptPassword()
		if err != nil {
			return entities.User{}, err
		}
		oldUser.Password = encryptedUser.Password
	}

	oldUser.Email = newUser.Email
	oldUser.Token = newUser.Token
	oldUser.ExpiresInSeconds = newUser.ExpiresInSeconds
	oldUser.RefreshToken = newUser.RefreshToken
	oldUser.RefreshExpiresInSeconds = newUser.RefreshExpiresInSeconds
	db.store.Users[newUser.ID] = oldUser
	return oldUser, nil
}

func (db *DB) UpdateUserTokens(newUser entities.User) (entities.User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	oldUser, exits := db.store.Users[newUser.ID]
	if !exits {
		return entities.User{}, ErrDoesNotExist
	}

	oldUser.Token = newUser.Token
	oldUser.ExpiresInSeconds = newUser.ExpiresInSeconds
	oldUser.RefreshToken = newUser.RefreshToken
	oldUser.RefreshExpiresInSeconds = newUser.RefreshExpiresInSeconds
	db.store.Users[newUser.ID] = oldUser
	return oldUser, nil
}

func (db *DB) UpdateUserEmailAndPassword(newUser entities.User) (entities.User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	oldUser, exits := db.store.Users[newUser.ID]
	if !exits {
		return entities.User{}, ErrDoesNotExist
	}

	encryptedUser, err := newUser.EncryptPassword()
	if err != nil {
		return entities.User{}, err
	}
	oldUser.Email = newUser.Email
	oldUser.Password = encryptedUser.Password
	return oldUser, nil
}

func (db *DB) UpdateUserRedStatus(userID int, status bool) (entities.User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	user, exits := db.store.Users[userID]
	if !exits {
		return entities.User{}, ErrDoesNotExist
	}

	user.IsChirpyRed = status
	db.store.Users[userID] = user
	return db.store.Users[userID], nil
}

func (db *DB) Reset() error {
	db.mux.Lock()
	defer db.mux.Unlock()

	err := os.Remove(db.path)
	if err != nil {
		return err
	}
	clear(db.store.Chirps)
	return nil
}

func (db *DB) loadDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := os.ReadFile(db.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &db.store)
}

func (db *DB) writeDB() error {
	db.mux.RLock()
	defer db.mux.RUnlock()

	data, err := json.Marshal(db.store)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, 0644)
	return err
}
