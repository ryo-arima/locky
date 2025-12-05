package server

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/server/controller/internal"
	"github.com/ryo-arima/locky/pkg/server/controller/private"
	"github.com/ryo-arima/locky/pkg/server/controller/public"
	"github.com/ryo-arima/locky/pkg/server/repository"
	"github.com/ryo-arima/locky/pkg/server/share"
	"github.com/ryo-arima/locky/pkg/server/usecase"
)

func InitRouter(conf config.BaseConfig) *gin.Engine {
	// Set Gin to use our custom logger
	gin.DefaultWriter = share.NewGinLoggerWriter(conf.Logger)
	gin.DefaultErrorWriter = share.NewGinLoggerWriter(conf.Logger)

	redisClient, err := repository.NewRedisClient(conf.YamlConfig.Redis)
	if err != nil {
		panic(err)
	}

	// Casbin initialization: app-wide (locky) + group/resource permissions (resources)
	appEnforcer, err := casbin.NewEnforcer("etc/casbin/locky/model.conf", "etc/casbin/locky/policy.csv")
	if err != nil {
		panic(err)
	}
	if err := appEnforcer.LoadPolicy(); err != nil {
		log.Fatalf("failed to load app casbin policy: %v", err)
	}

	resourceEnforcer, err := casbin.NewEnforcer("etc/casbin/resources/model.conf", "etc/casbin/resources/policy.csv")
	if err != nil {
		panic(err)
	}
	if err := resourceEnforcer.LoadPolicy(); err != nil {
		log.Fatalf("failed to load resource casbin policy: %v", err)
	}

	userRepository := repository.NewUserRepository(conf)
	commonRepository := repository.NewCommonRepository(conf, redisClient)

	// Initialize usecase layer
	userUsecase := usecase.NewUserUsecase(userRepository)

	userControllerForPublic := public.NewUserController(userUsecase, commonRepository, conf)
	userControllerForInternal := internal.NewUserController(userUsecase, commonRepository)
	userControllerForPrivate := private.NewUserController(userUsecase, commonRepository)

	groupRepository := repository.NewGroupRepository(conf)
	groupControllerForInternal := internal.NewGroupController(groupRepository, commonRepository)
	groupControllerForPrivate := private.NewGroupController(groupRepository, commonRepository)

	memberRepository := repository.NewMemberRepository(conf)
	memberControllerForInternal := internal.NewMemberController(memberRepository, commonRepository)
	memberControllerForPrivate := private.NewMemberController(memberRepository, commonRepository)

	roleRepository := repository.NewRoleRepository(appEnforcer, resourceEnforcer)
	roleControllerForInternal := internal.NewRoleController(roleRepository, appEnforcer)
	roleControllerForPrivate := private.NewRoleController(roleRepository, appEnforcer)

	// CommonController for authentication endpoints
	commonControllerForPublic := public.NewCommonController(userRepository, commonRepository)

	router := gin.Default()

	// Health check endpoint (no authentication required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	loggerMW := share.LoggerWithConfig(conf)
	requestIDMW := share.RequestID()

	// OpenStack Keystone-style API versioning and structure
	// v1 API with proper versioning
	v1 := router.Group("/v1")
	v1.Use(requestIDMW) // Apply RequestID to all v1 routes

	// Authentication endpoints (Keystone-style, no middleware for login)
	auth := v1.Group("/share/common/auth")
	auth.Use(loggerMW)
	{
		auth.POST("/tokens", commonControllerForPublic.Login)                 // Issue token (login)
		auth.DELETE("/tokens", commonControllerForPublic.Logout)              // Revoke token (logout)
		auth.GET("/tokens/validate", commonControllerForPublic.ValidateToken) // Validate token
		auth.POST("/tokens/refresh", commonControllerForPublic.RefreshToken)  // Refresh token

		// GetUserInfo requires authentication middleware
		authWithMW := auth.Group("")
		authWithMW.Use(share.ForInternal(commonRepository, appEnforcer))
		authWithMW.GET("/tokens/user", commonControllerForPublic.GetUserInfo) // Get user info from token
	}

	// Public API - No authentication required (read-only discovery)
	publicAPI := v1.Group("/public")
	publicAPI.Use(loggerMW, share.ForPublic(conf))

	// Internal API - Authentication required (standard operations)
	internalAPI := v1.Group("/internal")
	internalAPI.Use(loggerMW, share.ForInternal(commonRepository, appEnforcer))

	// Private API - Administrative operations (Keystone admin endpoints style)
	privateAPI := v1.Group("/private")
	privateAPI.Use(loggerMW, share.ForPrivate(commonRepository, appEnforcer))

	// ============ USER ENDPOINTS ============
	// Public: User registration (POST uses singular)
	publicAPI.POST("/user", userControllerForPublic.CreateUser)
	// Internal: Standard user operations (GET plural, mutating singular)
	internalAPI.GET("/users", share.CasbinAuthorization(appEnforcer, "users", "read"), userControllerForInternal.GetUsers)
	internalAPI.GET("/users/count", share.CasbinAuthorization(appEnforcer, "users", "read"), userControllerForInternal.CountUsers)
	internalAPI.PUT("/user/:id", share.CasbinAuthorization(appEnforcer, "users", "write"), userControllerForInternal.UpdateUser)
	internalAPI.DELETE("/user/:id", share.CasbinAuthorization(appEnforcer, "users", "write"), userControllerForInternal.DeleteUser)
	// Private: Administrative user management
	privateAPI.GET("/users", share.CasbinAuthorization(appEnforcer, "users", "read"), userControllerForPrivate.GetUsers)
	privateAPI.GET("/users/count", share.CasbinAuthorization(appEnforcer, "users", "read"), userControllerForPrivate.CountUsers)
	privateAPI.POST("/user", share.CasbinAuthorization(appEnforcer, "users", "write"), userControllerForPrivate.CreateUser)
	privateAPI.PUT("/user/:id", share.CasbinAuthorization(appEnforcer, "users", "write"), userControllerForPrivate.UpdateUser)
	privateAPI.DELETE("/user/:id", share.CasbinAuthorization(appEnforcer, "users", "write"), userControllerForPrivate.DeleteUser)

	// ============ GROUP ENDPOINTS ============
	internalAPI.GET("/groups", share.CasbinAuthorization(appEnforcer, "groups", "read"), groupControllerForInternal.GetGroups)
	internalAPI.GET("/groups/count", share.CasbinAuthorization(appEnforcer, "groups", "read"), groupControllerForInternal.CountGroups)
	internalAPI.POST("/group", share.CasbinAuthorization(appEnforcer, "groups", "write"), groupControllerForInternal.CreateGroup)
	internalAPI.PUT("/group/:id", share.CasbinAuthorization(appEnforcer, "groups", "write"), groupControllerForInternal.UpdateGroup)
	internalAPI.DELETE("/group/:id", share.CasbinAuthorization(appEnforcer, "groups", "write"), groupControllerForInternal.DeleteGroup)
	privateAPI.GET("/groups", share.CasbinAuthorization(appEnforcer, "groups", "read"), groupControllerForPrivate.GetGroups)
	privateAPI.GET("/groups/count", share.CasbinAuthorization(appEnforcer, "groups", "read"), groupControllerForPrivate.CountGroups)
	privateAPI.POST("/group", share.CasbinAuthorization(appEnforcer, "groups", "write"), groupControllerForPrivate.CreateGroup)
	privateAPI.PUT("/group/:id", share.CasbinAuthorization(appEnforcer, "groups", "write"), groupControllerForPrivate.UpdateGroup)
	privateAPI.DELETE("/group/:id", share.CasbinAuthorization(appEnforcer, "groups", "write"), groupControllerForPrivate.DeleteGroup)

	// ============ MEMBER ENDPOINTS ============
	internalAPI.GET("/members", share.CasbinAuthorization(appEnforcer, "members", "read"), memberControllerForInternal.GetMembers)
	internalAPI.GET("/members/count", share.CasbinAuthorization(appEnforcer, "members", "read"), memberControllerForInternal.CountMembers)
	internalAPI.POST("/member", share.CasbinAuthorization(appEnforcer, "members", "write"), memberControllerForInternal.CreateMember)
	internalAPI.PUT("/member/:id", share.CasbinAuthorization(appEnforcer, "members", "write"), memberControllerForInternal.UpdateMember)
	internalAPI.DELETE("/member/:id", share.CasbinAuthorization(appEnforcer, "members", "write"), memberControllerForInternal.DeleteMember)
	privateAPI.GET("/members", share.CasbinAuthorization(appEnforcer, "members", "read"), memberControllerForPrivate.GetMembers)
	privateAPI.GET("/members/count", share.CasbinAuthorization(appEnforcer, "members", "read"), memberControllerForPrivate.CountMembers)
	privateAPI.POST("/member", share.CasbinAuthorization(appEnforcer, "members", "write"), memberControllerForPrivate.CreateMember)
	privateAPI.PUT("/member/:id", share.CasbinAuthorization(appEnforcer, "members", "write"), memberControllerForPrivate.UpdateMember)
	privateAPI.DELETE("/member/:id", share.CasbinAuthorization(appEnforcer, "members", "write"), memberControllerForPrivate.DeleteMember)

	// ===== ROLE (policy driven) =====
	internalAPI.GET("/roles", share.CasbinAuthorization(appEnforcer, "roles", "read"), roleControllerForInternal.ListRoles)
	privateAPI.GET("/roles", share.CasbinAuthorization(appEnforcer, "roles", "read"), roleControllerForPrivate.ListRoles)
	privateAPI.POST("/role", share.CasbinAuthorization(appEnforcer, "roles", "write"), roleControllerForPrivate.CreateRole)
	privateAPI.PUT("/role/:id", share.CasbinAuthorization(appEnforcer, "roles", "write"), roleControllerForPrivate.UpdateRole)
	privateAPI.DELETE("/role/:id", share.CasbinAuthorization(appEnforcer, "roles", "write"), roleControllerForPrivate.DeleteRole)

	return router
}
