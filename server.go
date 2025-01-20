package goflex

import (
	"encoding/xml"
	"fmt"
)

// ServerService describes the Server endpoints
type ServerService interface {
	Identity() (*IdentityResponse, error)
	MachineID() (string, error)
}

// ServerServiceOp is the operator for the ServerService
type ServerServiceOp struct {
	p             *Plex
	identityCache *IdentityResponse
}

// MachineID returns the ServerID
func (svc *ServerServiceOp) MachineID() (string, error) {
	got, err := svc.Identity()
	if err != nil {
		return "", err
	}
	return got.MachineIdentifier, nil
}

func (svc *ServerServiceOp) updateIdentityCache() error {
	req := mustNewRequest("GET", fmt.Sprintf("%v/identity", svc.p.baseURL))
	var ret IdentityResponse
	if err := svc.p.sendRequest(req, &ret); err != nil {
		return err
	}
	svc.identityCache = &ret
	return nil
}

// Identity returns the server identity
func (svc *ServerServiceOp) Identity() (*IdentityResponse, error) {
	if svc.identityCache == nil {
		if err := svc.updateIdentityCache(); err != nil {
			return nil, err
		}
	}
	return svc.identityCache, nil
}

// IdentityResponse is the response back from the identity endpoint
type IdentityResponse struct {
	XMLName           xml.Name `xml:"MediaContainer"`
	Text              string   `xml:",chardata"`
	Size              string   `xml:"size,attr"`
	APIVersion        string   `xml:"apiVersion,attr"`
	Claimed           string   `xml:"claimed,attr"`
	MachineIdentifier string   `xml:"machineIdentifier,attr"`
	Version           string   `xml:"version,attr"`
}
