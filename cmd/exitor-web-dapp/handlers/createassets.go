package handlers
// To change Createassets to createassets
import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/transaction"

	"exitor-dapp/internal/Createasset"
	"exitor-dapp/internal/platform/auth"
	"exitor-dapp/internal/platform/datatable"
	"exitor-dapp/internal/platform/web"
	"exitor-dapp/internal/platform/web/webcontext"
	"exitor-dapp/internal/platform/web/weberror"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"
)

// Createassets represents the Createasset API method handler set.
type Createassets struct {
	CreateassetRepo *createasset.Repository
	Redis         *redis.Client
	Renderer      web.Renderer
}

func urlCreateassetsIndex() string {
	return fmt.Sprintf("/createassets")
}

func urlCreateassetsCreate() string {
	return fmt.Sprintf("/createassets/create")
}

func urlCreateassetsView(createdassetID string) string {
	return fmt.Sprintf("/createassets/%s", createdassetID)
}

func urlCreateassetsUpdate(createdassetID string) string {
	return fmt.Sprintf("/createassets/%s/update", createdassetID)
}

// Index handles listing all the Createassets for the current account.
func (h *Createassets) Index(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	statusOpts := web.NewEnumResponse(ctx, nil, Createasset.Createassetstatus_ValuesInterface()...)

	statusFilterItems := []datatable.FilterOptionItem{}
	for _, opt := range statusOpts.Options {
		statusFilterItems = append(statusFilterItems, datatable.FilterOptionItem{
			Display: opt.Title,
			Value:   opt.Value,
		})
	}

	// Below, we will transform to represent the parameters required for asset creation on Algorand
	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "name", Title: "Createasset", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
		{Field: "status", Title: "Status", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "All Statuses", FilterItems: statusFilterItems},
		{Field: "updated_at", Title: "Last Updated", Visible: true, Searchable: true, Orderable: true, Filterable: false},
		{Field: "created_at", Title: "Created", Visible: true, Searchable: true, Orderable: true, Filterable: false},
	}

	mapFunc := func(q *Createasset.Createasset, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
		for i := 0; i < len(cols); i++ {
			col := cols[i]
			var v datatable.ColumnValue
			switch col.Field {
			case "id":
				v.Value = fmt.Sprintf("%s", q.ID)
			case "name":
				v.Value = q.Name
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlCreateassetsView(q.ID), v.Value)
			case "status":
				v.Value = q.Status.String()

				var subStatusClass string
				var subStatusIcon string
				switch q.Status {
				case Createasset.Createassetstatus_Active:
					subStatusClass = "text-green"
					subStatusIcon = "far fa-dot-circle"
				case Createasset.Createassetstatus_Disabled:
					subStatusClass = "text-orange"
					subStatusIcon = "far fa-circle"
				}

				v.Formatted = fmt.Sprintf("<span class='cell-font-status %s'><i class='%s mr-1'></i>%s</span>", subStatusClass, subStatusIcon, web.EnumValueTitle(v.Value))
			case "created_at":
				dt := web.NewTimeResponse(ctx, q.CreatedAt)
				v.Value = dt.Local
				v.Formatted = fmt.Sprintf("<span class='cell-font-date'>%s</span>", v.Value)
			case "updated_at":
				dt := web.NewTimeResponse(ctx, q.UpdatedAt)
				v.Value = dt.Local
				v.Formatted = fmt.Sprintf("<span class='cell-font-date'>%s</span>", v.Value)
			default:
				return resp, errors.Errorf("Failed to map value for %s.", col.Field)
			}
			resp = append(resp, v)
		}

		return resp, nil
	}

	loadFunc := func(ctx context.Context, sorting string, fields []datatable.DisplayField) (resp [][]datatable.ColumnValue, err error) {
		res, err := h.CreateassetRepo.Find(ctx, claims, Createasset.CreateassetFindRequest{
			Where: "account_id = ?",
			Args:  []interface{}{claims.Audience},
			Order: strings.Split(sorting, ","),
		})
		if err != nil {
			return resp, err
		}

		for _, a := range res {
			l, err := mapFunc(a, fields)
			if err != nil {
				return resp, errors.Wrapf(err, "Failed to map Createasset for display.")
			}

			resp = append(resp, l)
		}

		return resp, nil
	}

	dt, err := datatable.New(ctx, w, r, h.Redis, fields, loadFunc)
	if err != nil {
		return err
	}

	if dt.HasCache() {
		return nil
	}

	if ok, err := dt.Render(); ok {
		if err != nil {
			return err
		}
		return nil
	}

	data := map[string]interface{}{
		"datatable":           dt.Response(),
		"urlCreateassetsCreate": urlCreateassetsCreate(),
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "createassets-index.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

const algodAddress = 
const psToken = ""

// Create a throw-away account for this example - 
// check that it has funds before running the program
const mn = "..." // To include mnemonic
const ownerAddress = "..."// can also be derived from mnemonic, I will hardcode to make it easier

// Will create an algod client in the main.go file
// It has to pass

func () initAlgorand()

func () assetAlgo(walletAddress, assetName, unitName, supp)

// Create handles creating a new Asset for the account.
// Also includes the Algorand Implementation for Asset Creation which has to succeed as well
func (h *Createassets) Create(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	//
	req := new(Createasset.CreateassetCreateRequest)
	data := make(map[string]interface{})
	f := func() (bool, error) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				return false, err
			}

			decoder := schema.NewDecoder()
			decoder.IgnoreUnknownKeys(true)

			if err := decoder.Decode(req, r.PostForm); err != nil {
				return false, err
			}
			req.AccountID = claims.Audience

			usr, err := h.CreateassetRepo.Create(ctx, claims, *req, ctxValues.Now)
			if err != nil {
				switch errors.Cause(err) {
				default:
					if verr, ok := weberror.NewValidationError(ctx, err); ok {
						data["validationErrors"] = verr.(*weberror.Error)
						return false, nil
					} else {
						return false, err
					}
				}
			}

			// We make a transaction during the createasset operation
			coinTotalIssuance := uint64(1000000)
			coinDecimalsForDisplay := uint32(0)
			accountsAreDefaultFrozen := false
			managerAddress := ownerAddress
			assetReserveAddress := ""
			addressWithFreezingPrivileges := ownerAddress
			addressWithClawbackPrivileges := ownerAddress
			assetUnitName :+ "biztoken"
			assetUrl := ""
			assetMetadataHash := ""
			tx, err := transaction.MakeAssetCreateTxn(ownerAddress, txParams.Fee, txParams.LastRound, txParams.LastRound+10, nil,
			txParams.GenesisID, base64.stdEncoding.EncodeToString(txParams.GenesisHash),
			coinTotalIssuance, coinDecimalsForDisplay, accountsAreDefaultFrozen, managerAddress,
			assetReserveAddress, addressWithFreezingPrivileges, addressWithClawbackPrivileges,
			assetUnitName, assetName, assetUrl, assetMetadataHash)

			if err != nil {
				fmt.Printf("Error creating transaction: %s\n", err)
				return
			}

			// Sign the Transaction
			_, bytes, err := crypto.SignTran

			// Display a success message to the Createasset.
			webcontext.SessionFlashSuccess(ctx,
				"Createasset Created",
				"Createasset successfully created.")

			return true, web.Redirect(ctx, w, r, urlCreateassetsView(usr.ID), http.StatusFound)
		}

		return false, nil
	}

	end, err := f()
	if err != nil {
		return web.RenderError(ctx, w, r, err, h.Renderer, TmplLayoutBase, TmplContentErrorGeneric, web.MIMETextHTMLCharsetUTF8)
	} else if end {
		return nil
	}

	data["form"] = req

	if verr, ok := weberror.NewValidationError(ctx, webcontext.Validator().Struct(Createasset.CreateassetCreateRequest{})); ok {
		data["validationDefaults"] = verr.(*weberror.Error)
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "Createassets-create.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// View handles displaying a Createasset.
func (h *Createassets) View(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	CreateassetID := params["Createasset_id"]

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	f := func() (bool, error) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				return false, err
			}

			switch r.PostForm.Get("action") {
			case "archive":
				err = h.CreateassetRepo.Archive(ctx, claims, Createasset.CreateassetArchiveRequest{
					ID: CreateassetID,
				}, ctxValues.Now)
				if err != nil {
					return false, err
				}

				webcontext.SessionFlashSuccess(ctx,
					"Createasset Archive",
					"Createasset successfully archive.")

				return true, web.Redirect(ctx, w, r, urlCreateassetsIndex(), http.StatusFound)
			}
		}

		return false, nil
	}

	end, err := f()
	if err != nil {
		return web.RenderError(ctx, w, r, err, h.Renderer, TmplLayoutBase, TmplContentErrorGeneric, web.MIMETextHTMLCharsetUTF8)
	} else if end {
		return nil
	}

	prj, err := h.CreateassetRepo.ReadByID(ctx, claims, CreateassetID)
	if err != nil {
		return err
	}
	data["Createasset"] = prj.Response(ctx)
	data["urlCreateassetsView"] = urlCreateassetsView(CreateassetID)
	data["urlCreateassetsUpdate"] = urlCreateassetsUpdate(CreateassetID)

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "Createassets-view.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// Update handles updating a Createasset for the account.
func (h *Createassets) Update(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	CreateassetID := params["Createasset_id"]

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	//
	req := new(Createasset.CreateassetUpdateRequest)
	data := make(map[string]interface{})
	f := func() (bool, error) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				return false, err
			}

			decoder := schema.NewDecoder()
			decoder.IgnoreUnknownKeys(true)

			if err := decoder.Decode(req, r.PostForm); err != nil {
				return false, err
			}
			req.ID = CreateassetID

			err = h.CreateassetRepo.Update(ctx, claims, *req, ctxValues.Now)
			if err != nil {
				switch errors.Cause(err) {
				default:
					if verr, ok := weberror.NewValidationError(ctx, err); ok {
						data["validationErrors"] = verr.(*weberror.Error)
						return false, nil
					} else {
						return false, err
					}
				}
			}

			// Display a success message to the Createasset.
			webcontext.SessionFlashSuccess(ctx,
				"Createasset Updated",
				"Createasset successfully updated.")

			return true, web.Redirect(ctx, w, r, urlCreateassetsView(req.ID), http.StatusFound)
		}

		return false, nil
	}

	end, err := f()
	if err != nil {
		return web.RenderError(ctx, w, r, err, h.Renderer, TmplLayoutBase, TmplContentErrorGeneric, web.MIMETextHTMLCharsetUTF8)
	} else if end {
		return nil
	}

	prj, err := h.CreateassetRepo.ReadByID(ctx, claims, CreateassetID)
	if err != nil {
		return err
	}
	data["Createasset"] = prj.Response(ctx)

	data["urlCreateassetsView"] = urlCreateassetsView(CreateassetID)

	if req.ID == "" {
		req.Name = &prj.Name
		req.Status = &prj.Status
	}
	data["form"] = req

	if verr, ok := weberror.NewValidationError(ctx, webcontext.Validator().Struct(Createasset.CreateassetUpdateRequest{})); ok {
		data["validationDefaults"] = verr.(*weberror.Error)
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "Createassets-update.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}
