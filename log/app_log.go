package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"gitlab.com/apus-backend/base-service/config"
	"os"
)

type AppLog struct {
}

func NewAppLog(c config.Config) (*AppLog, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	level, err := zerolog.ParseLevel(c.Log.Level)
	if err != nil {
		return nil, err
	}
	zerolog.SetGlobalLevel(level)
	/*zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}*/
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	return &AppLog{}, nil
}
