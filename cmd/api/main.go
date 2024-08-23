package main

import (
	"context"

	"github.com/andrepostiga/team-cron-notifier/src/app"
)

func main() {
	ctx := context.Background()
	app.BuildApi(ctx)
}
