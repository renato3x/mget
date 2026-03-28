package downloader

type Provider interface {
	CanHandle(url string) bool
	Handle(url string) error
}
