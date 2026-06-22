package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func textMax(n int) *int { return &n }

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)
		if err := createOrganizations(dao); err != nil {
			return err
		}
		if err := extendUsers(dao); err != nil {
			return err
		}
		if err := createCustomers(dao); err != nil {
			return err
		}
		return createOrdersAndWindows(dao)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)
		for _, name := range []string{"window_items", "measuring_orders", "customers", "organizations"} {
			c, err := dao.FindCollectionByNameOrId(name)
			if err == nil {
				if err := dao.DeleteCollection(c); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

var windowTypes = []string{
	"single_fixed",
	"single_tilt_turn",
	"single_side_hung",
	"single_tilt_only",
	"double_fixed_operable",
	"double_two_tilt_turn",
	"double_french_casement",
	"triple_fixed_center_operable_fixed",
	"triple_operable_fixed_operable",
	"triple_mixed",
	"multi_4plus",
	"multi_transom",
	"multi_mullion",
}

func createOrganizations(dao *daos.Dao) error {
	if _, err := dao.FindCollectionByNameOrId("organizations"); err == nil {
		return nil
	}

	collection := &models.Collection{}
	collection.Name = "organizations"
	collection.Type = models.CollectionTypeBase
	collection.ListRule = types.Pointer(`@request.auth.id != ""`)
	collection.ViewRule = types.Pointer(`@request.auth.id != ""`)
	collection.CreateRule = types.Pointer(`@request.auth.id != ""`)
	collection.UpdateRule = types.Pointer(`@request.auth.id != ""`)
	collection.DeleteRule = types.Pointer(`@request.auth.id != "" && @request.auth.role = "admin"`)

	collection.Schema = schema.NewSchema(
		&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true, Options: &schema.TextOptions{Max: textMax(200)}},
		&schema.SchemaField{Name: "slug", Type: schema.FieldTypeText, Required: true, Options: &schema.TextOptions{Max: textMax(100)}},
		&schema.SchemaField{Name: "branding", Type: schema.FieldTypeJson},
		&schema.SchemaField{Name: "license_key", Type: schema.FieldTypeText, Options: &schema.TextOptions{Max: textMax(200)}},
	)
	collection.Indexes = types.JsonArray[string]{
		"CREATE UNIQUE INDEX idx_organizations_slug ON organizations (slug)",
	}

	return dao.SaveCollection(collection)
}

func extendUsers(dao *daos.Dao) error {
	collection, err := dao.FindCollectionByNameOrId("users")
	if err != nil {
		return err
	}

	if collection.Schema.GetFieldByName("organization") != nil {
		return nil
	}

	orgs, err := dao.FindCollectionByNameOrId("organizations")
	if err != nil {
		return err
	}

	collection.Schema.AddField(&schema.SchemaField{
		Name:     "organization",
		Type:     schema.FieldTypeRelation,
		Required: false,
		Options: &schema.RelationOptions{
			CollectionId:  orgs.Id,
			CascadeDelete: false,
			MaxSelect:     types.Pointer(1),
		},
	})
	collection.Schema.AddField(&schema.SchemaField{
		Name: "role",
		Type: schema.FieldTypeSelect,
		Options: &schema.SelectOptions{
			Values:    []string{"field", "reviewer", "admin"},
			MaxSelect: 1,
		},
	})

	return dao.SaveCollection(collection)
}

func createCustomers(dao *daos.Dao) error {
	if _, err := dao.FindCollectionByNameOrId("customers"); err == nil {
		return nil
	}

	orgs, err := dao.FindCollectionByNameOrId("organizations")
	if err != nil {
		return err
	}

	collection := &models.Collection{}
	collection.Name = "customers"
	collection.Type = models.CollectionTypeBase
	collection.ListRule = types.Pointer(`@request.auth.id != "" && organization = @request.auth.organization`)
	collection.ViewRule = types.Pointer(`@request.auth.id != "" && organization = @request.auth.organization`)
	collection.CreateRule = types.Pointer(`@request.auth.id != "" && @request.data.organization = @request.auth.organization`)
	collection.UpdateRule = types.Pointer(`@request.auth.id != "" && organization = @request.auth.organization`)
	collection.DeleteRule = types.Pointer(`@request.auth.id != "" && @request.auth.role = "admin"`)

	collection.Schema = schema.NewSchema(
		&schema.SchemaField{
			Name: "organization", Type: schema.FieldTypeRelation, Required: true,
			Options: &schema.RelationOptions{CollectionId: orgs.Id, MaxSelect: types.Pointer(1)},
		},
		&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true, Options: &schema.TextOptions{Max: textMax(200)}},
		&schema.SchemaField{Name: "address", Type: schema.FieldTypeText, Options: &schema.TextOptions{Max: textMax(500)}},
		&schema.SchemaField{Name: "contact", Type: schema.FieldTypeText, Options: &schema.TextOptions{Max: textMax(200)}},
		&schema.SchemaField{Name: "email", Type: schema.FieldTypeEmail},
		&schema.SchemaField{Name: "phone", Type: schema.FieldTypeText, Options: &schema.TextOptions{Max: textMax(50)}},
		&schema.SchemaField{Name: "notes", Type: schema.FieldTypeText, Options: &schema.TextOptions{Max: textMax(2000)}},
	)
	collection.Indexes = types.JsonArray[string]{
		"CREATE INDEX idx_customers_org ON customers (organization)",
	}

	return dao.SaveCollection(collection)
}

func createOrdersAndWindows(dao *daos.Dao) error {
	orgs, err := dao.FindCollectionByNameOrId("organizations")
	if err != nil {
		return err
	}
	customers, err := dao.FindCollectionByNameOrId("customers")
	if err != nil {
		return err
	}
	users, err := dao.FindCollectionByNameOrId("users")
	if err != nil {
		return err
	}

	if _, err := dao.FindCollectionByNameOrId("measuring_orders"); err != nil {
		orders := &models.Collection{}
		orders.Name = "measuring_orders"
		orders.Type = models.CollectionTypeBase
		orders.ListRule = types.Pointer(`@request.auth.id != "" && organization = @request.auth.organization`)
		orders.ViewRule = types.Pointer(`@request.auth.id != "" && organization = @request.auth.organization`)
		orders.CreateRule = types.Pointer(`@request.auth.id != "" && @request.data.organization = @request.auth.organization`)
		orders.UpdateRule = types.Pointer(`@request.auth.id != "" && organization = @request.auth.organization`)
		orders.DeleteRule = types.Pointer(`@request.auth.id != "" && @request.auth.role = "admin"`)

		orders.Schema = schema.NewSchema(
			&schema.SchemaField{
				Name: "organization", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: orgs.Id, MaxSelect: types.Pointer(1)},
			},
			&schema.SchemaField{
				Name: "customer", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: customers.Id, MaxSelect: types.Pointer(1)},
			},
			&schema.SchemaField{
				Name: "status", Type: schema.FieldTypeSelect,
				Options: &schema.SelectOptions{
					Values:    []string{"draft", "in_progress", "review", "exported", "archived"},
					MaxSelect: 1,
				},
			},
			&schema.SchemaField{Name: "notes", Type: schema.FieldTypeText, Options: &schema.TextOptions{Max: textMax(5000)}},
			&schema.SchemaField{
				Name: "created_by", Type: schema.FieldTypeRelation,
				Options: &schema.RelationOptions{CollectionId: users.Id, MaxSelect: types.Pointer(1)},
			},
		)
		orders.Indexes = types.JsonArray[string]{
			"CREATE INDEX idx_orders_org ON measuring_orders (organization)",
		}
		if err := dao.SaveCollection(orders); err != nil {
			return err
		}
	}

	orderColl, err := dao.FindCollectionByNameOrId("measuring_orders")
	if err != nil {
		return err
	}

	if _, err := dao.FindCollectionByNameOrId("window_items"); err == nil {
		return nil
	}

	windows := &models.Collection{}
	windows.Name = "window_items"
	windows.Type = models.CollectionTypeBase
	windows.ListRule = types.Pointer(`@request.auth.id != ""`)
	windows.ViewRule = types.Pointer(`@request.auth.id != ""`)
	windows.CreateRule = types.Pointer(`@request.auth.id != ""`)
	windows.UpdateRule = types.Pointer(`@request.auth.id != ""`)
	windows.DeleteRule = types.Pointer(`@request.auth.id != ""`)

	windows.Schema = schema.NewSchema(
		&schema.SchemaField{
			Name: "order", Type: schema.FieldTypeRelation, Required: true,
			Options: &schema.RelationOptions{
				CollectionId: orderColl.Id, MaxSelect: types.Pointer(1), CascadeDelete: true,
			},
		},
		&schema.SchemaField{Name: "label", Type: schema.FieldTypeText, Options: &schema.TextOptions{Max: textMax(100)}},
		&schema.SchemaField{
			Name: "window_type", Type: schema.FieldTypeSelect, Required: true,
			Options: &schema.SelectOptions{Values: windowTypes, MaxSelect: 1},
		},
		&schema.SchemaField{
			Name: "status", Type: schema.FieldTypeSelect,
			Options: &schema.SelectOptions{
				Values:    []string{"pending", "photos", "processing", "review", "done"},
				MaxSelect: 1,
			},
		},
		&schema.SchemaField{
			Name: "inner_photos", Type: schema.FieldTypeFile,
			Options: &schema.FileOptions{MaxSelect: 10, MaxSize: 10485760},
		},
		&schema.SchemaField{
			Name: "outer_photos", Type: schema.FieldTypeFile,
			Options: &schema.FileOptions{MaxSelect: 10, MaxSize: 10485760},
		},
		&schema.SchemaField{Name: "cv_result", Type: schema.FieldTypeJson},
		&schema.SchemaField{Name: "overrides", Type: schema.FieldTypeJson},
		&schema.SchemaField{Name: "confidence", Type: schema.FieldTypeNumber},
	)
	windows.Indexes = types.JsonArray[string]{
		"CREATE INDEX idx_windows_order ON window_items (order)",
	}

	return dao.SaveCollection(windows)
}
