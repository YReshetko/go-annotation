package handlers

import (
	"fmt"
	"net/http"
)

// Container struct for handling container requests
// @Rest(path = "/api/v1/containers/{container_id}")
type Container struct{}

// @Rest(method = "GET")
func (c *Container) getContainer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `[
							{"title":"Container 1","id":"1", "text":"some text 1", "buttonText": "click me 1"},
							{"title":"Container 2","id":"2", "text":"some text 2", "buttonText": "click me 2"},
							{"title":"Container 3","id":"3", "text":"some text 3", "buttonText": "click me 3"},
							{"title":"Container 4","id":"4", "text":"some text 4", "buttonText": "click me 4"},
							{"title":"Container 5","id":"5", "text":"some text 5", "buttonText": "click me 5"}
					]`)
}

// @Rest(method = "PATCH")
func (c *Container) patchContainer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `[
							{"title":"Container 1","id":"1", "text":"some text 1", "buttonText": "click me 1"},
							{"title":"Container 2","id":"2", "text":"some text 2", "buttonText": "click me 2"},
							{"title":"Container 3","id":"3", "text":"some text 3", "buttonText": "click me 3"},
							{"title":"Container 4","id":"4", "text":"some text 4", "buttonText": "click me 4"},
							{"title":"Container 5","id":"5", "text":"some text 5", "buttonText": "click me 5"}
					]`)
}

// @Rest(method = "DELETE")
func (c *Container) deleteContainer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `[
							{"title":"Container 1","id":"1", "text":"some text 1", "buttonText": "click me 1"},
							{"title":"Container 2","id":"2", "text":"some text 2", "buttonText": "click me 2"},
							{"title":"Container 3","id":"3", "text":"some text 3", "buttonText": "click me 3"},
							{"title":"Container 4","id":"4", "text":"some text 4", "buttonText": "click me 4"},
							{"title":"Container 5","id":"5", "text":"some text 5", "buttonText": "click me 5"}
					]`)
}
