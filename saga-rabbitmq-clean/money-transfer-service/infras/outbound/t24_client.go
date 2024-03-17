package outbound

import (
	"context"
	"encoding/xml"
	"math/rand"

	"github.com/lengocson131002/go-clean/bootstrap"
	"github.com/lengocson131002/go-clean/pkg/trace"
	"github.com/lengocson131002/go-clean/usecase/outbound"
)

const (
	OPEN_CURRENT_ACCOUNT = "OpenCurrentAccount"
)

type t24MqClient struct {
	t24Cfg *bootstrap.T24Config
	// mRepo  data.MasterDataRepository
	// xslt   xslt.Xslt
	tracer trace.Tracer
}

func NewT24MqClient(
	t24Config *bootstrap.T24Config,
	// xslt xslt.Xslt,
	// mRepo data.MasterDataRepository,
	tracer trace.Tracer,
) outbound.T24MQClient {
	return &t24MqClient{
		t24Cfg: t24Config,
		// mRepo:  mRepo,
		// xslt:   xslt,
		tracer: tracer,
	}
}

type t24MQOpenAccountXmlRequest struct {
	XMLName         xml.Name `xml:"ROOT"` // root xml
	CIF             int      `xml:"CIF"`
	AccountTitle    string   `xml:"accountTitle"`
	ShortName       string   `xml:"shortName"`
	Category        string   `xml:"category"`
	RmCode          string   `xml:"rmCode"`
	BranchCode      string   `xml:"branchCode"`
	PostingRestrict string   `xml:"postingRestrict"`
	Program         string   `xml:"program"`
	Currency        string   `xml:"currency"`
	T24User         string   `xml:"t24User"`
}

type templateEntity struct {
	name     string `db:"template_name"`
	request  string `db:"template_request"`
	response string `db:"template_response"`
}

// ExceuteOpenAccount implements outbound.T24MQClient.
func (c *t24MqClient) ExceuteOpenAccount(ctx context.Context, request *outbound.T24MQOpenAccountRequest) (*outbound.T24MQOpenAccountResponse, error) {
	return &outbound.T24MQOpenAccountResponse{
		CIF:    rand.Intn(1000),
		Status: "Success",
	}, nil

}
