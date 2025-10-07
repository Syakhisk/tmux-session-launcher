package rpc

// EmptyParams is used for methods that don't require parameters
type EmptyParams struct{}

type OpenInParams struct {
	Category  string `json:"category"`
	Path      string `json:"path"`
}

type ModeResponse struct {
	Mode string `json:"mode"`
}

type ContentResponse struct {
	Content string `json:"content"`
}

type EmptyResponse struct{}
