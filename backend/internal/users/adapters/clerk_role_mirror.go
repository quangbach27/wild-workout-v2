// Package adapters holds the users module's outbound integrations that are not the database.
package adapters

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/clerk/clerk-sdk-go/v2/user"

	"backend/internal/shared"
)

// ClerkRoleMirror writes the role into the Clerk public metadata (only the backend can), which the
// "Customize session token" template exposes as the roles claim.
type ClerkRoleMirror struct{}

func (ClerkRoleMirror) MirrorRole(ctx context.Context, userID string, role shared.Role) error {
	metadata, err := json.Marshal(map[string]any{"roles": []string{role.String()}})
	if err != nil {
		return fmt.Errorf("marshaling clerk metadata: %w", err)
	}

	raw := json.RawMessage(metadata)
	if _, err := user.UpdateMetadata(ctx, userID, &user.UpdateMetadataParams{PublicMetadata: &raw}); err != nil {
		return fmt.Errorf("updating clerk public metadata: %w", err)
	}

	return nil
}
