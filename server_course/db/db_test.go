package db

import (
	"server_course/entities"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDB_writeDB(t *testing.T) {
	db, err := NewDB(".")
	assert.NoError(t, err)

	chirp := entities.Chirp{
		Body: "Hello World",
	}
	_, err = db.StoreChirp(chirp)
	assert.NoError(t, err)

	chirps, err := db.GetChirps()
	assert.NoError(t, err)
	assert.Equal(t, chirps[0], chirp)

	db2, err := NewDB(".")
	assert.NoError(t, err)
	chirps2, err := db2.GetChirps()
	assert.NoError(t, err)
	assert.Equal(t, chirps2[0], chirp)
	assert.Equal(t, chirps, chirps2)
}
