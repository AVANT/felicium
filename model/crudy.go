package model

import (
	"time"
)

type crudy interface {
	GetConnection() *Connection
	SetCreatedAt(time.Time)
	GetCreatedAt() time.Time
	SetUpdatedAt(time.Time)
	GetUpdatedAt() time.Time
	GetId() string
	GetRev() string
	SetId(string)
	SetRev(string)
}

func createOrUpdateCrudy(toAlter crudy) (string, string, error) {
	toAlter.SetUpdatedAt(time.Now())
	if toAlter.GetCreatedAt().IsZero() {
		toAlter.SetCreatedAt(time.Now())
	}

	//need to be careful here not to cause a conflict.
	//this is replacing the sync value hash calls so all crudy have the ability
	//ensures that accessors are used
	if toAlter.GetId() != "" {
		toAlter.SetId(toAlter.GetId())
		toAlter.SetRev(toAlter.GetRev())
	}
	id, rev, err := toAlter.GetConnection().Insert(toAlter)
	if err != nil {
		return "", "", err
	}
	//update the values after save.
	//this will ensure that all values are set with accessors.
	toAlter.SetId(id)
	toAlter.SetRev(rev)
	return id, rev, nil
}

func deleteCrudy(toDelete crudy) error {
	err := toDelete.GetConnection().Delete(toDelete.GetId(), toDelete.GetRev())
	return err
}
