package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io"
	"log"
	"net/http"
	"terraform-provider-trust-store/tools"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &trustStoreResource{}
	_ resource.ResourceWithConfigure = &trustStoreResource{}
)

func NewTrustStoreResource() resource.Resource {
	return &trustStoreResource{}
}

type trustStoreResource struct {
	client *tools.TrustStoreClient
}

type trustStoreModel struct {
	ID           types.String `tfsdk:"id"`
	SerialNumber types.String `tfsdk:"serial_number"`
	Certificate  types.String `tfsdk:"certificate"`
	Status       types.String `tfsdk:"status"`
	Issuer       types.String `tfsdk:"issuer"`
	Signature    types.String `tfsdk:"signature"`
	UploadedOn   types.String `tfsdk:"uploaded_on"`
	UploadedAt   types.String `tfsdk:"uploaded_at"`
	ExpiresOn    types.String `tfsdk:"expires_on"`
}

type trustStoreJsonModel struct {
	Result struct {
		ID           string `json:"id"`
		SerialNumber string `json:"serial_number"`
		Certificate  string `json:"certificate"`
		Status       string `json:"status"`
		Issuer       string `json:"issuer"`
		Signature    string `json:"signature"`
		UploadedOn   string `json:"uploaded_on"`
		UploadedAt   string `json:"uploaded_at"`
		ExpiresOn    string `json:"expires_on"`
	}
}

// Configure adds the provider configured client to the resource.
func (r *trustStoreResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*tools.TrustStoreClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Metadata returns the resource type name.
func (r *trustStoreResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_origin"
}

// Schema defines the schema for the resource.
func (r *trustStoreResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"serial_number": schema.StringAttribute{
				Computed: true,
			},
			"certificate": schema.StringAttribute{
				Required: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"issuer": schema.StringAttribute{
				Computed: true,
			},
			"signature": schema.StringAttribute{
				Computed: true,
			},
			"uploaded_on": schema.StringAttribute{
				Computed: true,
			},
			"uploaded_at": schema.StringAttribute{
				Computed: true,
			},
			"expires_on": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create a new resource.
func (r *trustStoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan trustStoreModel
	jsonModel := &trustStoreJsonModel{}
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	values := map[string]string{"certificate": plan.Certificate.ValueString()}
	jsonData, _ := json.Marshal(values)

	request, err := http.NewRequest("POST", r.client.Url, bytes.NewBuffer(jsonData))
	request.Header.Add("Authorization", "Bearer "+r.client.Token)
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		resp.Diagnostics.AddError("Cannot send post request", err.Error())
	}

	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	if response.StatusCode != 201 {
		resp.Diagnostics.AddError("Not created", bodyString)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	_ = json.Unmarshal(bodyBytes, &jsonModel)

	plan.ID = types.StringValue(jsonModel.Result.ID)
	plan.SerialNumber = types.StringValue(jsonModel.Result.SerialNumber)
	plan.Status = types.StringValue(jsonModel.Result.Status)
	plan.Issuer = types.StringValue(jsonModel.Result.Issuer)
	plan.Signature = types.StringValue(jsonModel.Result.Signature)
	plan.UploadedOn = types.StringValue(jsonModel.Result.UploadedOn)
	plan.UploadedAt = types.StringValue(jsonModel.Result.UploadedAt)
	plan.ExpiresOn = types.StringValue(jsonModel.Result.ExpiresOn)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *trustStoreResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *trustStoreResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *trustStoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state trustStoreModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	request, err := http.NewRequest("DELETE", r.client.Url+"/"+state.ID.ValueString(), nil)
	request.Header.Add("Authorization", "Bearer "+r.client.Token)
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		resp.Diagnostics.AddError("Cannot send delete request", err.Error())
	}

	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	if response.StatusCode != 200 {
		resp.Diagnostics.AddError("Not deleted", bodyString)
	}

	if resp.Diagnostics.HasError() {
		return
	}
}
