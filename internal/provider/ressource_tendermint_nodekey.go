package provider

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

type PrivKeyJSON struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type NodeKeyJSON struct {
	PrivKey PrivKeyJSON `json:"priv_key"`
}

func resourceTendermintNodeKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTendermintNodeKeyCreate,
		ReadContext:   resourceTendermintNodeKeyRead,
		DeleteContext: resourceTendermintNodeKeyDelete,

		Schema: map[string]*schema.Schema{
			"node_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tendermint node ID.",
			},
			"node_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Tendermint node key.",
			},
			"node_key_json": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Tendermint node key in the JSON format expected by a node.",
			},
		},
	}
}

func resourceTendermintNodeKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	privKey := ed25519.GenPrivKey()
	pubKey := privKey.PubKey()

	nodeId := strings.ToLower(pubKey.Address().String())
	privKeyBase64 := base64.StdEncoding.EncodeToString(privKey.Bytes())

	nodeKeyJSON := NodeKeyJSON{
		PrivKey: PrivKeyJSON{
			Type:  "tendermint/PrivKeyEd25519",
			Value: privKeyBase64,
		},
	}

	nodeKeyJSONString, err := json.Marshal(nodeKeyJSON)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("node_id", nodeId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("node_key", privKeyBase64); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("node_key_json", string(nodeKeyJSONString)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(nodeId)

	return diags
}

func resourceTendermintNodeKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceTendermintNodeKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
