package invite

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestInviteProjectsAreUnordered(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	var schemaResponse resource.SchemaResponse
	(&InviteResource{}).Schema(ctx, resource.SchemaRequest{}, &schemaResponse)

	if _, ok := schemaResponse.Schema.Attributes["projects"].(schema.SetNestedAttribute); !ok {
		t.Fatalf("projects attribute is %T, want schema.SetNestedAttribute", schemaResponse.Schema.Attributes["projects"])
	}

	elementType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":   types.StringType,
		"role": types.StringType,
	}}
	configured, configuredDiags := types.SetValueFrom(ctx, elementType, []InviteProjectModel{
		{ID: types.StringValue("project-a"), Role: types.StringValue("member")},
		{ID: types.StringValue("project-b"), Role: types.StringValue("owner")},
	})
	if configuredDiags.HasError() {
		t.Fatalf("creating configured projects set: %v", configuredDiags)
	}

	apiResponse, responseDiags := types.SetValueFrom(ctx, elementType, []InviteProjectModel{
		{ID: types.StringValue("project-b"), Role: types.StringValue("owner")},
		{ID: types.StringValue("project-a"), Role: types.StringValue("member")},
	})
	if responseDiags.HasError() {
		t.Fatalf("creating API projects set: %v", responseDiags)
	}

	if !configured.Equal(apiResponse) {
		t.Fatal("projects sets differ when only their element order changes")
	}
}
