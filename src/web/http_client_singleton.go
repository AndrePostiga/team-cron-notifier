package web

import (
	"net/http"
	"sync"
	"time"

	"github.com/andrepostiga/team-cron-notifier/src/config"
)

//var lock = &sync.Mutex{}

var singleton *http.Client

var once sync.Once

func NewClientSingleton(options *config.HttpClientConfig) *http.Client {
	once.Do(func() {
		singleton = &http.Client{
			Timeout: time.Duration(options.TimeoutInSeconds) * time.Second,
		}
	})

	//if singleton == nil {
	//	lock.Lock()
	//	defer lock.Unlock()
	//	singleton = &http.Client{
	//		Timeout: time.Duration(options.TimeoutInSeconds) * time.Second,
	//	}
	//}

	return singleton
}
