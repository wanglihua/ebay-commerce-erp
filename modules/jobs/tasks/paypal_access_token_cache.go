package tasks

import (
	"sync"
	"time"
)

var paypal_access_token_cache paypal_access_token_cache_t

type paypal_access_token_cache_t struct {
	lock                            sync.RWMutex
	paypal_access_token_cache_items map[string]paypal_access_token_cache_item
}

type paypal_access_token_cache_item struct {
	paypal_access_token         string
	paypal_access_token_expired time.Time
}

func (t paypal_access_token_cache_t) init() {
	t.paypal_access_token_cache_items = make(map[string]paypal_access_token_cache_item)
}

func (t paypal_access_token_cache_t) set(paypal_account string, item paypal_access_token_cache_item) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.paypal_access_token_cache_items[paypal_account] = item
}

func (t paypal_access_token_cache_t) get(paypal_account string) paypal_access_token_cache_item {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if item, has := t.paypal_access_token_cache_items[paypal_account]; has == true {
		return item
	} else {
		return paypal_access_token_cache_item{} //返回一个默认的空值
	}
}
