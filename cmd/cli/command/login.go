package command

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/keyring"
	"github.com/common-fate/ciem/tokenstore"
	"github.com/common-fate/clio"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
)

var Login = cli.Command{
	Name:  "login",
	Usage: "Log in to Common Fate Cloud",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "lazy", Usage: "When the lazy flag is used, a login flow will only be started when the access token is expired"},
		&cli.StringFlag{Name: "issuer", Usage: "The OIDC issuer"},
		&cli.StringFlag{Name: "client-id", Usage: "The OIDC client ID"},
	},
	Action: func(c *cli.Context) error {
		lf := LoginFlow{
			ClientID: c.String("client-id"),
			Issuer:   c.String("issuer"),
		}

		return lf.Login(c.Context)
	},
}

type LoginFlow struct {
	// Keyring optionally overrides the keyring that auth tokens are saved to.
	Keyring  keyring.Keyring
	ClientID string
	Issuer   string
}

type Response struct {
	// Err is set if there was an error which
	// prevented the flow from completing
	Err   error
	Token (*oauth2.Token)
}

func (lf LoginFlow) Login(ctx context.Context) error {
	// oldCfg, err := cliconfig.Load()
	// if err != nil {
	// 	return err
	// }
	// oldDefaultContext := oldCfg.CurrentOrEmpty()

	authResponse := make(chan Response)

	var g errgroup.Group

	emptyClientSecret := ""
	scopes := []string{"openid", "email"}
	redirectURI := "http://localhost:18900/auth/callback"
	options := []rp.Option{
		rp.WithVerifierOpts(rp.WithIssuedAtOffset(5 * time.Second)),
	}

	provider, err := rp.NewRelyingPartyOIDC(lf.Issuer, lf.ClientID, emptyClientSecret, redirectURI, scopes, options...)
	if err != nil {
		return fmt.Errorf("error creating provider %s", err.Error())
	}

	// create random state variable for OIDC flow
	state := func() string {
		return uuid.New().String()
	}

	tokenWriter := func(w http.ResponseWriter, r *http.Request, tokens *oidc.Tokens[*oidc.IDTokenClaims], state string, rp rp.RelyingParty) {
		authResponse <- Response{
			Token: tokens.Token,
		}
		_, _ = w.Write([]byte("success! You can now close this window"))
	}

	r := chi.NewRouter()
	r.Handle("/auth/callback", rp.CodeExchangeHandler(tokenWriter, provider))
	r.Handle("/login", rp.AuthURLHandler(state, provider, rp.WithPromptURLParam("Welcome back!")))
	server := &http.Server{
		Addr:    ":18900",
		Handler: r,
	}

	// run the auth server on localhost
	g.Go(func() error {
		clio.Debugw("starting HTTP server", "address", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
		clio.Debugw("auth server closed")
		return nil
	})

	// read the returned ID token from Cognito
	g.Go(func() error {
		res := <-authResponse

		err := server.Shutdown(ctx)
		if err != nil {
			return err
		}

		// check that the auth flow didn't error out
		if res.Err != nil {
			return err
		}

		ts := tokenstore.New(tokenstore.WithKeyring(lf.Keyring))
		err = ts.Save(res.Token)
		if err != nil {
			return err
		}

		clio.Successf("Successfully logged in")

		return nil
	})

	// open the browser and read the token
	g.Go(func() error {
		url := "http://localhost:18900/login"
		clio.Infof("Opening your web browser to: %s", url)
		err := browser.OpenURL(url)
		if err != nil {
			clio.Errorf("error opening browser: %s", err)
		}
		return nil
	})

	err = g.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (lf *LoginFlow) oidcProvider() (rp.RelyingParty, error) {
	emptyClientSecret := ""
	scopes := []string{"openid", "email"}
	redirectURI := "http://localhost:18900/auth/callback"
	options := []rp.Option{
		rp.WithVerifierOpts(rp.WithIssuedAtOffset(5 * time.Second)),
	}

	return rp.NewRelyingPartyOIDC(lf.Issuer, lf.ClientID, emptyClientSecret, redirectURI, scopes, options...)
}

func (lf *LoginFlow) RefreshToken(ctx context.Context) (*oauth2.Token, error) {
	ts := tokenstore.New(tokenstore.WithKeyring(lf.Keyring))
	tok, err := ts.Token()
	if err != nil {
		return nil, err
	}

	provider, err := lf.oidcProvider()
	if err != nil {
		return nil, err
	}

	c := provider.OAuthConfig()
	tok.Expiry = time.Now().Add(-time.Second * 10)
	src := tokenstore.NotifyRefreshTokenSource{
		New:       c.TokenSource(ctx, tok),
		T:         tok,
		SaveToken: ts.Save,
	}
	return src.Token()
}
