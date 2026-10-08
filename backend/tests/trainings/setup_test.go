//go:build component

package trainings_test

import (
	"context"
	"os"
	"testing"

	"backend/tests"
)

func TestMain(m *testing.M) {
	tests.StartServer(context.Background())

	os.Exit(m.Run())
}
