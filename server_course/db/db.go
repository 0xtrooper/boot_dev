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
	Chirps map[int]entities.Chirp `json:"chirps"`
	Users  map[int]entities.User  `json:"users"`
	Index  int                    `json:"index"`
}

type DB struct {
	store DBStructure
	path  string
	mux   *sync.RWMutex
}

func NewDB(p string) (*DB, error) {
	db := &DB{
		store: DBStructure{
			Chirps: make(map[int]entities.Chirp),
			Users:  make(map[int]entities.User),
			Index:  1,
		},
		path: p + "/database.json",
		mux:  &sync.RWMutex{},
	}

	return db, db.loadDB()
}

func (db *DB) StoreChirp(c entities.Chirp) (entities.Chirp, error) {
	db.mux.Lock()
	c.ID = db.store.Index // idk
	db.store.Chirps[db.store.Index] = c
	db.store.Index++
	db.mux.Unlock() // unlock manual cause writeDB relocks

	return c, db.writeDB()
}

func (db *DB) StoreUser(u entities.User) (entities.User, error) {
	db.mux.Lock()
	u.ID = db.store.Index // idk
	encryptedUser, err := u.EncryptPassword()
	if err != nil {
		return entities.User{}, err
	}
	db.store.Users[db.store.Index] = encryptedUser
	db.store.Index++
	db.mux.Unlock() // unlock manual cause writeDB relocks

	return u, db.writeDB()
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

func (db *DB) UpdateUser(u entities.User) (entities.User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	oldUser, exits := db.store.Users[u.ID]
	if !exits {
		return entities.User{}, ErrDoesNotExist
	}

	encryptedUser, err := u.EncryptPassword()
	if err != nil {
		return entities.User{}, err
	}
	oldUser.Email = u.Email
	oldUser.Password = encryptedUser.Password
	db.store.Users[u.ID] = oldUser
	return oldUser, nil
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
