package cache

import (
	"context"
	"fmt"
	"github.com/coderconquerer/social-todo/module/authentication/entity"
	"log"
	"sync"
)

// RealStore defines the interface for fetching users from the main data store.
type RealStore interface {
	FindAccount(ctx context.Context, conditions map[string]interface{}) (*entity.Account, error)
}

// userCaching wraps a real store with caching.
type userCaching struct {
	store     Cache
	realStore RealStore
	once      *sync.Once
}

// NewUserCaching creates a new user caching layer.
func NewUserCaching(store Cache, realStore RealStore) *userCaching {
	return &userCaching{
		store:     store,
		realStore: realStore,
		once:      new(sync.Once),
	}
}

// FindAccount tries to find a user in the cache first, then falls back to the real store.
func (uc *userCaching) FindAccount(ctx context.Context, conditions map[string]interface{}) (*entity.Account, error) {
	var account entity.Account

	// build cache key based on user id
	userId := conditions["Id"].(int)
	key := fmt.Sprintf("user-%d", userId)

	// try cache first
	err := uc.store.Get(ctx, key, &account)
	if err == nil && account.Id > 0 {
		return &account, nil
	}

	var userErr error
	uc.once.Do(func() {
		realUser, findErr := uc.realStore.FindAccount(ctx, conditions)
		if findErr != nil {
			userErr = findErr
			log.Println("FindUser error:", findErr)
			return
		}

		// update cache
		_ = uc.store.Set(ctx, key, realUser, 0)
		account = *realUser
	})

	if userErr != nil {
		return nil, userErr
	}

	return &account, nil
}
