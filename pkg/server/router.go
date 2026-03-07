package server

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/server/controller"
	"github.com/ryo-arima/locky/pkg/server/repository"
	"github.com/ryo-arima/locky/pkg/server/share"
	"github.com/ryo-arima/locky/pkg/server/usecase"
)

func InitRouter(conf config.BaseConfig) *gin.Engine {
	// Initialize global server logger
	if logger, ok := conf.Logger.(*share.ServerLogger); ok {
		share.SetServerLogger(logger)
	}

	// Set Gin to use our custom logger (if available)
	if logger, ok := conf.Logger.(share.LoggerInterface); ok {
		gin.DefaultWriter = share.NewGinLoggerWriter(logger)
		gin.DefaultErrorWriter = share.NewGinLoggerWriter(logger)
	}

	redisClient, err := repository.NewRedisClient(conf.YamlConfig.Redis)
	if err != nil {
		panic(err)
	}

	// Casbin initialization: app-wide (locky) for group/resource permissions
	appEnforcer, err := casbin.NewEnforcer("etc/casbin/locky/model.conf", "etc/casbin/locky/policy.csv")
	if err != nil {
		panic(err)
	}
	if err := appEnforcer.LoadPolicy(); err != nil {
		log.Fatalf("failed to load app casbin policy: %v", err)
	}

	userRepository := repository.NewUser(conf)
	commonRepository := repository.NewCommon(conf, redisClient)

	// Initialize usecase layer
	userUsecase := usecase.NewUser(userRepository)
	commonUsecase := usecase.NewCommon(commonRepository)

	publicUserController := controller.NewUserPublic(userUsecase, commonUsecase, conf)
	internalUserController := controller.NewUserInternal(userUsecase)
	privateUserController := controller.NewUserPrivate(userUsecase)

	groupRepository := repository.NewGroup(conf)
	memberRepository := repository.NewMember(conf)
	groupUsecase := usecase.NewGroup(groupRepository, memberRepository, conf.DBConnection)
	internalGroupController := controller.NewGroupInternal(groupUsecase, commonUsecase)
	privateGroupController := controller.NewGroupPrivate(groupUsecase, commonUsecase)

	memberUsecase := usecase.NewMember(memberRepository)
	internalMemberController := controller.NewMemberInternal(memberUsecase)
	privateMemberController := controller.NewMemberPrivate(memberUsecase)

	roleRepository := repository.NewRole(appEnforcer)
	roleUsecase := usecase.NewRole(roleRepository)
	internalRoleController := controller.NewRoleInternal(roleUsecase, appEnforcer)
	privateRoleController := controller.NewRolePrivate(roleUsecase, appEnforcer)

	resourceRepository := repository.NewResource(conf)
	resourceUsecase := usecase.NewResource(resourceRepository, memberRepository)
	internalResourceController := controller.NewResourceInternal(resourceUsecase)
	privateResourceController := controller.NewResourcePrivate(resourceUsecase)

	// CommonController for authentication endpoints
	commonShareController := controller.NewCommonShare(userUsecase, commonUsecase)

	router := gin.Default()

	// Health check endpoint (no authentication required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	var loggerMW gin.HandlerFunc
	if logger, ok := conf.Logger.(share.LoggerInterface); ok {
		loggerMW = share.LoggerWithConfig(logger)
	} else {
		loggerMW = func(c *gin.Context) { c.Next() } // No-op middleware
	}
	requestIDMW := share.RequestID()
	authz := share.CasbinAuthorization

	// OpenStack Keystone-style API versioning and structure
	// v1 API with proper versioning
	v1 := router.Group("/v1")
	v1.Use(requestIDMW) // Apply RequestID to all v1 routes

	// Share API - Authentication endpoints
	shareAPI := v1.Group("/share")
	shareAPI.Use(loggerMW)

	// Public API - No authentication required (read-only discovery)
	publicAPI := v1.Group("/public")
	publicAPI.Use(loggerMW, share.ForPublic(conf))

	// Internal API - Authentication required (standard operations)
	internalAPI := v1.Group("/internal")
	internalAPI.Use(loggerMW, share.ForInternal(commonRepository, appEnforcer))

	// Private API - Administrative operations (Keystone admin endpoints style)
	privateAPI := v1.Group("/private")
	privateAPI.Use(loggerMW, share.ForPrivate(commonRepository, appEnforcer))

	// ============ COMMON ENDPOINTS ============
	publicAPI.POST("/token", commonShareController.Login)
	shareAPI.POST("/token/refresh", share.ForShare(commonRepository, appEnforcer), commonShareController.RefreshToken)
	shareAPI.DELETE("/token", share.ForShare(commonRepository, appEnforcer), commonShareController.Logout)
	shareAPI.GET("/token/validate", share.ForShare(commonRepository, appEnforcer), commonShareController.ValidateToken)
	shareAPI.GET("/token/user", share.ForShare(commonRepository, appEnforcer), commonShareController.GetUserInfo)

	// ============ USER ENDPOINTS ============
	// Public: User registration (POST uses singular)
	publicAPI.POST("/user", publicUserController.CreateUser)
	// Internal: Standard user operations (GET plural, mutating singular)
	internalAPI.GET("/users", authz(appEnforcer, "users", "read"), internalUserController.GetUsers)
	internalAPI.GET("/users/count", authz(appEnforcer, "users", "read"), internalUserController.CountUsers)
	internalAPI.PUT("/user/:id", authz(appEnforcer, "users", "write"), internalUserController.UpdateUser)
	internalAPI.DELETE("/user/:id", authz(appEnforcer, "users", "write"), internalUserController.DeleteUser)
	// Private: Administrative user management
	privateAPI.GET("/users", privateUserController.GetUsers)
	privateAPI.GET("/users/count", privateUserController.CountUsers)
	privateAPI.POST("/user", privateUserController.CreateUser)
	privateAPI.PUT("/user/:id", privateUserController.UpdateUser)
	privateAPI.DELETE("/user/:id", privateUserController.DeleteUser)

	// ============ GROUP ENDPOINTS ============
	internalAPI.GET("/groups", authz(appEnforcer, "groups", "read"), internalGroupController.GetGroups)
	internalAPI.GET("/groups/count", authz(appEnforcer, "groups", "read"), internalGroupController.CountGroups)
	internalAPI.POST("/group", authz(appEnforcer, "groups", "write"), internalGroupController.CreateGroup)
	internalAPI.PUT("/group/:id", authz(appEnforcer, "groups", "write"), internalGroupController.UpdateGroup)
	internalAPI.DELETE("/group/:id", authz(appEnforcer, "groups", "write"), internalGroupController.DeleteGroup)
	privateAPI.GET("/groups", privateGroupController.GetGroups)
	privateAPI.GET("/groups/count", privateGroupController.CountGroups)
	privateAPI.POST("/group", privateGroupController.CreateGroup)
	privateAPI.PUT("/group/:id", privateGroupController.UpdateGroup)
	privateAPI.DELETE("/group/:id", privateGroupController.DeleteGroup)

	// ============ MEMBER ENDPOINTS ============
	internalAPI.GET("/members", authz(appEnforcer, "members", "read"), internalMemberController.GetMembers)
	internalAPI.GET("/members/count", authz(appEnforcer, "members", "read"), internalMemberController.CountMembers)
	internalAPI.POST("/member", authz(appEnforcer, "members", "write"), internalMemberController.CreateMember)
	internalAPI.PUT("/member/:id", authz(appEnforcer, "members", "write"), internalMemberController.UpdateMember)
	internalAPI.DELETE("/member/:id", authz(appEnforcer, "members", "write"), internalMemberController.DeleteMember)
	privateAPI.GET("/members", privateMemberController.GetMembers)
	privateAPI.GET("/members/count", privateMemberController.CountMembers)
	privateAPI.POST("/member", privateMemberController.CreateMember)
	privateAPI.PUT("/member/:id", privateMemberController.UpdateMember)
	privateAPI.DELETE("/member/:id", privateMemberController.DeleteMember)

	// ===== ROLE (policy driven) =====
	internalAPI.GET("/roles", authz(appEnforcer, "roles", "read"), internalRoleController.ListRoles)
	privateAPI.GET("/roles", privateRoleController.ListRoles)
	privateAPI.POST("/role", privateRoleController.CreateRole)
	privateAPI.PUT("/role/:id", privateRoleController.UpdateRole)
	privateAPI.DELETE("/role/:id", privateRoleController.DeleteRole)

	// ============ RESOURCE ENDPOINTS ============
	internalAPI.GET("/resources", authz(appEnforcer, "resources", "read"), internalResourceController.GetResources)
	internalAPI.GET("/resources/count", authz(appEnforcer, "resources", "read"), internalResourceController.CountResources)
	internalAPI.POST("/resource", authz(appEnforcer, "resources", "write"), internalResourceController.CreateResource)
	internalAPI.PUT("/resource/:id", authz(appEnforcer, "resources", "write"), internalResourceController.UpdateResource)
	internalAPI.DELETE("/resource/:id", authz(appEnforcer, "resources", "write"), internalResourceController.DeleteResource)
	privateAPI.GET("/resources", privateResourceController.GetResources)
	privateAPI.GET("/resources/count", privateResourceController.CountResources)
	privateAPI.POST("/resource", privateResourceController.CreateResource)
	privateAPI.PUT("/resource/:id", privateResourceController.UpdateResource)
	privateAPI.DELETE("/resource/:id", privateResourceController.DeleteResource)

	return router
}
