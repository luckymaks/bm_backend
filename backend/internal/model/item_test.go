package model_test

import (
	"bytes"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/olekukonko/tablewriter"
	"github.com/stretchr/testify/require"

	"github.com/luckymaks/bm_backend/backend/internal/model"
)

func renderGetItemOutputTable(item *model.GetItemOutput) string {
	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.Header("Field", "Value")
	_ = table.Append([]string{"ID", item.ID})
	_ = table.Append([]string{"Data", item.Data})
	_ = table.Append([]string{"CreatedAt", "[NORMALIZED]"})
	_ = table.Render()
	return buf.String()
}

func TestModelCreateItem(t *testing.T) {
	t.Run("returns validation error for empty id", func(t *testing.T) {
		ctx, mdl := setup(t)

		_, err := mdl.CreateItem(ctx, model.CreateItemInput{ID: "", Data: "test"})

		require.Error(t, err)
		require.ErrorIs(t, err, model.ErrValidation)
	})

	t.Run("creates item successfully", func(t *testing.T) {
		ctx, mdl := setup(t)

		output, err := mdl.CreateItem(ctx, model.CreateItemInput{
			ID:   "test-123",
			Data: "hello world",
		})

		require.NoError(t, err)
		require.Equal(t, "test-123", output.ID)
	})

	t.Run("can retrieve created item", func(t *testing.T) {
		ctx, mdl := setup(t)

		_, err := mdl.CreateItem(ctx, model.CreateItemInput{
			ID:   "test-456",
			Data: "stored data",
		})
		require.NoError(t, err)

		item, err := mdl.GetItem(ctx, model.GetItemInput{ID: "test-456"})

		require.NoError(t, err)
		require.Equal(t, "test-456", item.ID)
		require.Equal(t, "stored data", item.Data)
		require.NotEmpty(t, item.CreatedAt)
	})
}

func TestModelGetItem(t *testing.T) {
	t.Run("returns validation error for empty id", func(t *testing.T) {
		ctx, mdl := setup(t)

		_, err := mdl.GetItem(ctx, model.GetItemInput{ID: ""})

		require.Error(t, err)
		require.ErrorIs(t, err, model.ErrValidation)
	})

	t.Run("returns not found for missing item", func(t *testing.T) {
		ctx, mdl := setup(t)

		_, err := mdl.GetItem(ctx, model.GetItemInput{ID: "nonexistent"})

		require.Error(t, err)
		require.ErrorIs(t, err, model.ErrNotFound)
	})
}

func TestGetItemOutput(t *testing.T) {
	t.Parallel()

	t.Run("has correct struct fields", func(t *testing.T) {
		t.Parallel()
		item := model.GetItemOutput{
			ID:        "test-123",
			Data:      "hello world",
			CreatedAt: "2024-01-15T10:00:00Z",
		}

		snapshot := renderGetItemOutputTable(&item)
		snaps.WithConfig(snaps.Dir("snapshot")).MatchSnapshot(t, snapshot)
	})
}
