package client

type HTTPUtil struct{}

func NewHTTPUtil() *HTTPUtil {
    return &HTTPUtil{}
}

func (hu *HTTPUtil) Ping(url string) bool {
    return true // Mock ping
}