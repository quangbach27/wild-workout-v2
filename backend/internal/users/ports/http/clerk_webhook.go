package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/labstack/echo/v5"
	common "github.com/quangbach27/golang-common"
	commonHttp "github.com/quangbach27/golang-common/http"
	"github.com/quangbach27/golang-common/log"

	"backend/internal/shared"
	"backend/internal/users/app"
)

// ClerkWebhookPath is where Clerk delivers its events. The route is public: it is authenticated by
// the Svix signature, not by a bearer token; the "/public/" prefix makes the auth middleware skip it.
const ClerkWebhookPath = "/public/api/v1/webhooks/clerk"

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

func RegisterWebhook(router commonHttp.EchoRouter, webhookHandler *WebhookHandler) error {
	router.POST(ClerkWebhookPath, webhookHandler.handle)
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
		Roles []shared.Role `json:"roles"`
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
	if slices.Contains(d.PublicMetadata.Roles, shared.RoleTrainer) {
		return shared.RoleTrainer
	}

	return shared.RoleAttendee
}

func (w *WebhookHandler) handle(c *echo.Context) error {
	logger := log.FromContext(c.Request().Context())

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
	logger.Info("Handle Clerk Event", "event", event.Type)

	ctx := c.Request().Context()

	switch event.Type {
	case "user.created":
		err = w.createUser(ctx, event.Data)
	case "user.updated":
		err = w.updateUser(ctx, event.Data)
	case "user.deleted":
		err = w.service.DeleteUser(ctx, event.Data.ID)
	default:
		logger.Info("ignoring clerk webhook event", "type", event.Type)
	}
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (w *WebhookHandler) createUser(ctx context.Context, data clerkUserData) error {
	role := data.role()

	if err := w.service.OnboardUser(ctx, app.OnboardUserCmd{
		UUID:        data.ID,
		DisplayName: data.displayName(),
		Role:        role,
	}); err != nil {
		return err
	}

	return nil
}

func (w *WebhookHandler) updateUser(ctx context.Context, data clerkUserData) error {
	displayName := data.displayName()
	role := data.role()

	return w.service.UpdateUser(ctx, app.UpdateUserCmd{
		UUID:        data.ID,
		DisplayName: &displayName,
		Role:        &role,
	})
}
