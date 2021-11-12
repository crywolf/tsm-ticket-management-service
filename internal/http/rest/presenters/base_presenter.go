package presenters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"
	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters/hypermedia"
	"go.uber.org/zap"
)

// NewBasePresenter returns presentation service with basic functionality
func NewBasePresenter(logger *zap.SugaredLogger, serverAddr string) *BasePresenter {
	return &BasePresenter{
		logger:     logger,
		serverAddr: serverAddr,
	}
}

// BasePresenter must be included in all derived presenters via object composition
type BasePresenter struct {
	logger     *zap.SugaredLogger
	serverAddr string
}

// WriteError replies to the request with the specified error message and HTTP code.
func (p BasePresenter) WriteError(w http.ResponseWriter, error string, code int) {
	p.sendErrorJSON(w, error, code)
}

func (p BasePresenter) resourceToHypermediaLinks(hypermediaMapper hypermedia.Mapper, domainObject hypermedia.ActionsMapper) api.HypermediaLinks {
	hypermediaLinks := api.HypermediaLinks{}

	actions := hypermediaMapper.RoutesToHypermediaActionLinks()
	allowedActions := domainObject.AllowedActions()
	for _, action := range allowedActions {
		link := actions[action]
		href := strings.ReplaceAll(link.Href, "{uuid}", domainObject.UUID().String())
		hypermediaLinks[link.Name] = map[string]string{
			"href": href,
		}
	}

	return hypermediaLinks
}

// sendErrorJSON replies to the request with the specified error message and HTTP code.
// It encodes error string as JSON object {"error":"error_string"} and sets correct header.
// It does not otherwise end the request; the caller should ensure no further writes are done to 'w'.
// The error message should be plain text.
func (p BasePresenter) sendErrorJSON(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	errorJSON, _ := json.Marshal(error)
	_, _ = fmt.Fprintf(w, `{"error":%s}`+"\n", errorJSON)
}

// encodeJSON encodes 'v' to JSON and writes it to the 'w'. Also sets correct Content-Type header.
// It does not otherwise end the request; the caller should ensure no further writes are done to 'w'.
func (p BasePresenter) encodeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		eMsg := "could not encode JSON response"
		p.logger.Errorw(eMsg, "error", err)
		p.WriteError(w, eMsg, http.StatusInternalServerError)
		return
	}
}
