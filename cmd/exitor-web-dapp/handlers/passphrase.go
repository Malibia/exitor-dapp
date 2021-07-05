package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"

	"exitor-dapp/internal/platform/web"
	"exitor-dapp/internal/platform/web/webcontext"
	"exitor-dapp/internal/platform/web/weberror"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

type SignTx struct {
	Renderer web.Renderer
}

// Prompt the user for a passphrase
func (s SignTx) FlashForPassphrase(ctx contex.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	if r.Method == http.MethodGet {
		webcontext.SessionFlashInfo(ctx,
		"Please enter your Algorand Account Passphrase")
	}

	data := make(map[string]interface{})

	type SignAssetTx struct {
		walletAddress string `json:"algorand_wallet_address" validate:"required"`
		passphrase string `json:"passphrase" validate:"required"`
	}

	req := new(SignAssetTx)
	f := func() error {

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
					return err
			}

			decoder := schema.NewDecoder()
			if err := decoder.Decode(req, r.PostForm); err != nil {
					return err
			}

			if err := webcontext.Validator().Struct(req); err != nil {
					if ne, ok := weberror.NewValidationError(ctx, err); ok {
							data["validationErrors"] = ne.(*weberror.Error)
							return nil
					} else {
							return err
					}
			}
		}

		return nil
	}

	if err := f(); err != nil {
			return web.RenderError(ctx, w, r, err, h.Renderer, TmplLayoutBase, TmplContentErrorGeneric, web.MIMETextHTMLCharsetUTF8)
	}

	data["form"] = req

	if verr, ok := weberror.NewValidationError(ctx, webcontext.Validator().Struct(SignAssetTx{})); ok {
			data["validationDefaults"] = verr.(*weberror.Error)
	}

	return h.Renderer.Render(ctx, w, r,TmplLayoutBase, "flash-for-passphrase-message.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}
/*// Example represents the example pages
type Examples struct {
	Renderer web.Renderer
}*/

