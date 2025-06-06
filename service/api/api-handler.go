package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.PUT("/users/:username", rt.wrap(rt.setMyUserName))
	rt.router.GET("/users/:username", rt.getUserProfile)
	rt.router.PUT("/users/:username/follow", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:username/follow", rt.wrap(rt.unfollowUser))
	rt.router.PUT("/users/:username/ban", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:username/ban", rt.wrap(rt.unbanUser))
	rt.router.GET("/users/:username/stream", rt.getMyStream)
	rt.router.POST("/images/:imageurl", rt.wrap(rt.uploadImage))
	rt.router.DELETE("/images/:imageurl", rt.deletePhoto)
	rt.router.PUT("/images/:imageurl/like", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/images/:imageurl/like", rt.unlikePhoto)
	rt.router.PUT("/images/:imageurl/comment", rt.wrap(rt.addComment))
	rt.router.DELETE("/images/:imageurl/comment", rt.wrap(rt.removeComment))
	rt.router.GET("/images/:imageurl", rt.getImageInfo)
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
