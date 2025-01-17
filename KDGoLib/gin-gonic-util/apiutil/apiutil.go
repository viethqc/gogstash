package apiutil

import (
	"sort"

	"github.com/gin-gonic/gin"
)

// BindRouter bind all registed apis to router
func BindRouter(router gin.IRouter) {
	for _, apiconf := range registedAPIMethodPathMap {
		method := ParseMethod(apiconf.Method)
		if method == MethodAny {
			router.Any(apiconf.Path, apiconf.Handlers...)
		} else {
			router.Handle(apiconf.Method, apiconf.Path, apiconf.Handlers...)
		}
	}
}

// RegistedAPIMap return all registed API map, key generated by getMethodPath()
func RegistedAPIMap() map[string]apiConfig {
	return registedAPIMethodPathMap
}

// RegistedSortedAPIMap return all registed API sorted slice
func RegistedSortedAPIMap() []apiConfig {
	results := []apiConfig{}

	keys := []string{}
	for key := range registedAPICallerMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		apiconf := registedAPICallerMap[key]
		results = append(results, apiconf)
	}

	return results
}
