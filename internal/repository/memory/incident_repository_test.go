package memory

import (
	"context"
	"testing"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIncidentRepositoryMemory_AddingAndGettingIncident(t *testing.T) {
	clock := mocks.FixedClock{}
	repo := &IncidentRepositoryMemory{
		Clock: clock,
	}

	channelID := ref.ChannelID("e27ddcd0-0e1f-4bc5-93df-f6f04155beec")
	actorID := ref.ExternalUserUUID("f49d5fd5-8da4-4779-b5ba-32e78aa2c444")
	ctx := context.Background()

	inc1 := incident.Incident{
		Number:           "ABC123",
		ExternalID:       "some external ID",
		ShortDescription: "some short description",
		Description:      "some description",
	}
	err := inc1.CreatedUpdated.SetCreatedBy(actorID)
	require.NoError(t, err)
	err = inc1.CreatedUpdated.SetUpdatedBy(actorID)
	require.NoError(t, err)

	incID, err := repo.AddIncident(ctx, channelID, inc1)
	require.NoError(t, err)

	retInc, err := repo.GetIncident(ctx, channelID, incID)
	require.NoError(t, err)

	assert.Equal(t, incID, retInc.UUID())
	assert.Equal(t, inc1.Number, retInc.Number)
	assert.Equal(t, inc1.ExternalID, retInc.ExternalID)
	assert.Equal(t, inc1.ShortDescription, retInc.ShortDescription)
	assert.Equal(t, inc1.Description, retInc.Description)

	// test correct timestamps
	assert.NotEmpty(t, inc1.CreatedUpdated.CreatedBy())
	assert.Equal(t, inc1.CreatedUpdated.CreatedBy(), retInc.CreatedUpdated.CreatedBy())
	assert.Equal(t, clock.NowFormatted(), retInc.CreatedUpdated.CreatedAt())

	assert.NotEmpty(t, inc1.CreatedUpdated.UpdatedBy())
	assert.Equal(t, inc1.CreatedUpdated.UpdatedBy(), retInc.CreatedUpdated.UpdatedBy())
	assert.Equal(t, clock.NowFormatted(), retInc.CreatedUpdated.UpdatedAt())
}

func TestIncidentRepositoryMemory_ListIncidents(t *testing.T) {
	clock := mocks.FixedClock{}
	repo := &IncidentRepositoryMemory{
		Clock: clock,
	}

	channelID := ref.ChannelID("e27ddcd0-0e1f-4bc5-93df-f6f04155beec")
	actorID := ref.ExternalUserUUID("f49d5fd5-8da4-4779-b5ba-32e78aa2c444")
	actor2ID := ref.ExternalUserUUID("00271cb4-3716-4203-9124-1d2f515ae0b2")

	ctx := context.Background()

	inc1 := incident.Incident{
		Number:           "Bca258",
		ExternalID:       "some external ID",
		ShortDescription: "some short description",
		Description:      "some description",
	}
	err := inc1.CreatedUpdated.SetCreatedBy(actorID)
	require.NoError(t, err)
	err = inc1.CreatedUpdated.SetUpdatedBy(actorID)
	require.NoError(t, err)

	inc2 := incident.Incident{
		Number:           "CDB36478",
		ExternalID:       "some external ID 2",
		ShortDescription: "some short description 2",
		Description:      "some description 2",
	}
	err = inc2.CreatedUpdated.SetCreatedBy(actorID)
	require.NoError(t, err)
	err = inc2.CreatedUpdated.SetUpdatedBy(actor2ID)
	require.NoError(t, err)

	_, err = repo.AddIncident(ctx, channelID, inc1)
	require.NoError(t, err)

	_, err = repo.AddIncident(ctx, channelID, inc2)
	require.NoError(t, err)

	// first page
	incidentsList, err := repo.ListIncidents(ctx, channelID, 1, 10)
	require.NoError(t, err)

	list := incidentsList.Result

	assert.Len(t, list, 2)

	for i, retInc := range list {
		var inc incident.Incident
		switch i {
		case 0:
			inc = inc1
		case 1:
			inc = inc2
		}

		assert.NotEmpty(t, retInc.UUID)
		assert.Equal(t, inc.Number, retInc.Number)
		assert.Equal(t, inc.ExternalID, retInc.ExternalID)
		assert.Equal(t, inc.ShortDescription, retInc.ShortDescription)
		assert.Equal(t, inc.Description, retInc.Description)

		// test correct timestamps
		assert.NotEmpty(t, inc.CreatedUpdated.CreatedBy())
		assert.Equal(t, inc.CreatedUpdated.CreatedBy(), retInc.CreatedUpdated.CreatedBy())
		assert.Equal(t, clock.NowFormatted(), retInc.CreatedUpdated.CreatedAt())

		assert.NotEmpty(t, inc.CreatedUpdated.UpdatedBy())
		assert.Equal(t, inc.CreatedUpdated.UpdatedBy(), retInc.CreatedUpdated.UpdatedBy())
		assert.Equal(t, clock.NowFormatted(), retInc.CreatedUpdated.UpdatedAt())
	}

	// second page out of range
	incidentsList, err = repo.ListIncidents(ctx, channelID, 2, 10)
	require.NoError(t, err)

	list = incidentsList.Result
	assert.Len(t, list, 0)

	// first page with small number per page
	incidentsList, err = repo.ListIncidents(ctx, channelID, 1, 1)
	require.NoError(t, err)

	list = incidentsList.Result
	assert.Len(t, list, 1)
}
