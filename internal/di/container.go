package di

import (
	"avito-backend-bootcamp/internal/config"
	"avito-backend-bootcamp/internal/http/server"
	"avito-backend-bootcamp/internal/infra/cache"
	sender "avito-backend-bootcamp/internal/infra/email"
	"avito-backend-bootcamp/internal/infra/jwt"
	"avito-backend-bootcamp/internal/infra/repository/postgres"
	"avito-backend-bootcamp/internal/service/auth"
	emailsender "avito-backend-bootcamp/internal/service/email-sender"
	"avito-backend-bootcamp/internal/service/flat"
	"avito-backend-bootcamp/internal/service/house"
	sub "avito-backend-bootcamp/internal/service/subscription"
	dbUtil "avito-backend-bootcamp/pkg/utils/db"
	"context"
	"log/slog"

	_ "github.com/lib/pq"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

type Container struct {
	cfg *config.Config
	log *slog.Logger

	flatCache   *cache.TTLCache[int64, string]
	validator   *validator.Validate
	jwt         *jwt.Manager
	emailClient *sender.Sender
	repository  *postgres.Repository
	db          *sqlx.DB
	trManager   *manager.Manager

	flatService  *flat.Service
	houseService *house.Service
	subService   *sub.Service
	authService  *auth.Service
	emailService *emailsender.Service

	serverHTTP *server.Server
}

func New(cfg *config.Config, log *slog.Logger) *Container {
	return &Container{
		cfg: cfg,
		log: log,
	}
}

func (c *Container) GetTrManager() *manager.Manager {
	return get(&c.trManager, func() *manager.Manager {
		return manager.Must(trmsqlx.NewDefaultFactory(c.GetDB()))
	})
}

func (c *Container) GetDB() *sqlx.DB {
	return get(&c.db, func() *sqlx.DB {
		db, err := sqlx.Open(
			"postgres",
			dbUtil.BuildDSN(
				c.cfg.DB.Username,
				c.cfg.DB.Password,
				c.cfg.DB.Name,
				c.cfg.DB.Host,
				c.cfg.DB.Port,
			),
		)
		if err != nil {
			panic(err)
		}
		return db
	})
}

func (c *Container) GetFlatCache() *cache.TTLCache[int64, string] {
	return get(&c.flatCache, func() *cache.TTLCache[int64, string] {
		return cache.NewTTL[int64, string](c.cfg.Cache.TTL)
	})
}

func (c *Container) GetValidator() *validator.Validate {
	return get(&c.validator, func() *validator.Validate {
		return validator.New()
	})
}

func (c *Container) GetJwtManager() *jwt.Manager {
	return get(&c.jwt, func() *jwt.Manager {
		return jwt.New(c.cfg.JWT.SecretKey, c.cfg.JWT.TokenTTL)
	})
}

func (c *Container) GetEmailClient() *sender.Sender {
	return get(&c.emailClient, func() *sender.Sender {
		return sender.New()
	})
}

func (c *Container) GetRepository() *postgres.Repository {
	return get(&c.repository, func() *postgres.Repository {
		db, err := postgres.New(context.Background(), &c.cfg.DB)
		if err != nil {
			panic(err)
		}
		return db
	})
}

func (c *Container) GetAuthService() *auth.Service {
	return get(&c.authService, func() *auth.Service {
		return auth.New(
			c.log,
			c.GetJwtManager(),
			c.GetRepository(),
		)
	})
}

func (c *Container) GetFlatService() *flat.Service {
	return get(&c.flatService, func() *flat.Service {
		return flat.New(
			c.log,
			c.GetRepository(),
			c.GetRepository(),
			c.GetFlatCache(),
			c.GetTrManager(),
		)
	})
}

func (c *Container) GetHouseService() *house.Service {
	return get(&c.houseService, func() *house.Service {
		return house.New(
			c.log,
			c.GetRepository(),
		)
	})
}

func (c *Container) GetSenderService() *emailsender.Service {
	return get(&c.emailService, func() *emailsender.Service {
		return emailsender.New(
			c.log,
			c.GetEmailClient(),
			c.GetRepository(),
			c.GetRepository(),
			c.GetRepository(),
		)
	})
}

func (c *Container) GetSubsciptionService() *sub.Service {
	return get(&c.subService, func() *sub.Service {
		return sub.New(
			c.log,
			c.GetRepository(),
		)
	})
}

func (c *Container) GetHTTPServer() *server.Server {
	return get(&c.serverHTTP, func() *server.Server {
		srv, err := server.New(
			c.cfg,
			c.log,
			c.GetValidator(),
			c.GetAuthService(),
			c.GetFlatService(),
			c.GetHouseService(),
			c.GetSubsciptionService(),
			c.GetJwtManager(),
		)
		if err != nil {
			panic(err)
		}
		return srv
	})
}

func get[T comparable](obj *T, builder func() T) T {
	if *obj != *new(T) {
		return *obj
	}

	*obj = builder()
	return *obj
}
