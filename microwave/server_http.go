package microwave

import "net/http"

func NewHTTPServer() *http.Server {
	// TODO
	return &http.Server{}
}

func (s HTTPWrapper) Run(mw *Microwave) {
	mw.logger.Info("http server running").WithField("port", s.Port).Send()
	// TODO
}
