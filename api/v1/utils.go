package v1

import (
	"errors"
	"fmt"
	log "github.com/DataDrake/waterlog"
	"github.com/valyala/fasthttp"
	"io"
	"io/ioutil"
)

func formURI(part string) string {
	return fmt.Sprintf("http://localhost.localdomain:0/%s", part)
}

func readError(in io.Reader) error {
	raw, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	return errors.New(string(raw))
}

func writeErrorString(ctx *fasthttp.RequestCtx, e string, code int) {
	log.Errorln(e)
	ctx.Error(e, code)
}
