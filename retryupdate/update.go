// +build !solution

package retryupdate

import (
	"errors"
	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	var authErr *kvapi.AuthError
	var conflictError *kvapi.ConflictError
	getParams := kvapi.GetRequest{
		Key: key,
	}

all:
	for {
		var oldValue *string
		var oldVersion uuid.UUID

		for {
			getValue, err := c.Get(&getParams)
			if err == nil {
				oldValue, oldVersion = &getValue.Value, getValue.Version
				break
			} else if errors.Is(err, kvapi.ErrKeyNotFound) {
				break
			} else if errors.As(err, &authErr) {
				return err
			}
		}

		newValue, err := updateFn(oldValue)
		if err != nil {
			return err
		}

		set := kvapi.SetRequest{
			Key:        key,
			Value:      newValue,
			OldVersion: oldVersion,
			NewVersion: uuid.Must(uuid.NewV4()),
		}

		for {
			_, err = c.Set(&set)
			if errors.As(err, &conflictError) {
				break
			} else if err == nil {
				break all
			}
		}
	}

	return nil
}
