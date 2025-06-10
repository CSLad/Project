package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	//rt.router.GET("/users/:username", rt.wrap(rt.getUserProfile))
	//rt.router.DELETE("/images/:imageurl/like", rt.wrap(rt.unlikePhoto))

	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.PUT("/users/:username", rt.wrap(rt.setMyUserName))
	rt.router.PUT("/users/:username/follow", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:username/follow", rt.wrap(rt.unfollowUser))
	rt.router.PUT("/users/:username/ban", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:username/ban", rt.wrap(rt.unbanUser))
	rt.router.GET("/users/:username/stream", rt.wrap(rt.getMyStream))
	rt.router.POST("/images", rt.wrap(rt.uploadImage))
	rt.router.DELETE("/images/:imageid", rt.wrap(rt.deletePhoto))
	rt.router.PUT("/images/:imageid/like", rt.wrap(rt.likePhoto))
	rt.router.PUT("/images/:imageid/comment", rt.wrap(rt.addComment))
	rt.router.DELETE("/images/:imageid/comment", rt.wrap(rt.removeComment))
	rt.router.GET("/images/:imageid", rt.wrap(rt.getImageInfo))
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
