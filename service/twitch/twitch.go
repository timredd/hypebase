package twitch

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/nicklaw5/helix"
	"hypebase/ent"
	"hypebase/ent/servicetwitch"
)

type Twitch struct {
	// AppClientID is the client ID of Hypebase
	AppClientID string

	// UserAccessToken is the access token of the user of the connected Twitch account
	UserAccessToken string

	// State is a unique OAuth 2.0 opaque token for avoidance of CSRF attacks
	State string

	// Client is the API client used to access the Twitch API
	Client *helix.Client

	// DB is the database used by the service
	DB *ent.Client
}

func New(clientID string) (*Twitch, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:    AppClientID,
		RedirectURI: RedirectURI,
	})
	if err != nil {
		return nil, err
	}

	return &Twitch{
		AppClientID: clientID,
		Client:      client,
	}, nil
}

func (t *Twitch) GetAuthorizationURL() (string, error) {
	t.State = genState()

	url := t.Client.GetAuthorizationURL(&helix.AuthorizationURLParams{
		ResponseType: "code",
		Scopes:       DefaultScopes,
		State:        t.State,
		ForceVerify:  false,
	})

	return url, nil
}

func (t *Twitch) GetAccessToken(ctx context.Context, authCode, state string) error {
	if t.State != state {
		return fmt.Errorf("error validating app state; sent: %s, rec: %s", t.State, state)
	}

	resp, err := t.Client.RequestUserAccessToken(authCode)
	if err != nil {
		return err
	}

	t.Client.SetUserAccessToken(resp.Data.AccessToken)

	_, err = t.DB.ServiceTwitch.Create().
		SetAccessToken(resp.Data.AccessToken).
		SetRefreshToken(resp.Data.RefreshToken).
		SetScopes(resp.Data.Scopes).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("error writing access and refresh tokens to DB: %e", err)
	}

	return nil
}

func (t *Twitch) RefreshAccessToken(ctx context.Context) error {
	rec, err := t.DB.ServiceTwitch.Query().
		Where(servicetwitch.AccessToken(t.UserAccessToken)).Only(ctx)
	if err != nil {
		return fmt.Errorf("error reading user's record from db: %e", err)
	} else if rec.AccessToken == "" {
		return fmt.Errorf("error finding user in db: %e", err)
	}

	resp, err := t.Client.RefreshUserAccessToken(rec.RefreshToken)
	if err != nil {
		return err
	}

	t.Client.SetUserAccessToken(resp.Data.AccessToken)
	return nil
}

func (t *Twitch) RevokeAccessToken(ctx context.Context) error {
	rec, err := t.DB.ServiceTwitch.Query().
		Where(servicetwitch.AccessToken(t.UserAccessToken)).Only(ctx)
	if err != nil {
		return fmt.Errorf("error reading user's record from db: %e", err)
	} else if rec.AccessToken == "" {
		return fmt.Errorf("error finding user in db: %e", err)
	}

	_, err = t.Client.RevokeUserAccessToken(rec.AccessToken)
	if err != nil {
		return err
	}

	return nil
}

func (t *Twitch) ValidateAccessToken(ctx context.Context) (bool, error) {
	rec, err := t.DB.ServiceTwitch.Query().
		Where(servicetwitch.AccessToken(t.UserAccessToken)).Only(ctx)
	if err != nil {
		return false, fmt.Errorf("error reading user's record from db: %e", err)
	} else if rec.AccessToken == "" {
		return false, fmt.Errorf("error finding user in db: %e", err)
	}

	isValid, _, err := t.Client.ValidateToken(rec.AccessToken)
	if err != nil {
		return false, err
	}

	return isValid, nil
}

func genState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
