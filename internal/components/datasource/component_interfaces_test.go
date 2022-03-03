//go:build integration
// +build integration

package datasource

import (
	"context"
	"strconv"
	"testing"

	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/sqlstore"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"
)

func TestStoreDSStoreCRUD(t *testing.T) {
	t.Run("Insert / Update / Delete with Gets", func(t *testing.T) {
		sqlStore := sqlstore.InitTestDB(t)
		ctx := context.Background()
		dsStore := ProvideDataSourceSchemaStore(sqlStore)

		uid := "MySpecialUIDisDEARtoMe"
		jd := make(map[string]interface{})
		jd["test"] = "test"

		modelToInsert := Datasource{
			Spec: Model{
				JsonData: jd,
			},
		}
		modelToInsert.Name = "Test"
		modelToInsert.UID = types.UID(uid)

		// Insert
		err := dsStore.Insert(ctx, modelToInsert)
		require.NoError(t, err)

		// Get
		fetched, err := dsStore.Get(ctx, uid)
		require.NoError(t, err)

		fetchedDS, ok := fetched.(Datasource)
		require.True(t, ok)

		modelToInsertWithVersionBumped := Datasource{
			Spec: Model{
				JsonData: jd,
			},
		}
		modelToInsertWithVersionBumped.Name = "Test"
		modelToInsertWithVersionBumped.UID = types.UID(uid)
		modelToInsertWithVersionBumped.ResourceVersion = strconv.Itoa(1)

		require.Equal(t, modelToInsertWithVersionBumped, fetchedDS)

		// Update
		modelForUpdate := Datasource{
			Spec: Model{
				JsonData: jd,
				Type:     "slothFactory",
			},
		}
		modelForUpdate.Name = "Test"
		modelForUpdate.UID = types.UID(uid)
		modelForUpdate.ResourceVersion = fetchedDS.ResourceVersion // We are manually setting version

		err = dsStore.Update(ctx, modelForUpdate)
		require.NoError(t, err)

		// Get updated
		modelForUpdateWithVersionBump := Datasource{
			Spec: Model{
				JsonData: jd,
				Type:     "slothFactory",
			},
		}
		modelForUpdateWithVersionBump.Name = "Test"
		modelForUpdateWithVersionBump.UID = types.UID(uid)
		rv, err := strconv.Atoi(fetchedDS.ResourceVersion)
		require.NoError(t, err)
		modelForUpdateWithVersionBump.ResourceVersion = strconv.Itoa(rv + 1) // We are manually setting version

		fetchedUpdatedDS, err := dsStore.Get(ctx, uid)
		require.NoError(t, err)
		require.Equal(t, modelForUpdateWithVersionBump, fetchedUpdatedDS)

		// Delete it
		err = dsStore.Delete(ctx, uid)
		require.NoError(t, err)

		_, err = dsStore.Get(ctx, uid)
		require.ErrorIs(t, err, models.ErrDataSourceNotFound)
	})
}