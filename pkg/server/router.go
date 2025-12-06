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
	// Set Gin to use our custom logger (if available)
	if logger, ok := conf.Logger.(share.LoggerInterface); ok {
		gin.DefaultWriter = share.NewGinLoggerWriter(logger)
		gin.DefaultErrorWriter = share.NewGinLoggerWriter(logger)
	}

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

	userRepository := repository.NewUser(conf)
	commonRepository := repository.NewCommon(conf, redisClient)

	// Initialize usecase layer
	userUsecase := usecase.NewUser(userRepository)
	commonUsecase := usecase.NewCommon(commonRepository)

	publicUserController := controller.NewUserPublic(userUsecase, commonUsecase, conf)
	internalUserController := controller.NewUserInternal(userUsecase)
	privateUserController := controller.NewUserPrivate(userUsecase)

	groupRepository := repository.NewGroup(conf)
	groupUsecase := usecase.NewGroup(groupRepository)
	internalGroupController := controller.NewGroupInternal(groupUsecase, commonUsecase)
	privateGroupController := controller.NewGroupPrivate(groupUsecase, commonUsecase)

	memberRepository := repository.NewMember(conf)
	memberUsecase := usecase.NewMember(memberRepository)
	internalMemberController := controller.NewMemberInternal(memberUsecase)
	privateMemberController := controller.NewMemberPrivate(memberUsecase)

	roleRepository := repository.NewRole(appEnforcer, resourceEnforcer)
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
	privateAPI.GET("/users", authz(appEnforcer, "users", "read"), privateUserController.GetUsers)
	privateAPI.GET("/users/count", authz(appEnforcer, "users", "read"), privateUserController.CountUsers)
	privateAPI.POST("/user", authz(appEnforcer, "users", "write"), privateUserController.CreateUser)
	privateAPI.PUT("/user/:id", authz(appEnforcer, "users", "write"), privateUserController.UpdateUser)
	privateAPI.DELETE("/user/:id", authz(appEnforcer, "users", "write"), privateUserController.DeleteUser)

	// ============ GROUP ENDPOINTS ============
	internalAPI.GET("/groups", authz(appEnforcer, "groups", "read"), internalGroupController.GetGroups)
	internalAPI.GET("/groups/count", authz(appEnforcer, "groups", "read"), internalGroupController.CountGroups)
	internalAPI.POST("/group", authz(appEnforcer, "groups", "write"), internalGroupController.CreateGroup)
	internalAPI.PUT("/group/:id", authz(appEnforcer, "groups", "write"), internalGroupController.UpdateGroup)
	internalAPI.DELETE("/group/:id", authz(appEnforcer, "groups", "write"), internalGroupController.DeleteGroup)
	privateAPI.GET("/groups", authz(appEnforcer, "groups", "read"), privateGroupController.GetGroups)
	privateAPI.GET("/groups/count", authz(appEnforcer, "groups", "read"), privateGroupController.CountGroups)
	privateAPI.POST("/group", authz(appEnforcer, "groups", "write"), privateGroupController.CreateGroup)
	privateAPI.PUT("/group/:id", authz(appEnforcer, "groups", "write"), privateGroupController.UpdateGroup)
	privateAPI.DELETE("/group/:id", authz(appEnforcer, "groups", "write"), privateGroupController.DeleteGroup)

	// ============ MEMBER ENDPOINTS ============
	internalAPI.GET("/members", authz(appEnforcer, "members", "read"), internalMemberController.GetMembers)
	internalAPI.GET("/members/count", authz(appEnforcer, "members", "read"), internalMemberController.CountMembers)
	internalAPI.POST("/member", authz(appEnforcer, "members", "write"), internalMemberController.CreateMember)
	internalAPI.PUT("/member/:id", authz(appEnforcer, "members", "write"), internalMemberController.UpdateMember)
	internalAPI.DELETE("/member/:id", authz(appEnforcer, "members", "write"), internalMemberController.DeleteMember)
	privateAPI.GET("/members", authz(appEnforcer, "members", "read"), privateMemberController.GetMembers)
	privateAPI.GET("/members/count", authz(appEnforcer, "members", "read"), privateMemberController.CountMembers)
	privateAPI.POST("/member", authz(appEnforcer, "members", "write"), privateMemberController.CreateMember)
	privateAPI.PUT("/member/:id", authz(appEnforcer, "members", "write"), privateMemberController.UpdateMember)
	privateAPI.DELETE("/member/:id", authz(appEnforcer, "members", "write"), privateMemberController.DeleteMember)

	// ===== ROLE (policy driven) =====
	internalAPI.GET("/roles", authz(appEnforcer, "roles", "read"), internalRoleController.ListRoles)
	privateAPI.GET("/roles", authz(appEnforcer, "roles", "read"), privateRoleController.ListRoles)
	privateAPI.POST("/role", authz(appEnforcer, "roles", "write"), privateRoleController.CreateRole)
	privateAPI.PUT("/role/:id", authz(appEnforcer, "roles", "write"), privateRoleController.UpdateRole)
	privateAPI.DELETE("/role/:id", authz(appEnforcer, "roles", "write"), privateRoleController.DeleteRole)

	// ============ RESOURCE ENDPOINTS ============
	internalAPI.GET("/resources", authz(appEnforcer, "resources", "read"), internalResourceController.GetResources)
	internalAPI.GET("/resources/count", authz(appEnforcer, "resources", "read"), internalResourceController.CountResources)
	internalAPI.POST("/resource", authz(appEnforcer, "resources", "write"), internalResourceController.CreateResource)
	internalAPI.PUT("/resource/:id", authz(appEnforcer, "resources", "write"), internalResourceController.UpdateResource)
	internalAPI.DELETE("/resource/:id", authz(appEnforcer, "resources", "write"), internalResourceController.DeleteResource)
	privateAPI.GET("/resources", authz(appEnforcer, "resources", "read"), privateResourceController.GetResources)
	privateAPI.GET("/resources/count", authz(appEnforcer, "resources", "read"), privateResourceController.CountResources)
	privateAPI.POST("/resource", authz(appEnforcer, "resources", "write"), privateResourceController.CreateResource)
	privateAPI.PUT("/resource/:id", authz(appEnforcer, "resources", "write"), privateResourceController.UpdateResource)
	privateAPI.DELETE("/resource/:id", authz(appEnforcer, "resources", "write"), privateResourceController.DeleteResource)

	return router
}
