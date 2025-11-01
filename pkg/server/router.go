package server

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/server/controller"
	"github.com/ryo-arima/locky/pkg/server/middleware"
	"github.com/ryo-arima/locky/pkg/server/repository"
)

func InitRouter(conf config.BaseConfig) *gin.Engine {

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
	userControllerForPublic := controller.NewUserControllerForPublic(userRepository, commonRepository)
	userControllerForInternal := controller.NewUserControllerForInternal(userRepository, commonRepository)
	userControllerForPrivate := controller.NewUserControllerForPrivate(userRepository, commonRepository)

	groupRepository := repository.NewGroupRepository(conf)
	groupControllerForInternal := controller.NewGroupControllerForInternal(groupRepository, commonRepository)
	groupControllerForPrivate := controller.NewGroupControllerForPrivate(groupRepository, commonRepository)

	memberRepository := repository.NewMemberRepository(conf)
	memberControllerForInternal := controller.NewMemberControllerForInternal(memberRepository, commonRepository)
	memberControllerForPrivate := controller.NewMemberControllerForPrivate(memberRepository, commonRepository)

	roleRepository := repository.NewRoleRepository(appEnforcer, resourceEnforcer)
	roleControllerForInternal := controller.NewRoleControllerForInternal(roleRepository, appEnforcer)
	roleControllerForPrivate := controller.NewRoleControllerForPrivate(roleRepository, appEnforcer)

	// CommonController for authentication endpoints
	commonControllerForPublic := controller.NewCommonControllerForPublic(userRepository, commonRepository)

	router := gin.Default()

	loggerMW := middleware.LoggerWithConfig(conf)

	// OpenStack Keystone-style API versioning and structure
	// v1 API with proper versioning
	v1 := router.Group("/v1")

	// Authentication endpoints (Keystone-style, no middleware for login)
	auth := v1.Group("/share/common/auth")
	auth.Use(loggerMW)
	{
		auth.POST("/tokens", commonControllerForPublic.Login)                 // Issue token (login)
		auth.DELETE("/tokens", commonControllerForPublic.Logout)              // Revoke token (logout)
		auth.GET("/tokens/validate", commonControllerForPublic.ValidateToken) // Validate token
		auth.POST("/tokens/refresh", commonControllerForPublic.RefreshToken)  // Refresh token
		auth.GET("/tokens/user", commonControllerForPublic.GetUserInfo)       // Get user info from token
	}

	// Public API - No authentication required (read-only discovery)
	publicAPI := v1.Group("/public")
	publicAPI.Use(loggerMW, middleware.ForPublic(conf))

	// Internal API - Authentication required (standard operations)
	internalAPI := v1.Group("/internal")
	internalAPI.Use(loggerMW, middleware.ForInternal(commonRepository, appEnforcer))

	// Private API - Administrative operations (Keystone admin endpoints style)
	privateAPI := v1.Group("/private")
	privateAPI.Use(loggerMW, middleware.ForPrivate(commonRepository, appEnforcer))

	// ============ USER ENDPOINTS ============
	// Public: User registration (POST uses singular)
	publicAPI.POST("/user", userControllerForPublic.CreateUser)
	// Internal: Standard user operations (GET plural, mutating singular)
	internalAPI.GET("/users", middleware.CasbinAuthorization(appEnforcer, "users", "read"), userControllerForInternal.GetUsers)
	internalAPI.GET("/users/count", middleware.CasbinAuthorization(appEnforcer, "users", "read"), userControllerForInternal.CountUsers)
	internalAPI.PUT("/user/:id", middleware.CasbinAuthorization(appEnforcer, "users", "write"), userControllerForInternal.UpdateUser)
	internalAPI.DELETE("/user/:id", middleware.CasbinAuthorization(appEnforcer, "users", "write"), userControllerForInternal.DeleteUser)
	// Private: Administrative user management
	privateAPI.GET("/users", middleware.CasbinAuthorization(appEnforcer, "users", "read"), userControllerForPrivate.GetUsers)
	privateAPI.GET("/users/count", middleware.CasbinAuthorization(appEnforcer, "users", "read"), userControllerForPrivate.CountUsers)
	privateAPI.POST("/user", middleware.CasbinAuthorization(appEnforcer, "users", "write"), userControllerForPrivate.CreateUser)
	privateAPI.PUT("/user/:id", middleware.CasbinAuthorization(appEnforcer, "users", "write"), userControllerForPrivate.UpdateUser)
	privateAPI.DELETE("/user/:id", middleware.CasbinAuthorization(appEnforcer, "users", "write"), userControllerForPrivate.DeleteUser)

	// ============ GROUP ENDPOINTS ============
	internalAPI.GET("/groups", middleware.CasbinAuthorization(appEnforcer, "groups", "read"), groupControllerForInternal.GetGroups)
	internalAPI.GET("/groups/count", middleware.CasbinAuthorization(appEnforcer, "groups", "read"), groupControllerForInternal.CountGroups)
	internalAPI.POST("/group", middleware.CasbinAuthorization(appEnforcer, "groups", "write"), groupControllerForInternal.CreateGroup)
	internalAPI.PUT("/group/:id", middleware.CasbinAuthorization(appEnforcer, "groups", "write"), groupControllerForInternal.UpdateGroup)
	internalAPI.DELETE("/group/:id", middleware.CasbinAuthorization(appEnforcer, "groups", "write"), groupControllerForInternal.DeleteGroup)
	privateAPI.GET("/groups", middleware.CasbinAuthorization(appEnforcer, "groups", "read"), groupControllerForPrivate.GetGroups)
	privateAPI.GET("/groups/count", middleware.CasbinAuthorization(appEnforcer, "groups", "read"), groupControllerForPrivate.CountGroups)
	privateAPI.POST("/group", middleware.CasbinAuthorization(appEnforcer, "groups", "write"), groupControllerForPrivate.CreateGroup)
	privateAPI.PUT("/group/:id", middleware.CasbinAuthorization(appEnforcer, "groups", "write"), groupControllerForPrivate.UpdateGroup)
	privateAPI.DELETE("/group/:id", middleware.CasbinAuthorization(appEnforcer, "groups", "write"), groupControllerForPrivate.DeleteGroup)

	// ============ MEMBER ENDPOINTS ============
	internalAPI.GET("/members", middleware.CasbinAuthorization(appEnforcer, "members", "read"), memberControllerForInternal.GetMembers)
	internalAPI.GET("/members/count", middleware.CasbinAuthorization(appEnforcer, "members", "read"), memberControllerForInternal.CountMembers)
	internalAPI.POST("/member", middleware.CasbinAuthorization(appEnforcer, "members", "write"), memberControllerForInternal.CreateMember)
	internalAPI.PUT("/member/:id", middleware.CasbinAuthorization(appEnforcer, "members", "write"), memberControllerForInternal.UpdateMember)
	internalAPI.DELETE("/member/:id", middleware.CasbinAuthorization(appEnforcer, "members", "write"), memberControllerForInternal.DeleteMember)
	privateAPI.GET("/members", middleware.CasbinAuthorization(appEnforcer, "members", "read"), memberControllerForPrivate.GetMembers)
	privateAPI.GET("/members/count", middleware.CasbinAuthorization(appEnforcer, "members", "read"), memberControllerForPrivate.CountMembers)
	privateAPI.POST("/member", middleware.CasbinAuthorization(appEnforcer, "members", "write"), memberControllerForPrivate.CreateMember)
	privateAPI.PUT("/member/:id", middleware.CasbinAuthorization(appEnforcer, "members", "write"), memberControllerForPrivate.UpdateMember)
	privateAPI.DELETE("/member/:id", middleware.CasbinAuthorization(appEnforcer, "members", "write"), memberControllerForPrivate.DeleteMember)

	// ===== ROLE (policy driven) =====
	internalAPI.GET("/roles", middleware.CasbinAuthorization(appEnforcer, "roles", "read"), roleControllerForInternal.ListRoles)
	privateAPI.GET("/roles", middleware.CasbinAuthorization(appEnforcer, "roles", "read"), roleControllerForPrivate.ListRoles)
	privateAPI.POST("/role", middleware.CasbinAuthorization(appEnforcer, "roles", "write"), roleControllerForPrivate.CreateRole)
	privateAPI.PUT("/role/:id", middleware.CasbinAuthorization(appEnforcer, "roles", "write"), roleControllerForPrivate.UpdateRole)
	privateAPI.DELETE("/role/:id", middleware.CasbinAuthorization(appEnforcer, "roles", "write"), roleControllerForPrivate.DeleteRole)

	return router
}
