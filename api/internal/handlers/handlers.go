package handlers

import (
	"github.com/edgejay/pify-player/api/internal/constants"
	"github.com/edgejay/pify-player/api/internal/database"
	"github.com/edgejay/pify-player/api/internal/middlewares"
	"github.com/edgejay/pify-player/api/internal/services"
)

var spotifyService *services.SpotifyService = services.NewSpotifyService(services.GetSpotifyCredentials(), nil)

var userService *services.UserService = services.NewUserService(database.GetSQLiteDB())

var middlewareFactory *middlewares.MiddlewareFactory = middlewares.NewMiddlewareFactory(
	constants.COOKIE_SESSION_ID,
	userService,
	spotifyService,
)
