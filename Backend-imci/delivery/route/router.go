package route

import (
	"time"

	"github.com/Afomiat/Digital-IMCI/config"
	// "github.com/Afomiat/Digital-IMCI/delivery/middleware"
	// "github.com/Afomiat/Digital-IMCI/domain"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(env *config.Env, timeout time.Duration, db *pgxpool.Pool, r *gin.Engine) {
	// PublicRout := r.Group("")
	// NewSignUpRouter(env, timeout, db, PublicRout)
	// NewLogInRouter(env, timeout, db, PublicRout)
}
