package main

import (
	"errors"
	"math/rand"
	"strings"
	"sync"
)

type DAO interface {
	GetShortener(url string) string
	GetFullURL(short string) (string, error)
}

type _dao struct {
	keysToUrls map[string]string
	urlsToKeys map[string]string
	mtx        sync.Mutex
}

var _ DAO = (*_dao)(nil)

var KeyNotFound = errors.New("key not found")

func NewDAO() DAO {
	dao := _dao{
		keysToUrls: make(map[string]string),
		urlsToKeys: make(map[string]string),
	}
	return &dao
}

func (dao *_dao) GetShortener(url string) string {
	dao.mtx.Lock()
	defer dao.mtx.Unlock()

	key, urlsIsExist := dao.urlsToKeys[url]

	if !urlsIsExist {
		key = genRandomKey()
		for _, c := dao.keysToUrls[key]; c; _, c = dao.keysToUrls[key] {
			key = genRandomKey()
		}

		dao.keysToUrls[key] = url
		dao.urlsToKeys[url] = key
	}

	return key
}

func (dao *_dao) GetFullURL(short string) (string, error) {
	dao.mtx.Lock()
	defer dao.mtx.Unlock()

	url, contain := dao.keysToUrls[short]

	if !contain {
		return "", KeyNotFound
	}

	return url, nil
}

var (
	randomMax = 10 + 26*2
	keyLen    = 10
)

func genRandomKey() string {
	sb := strings.Builder{}

	for i := 0; i < keyLen; i++ {
		n := rand.Intn(randomMax)
		if n < 10 {
			sb.WriteByte('0' + byte(n))
		} else if n < 10+26 {
			sb.WriteByte('a' + byte(n-10))
		} else {
			sb.WriteByte('A' + byte(n-10-26))
		}
	}

	return sb.String()
}
