package property

import (
	"context"
	"kn-assignment/internal/log"

	"github.com/kelseyhightower/envconfig"
)

func InitProperty(ctx context.Context) {
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf(ctx, "read env error : %s", err.Error())
	}

}
