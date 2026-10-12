package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"backend/internal/shared"
	"backend/internal/users/app"

	"github.com/labstack/echo/v5"
	common "github.com/quangbach27/golang-common"
	commonHttp "github.com/quangbach27/golang-common/http"
)

// ClerkWebhookPath is where Clerk delivers its events. The route is public: it is authenticated by
// the Svix signature, not by a bearer token.
const ClerkWebhookPath = "/api/v1/webhooks/clerk"

const maxWebhookBodyBytes = 1 << 20

// RoleMirror copies the role of a user into the Clerk public metadata, so it ends up in the
// session token.
type RoleMirror interface {
	MirrorRole(ctx context.Context, userID string, role shared.Role) error
}

type WebhookVerifier interface {
	Verify(payload []byte, headers http.Header) error
}

type WebhookHandler struct {
	service  *app.Service
	verifier WebhookVerifier
}

func RegisterWebhook(globalRouter commonHttp.EchoRouter, webhookHandler *WebhookHandler) error {
	globalRouter.POST(ClerkWebhookPath, webhookHandler.handle)
	return nil
}

// NewClerkWebhook returns the handler for Clerk webhooks. signingSecret is the whsec_... secret of
// the endpoint. roleMirror may be nil to skip mirroring the role to Clerk.
func NewClerkWebhook(service *app.Service, webhookVerifier WebhookVerifier) (*WebhookHandler, error) {
	if service == nil {
		return nil, errors.New("service can't be nil")
	}

	if webhookVerifier == nil {
		return nil, errors.New("webhookVerifier can't be nil")
	}

	return &WebhookHandler{
		service:  service,
		verifier: webhookVerifier,
	}, nil
}

type clerkEvent struct {
	Type string        `json:"type"`
	Data clerkUserData `json:"data"`
}

type clerkUserData struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	PublicMetadata struct {
		Roles []string `json:"roles"`
	} `json:"public_metadata"`
}

func (d clerkUserData) displayName() string {
	if name := strings.TrimSpace(d.FirstName + " " + d.LastName); name != "" {
		return name
	}
	if d.Username != "" {
		return d.Username
	}

	return d.ID
}

// role is the role already assigned in Clerk, or attendee. Trainers can only be set from the
// backend side (public metadata), never by the user.
func (d clerkUserData) role() shared.Role {
	if slices.Contains(d.PublicMetadata.Roles, shared.RoleTrainer.String()) {
		return shared.RoleTrainer
	}

	return shared.RoleAttendee
}

func (w *WebhookHandler) handle(c *echo.Context) error {
	body, err := io.ReadAll(http.MaxBytesReader(c.Response(), c.Request().Body, maxWebhookBodyBytes))
	if err != nil {
		return common.NewInvalidInputError("invalid-webhook-body", "can't read the request body")
	}

	// Verify the signature (and that the timestamp is recent) over the raw body before parsing it.
	if err := w.verifier.Verify(body, c.Request().Header); err != nil {
		return common.NewUnauthorizedError("invalid-webhook-signature", "invalid webhook signature").
			WithInternalError(err)
	}

	var event clerkEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return common.NewInvalidInputError("invalid-webhook-body", "the request body is not valid json")
	}

	ctx := c.Request().Context()

	switch event.Type {
	case "user.created":
		err = w.userCreated(ctx, event.Data)
	case "user.updated":
		err = w.userUpdated(ctx, event.Data)
	default:
		// user.deleted and the rest: acknowledged, so Clerk does not retry. Deleting would orphan
		// the user's trainings, so it is left for a later decision.
		slog.InfoContext(ctx, "ignoring clerk webhook event", "type", event.Type)
	}
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// userCreated is idempotent: Clerk delivers at least once.
func (w *WebhookHandler) userCreated(ctx context.Context, data clerkUserData) error {
	role := data.role()

	if _, err := w.service.EnsureUser(ctx, app.OnboardUserCmd{
		UUID:        data.ID,
		DisplayName: data.displayName(),
		Role:        role,
	}); err != nil {
		return err
	}

	return nil
}

// userUpdated keeps the display name in sync. Events can arrive before user.created or for users
// created before the webhook existed, so an unknown user is created instead.
func (w *WebhookHandler) userUpdated(ctx context.Context, data clerkUserData) error {
	err := w.service.RenameUser(ctx, data.ID, data.displayName())

	var commonErr common.Error
	if errors.As(err, &commonErr) && commonErr.HttpErrorCode == http.StatusNotFound {
		return w.userCreated(ctx, data)
	}

	return err
}
