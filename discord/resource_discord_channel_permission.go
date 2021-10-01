package discord

import (
	"fmt"
	"strconv"

	"github.com/andersfylling/disgord"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/net/context"
)

func resourceDiscordChannelPermission() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChannelPermissionCreate,
		ReadContext:   resourceChannelPermissionRead,
		UpdateContext: resourceChannelPermissionUpdate,
		DeleteContext: resourceChannelPermissionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"channel_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateDiagFunc: func(val interface{}, path cty.Path) (diags diag.Diagnostics) {
					v := val.(string)

					if v != "role" && v != "user" {
						diags = append(diags, diag.Errorf("%s is not a valid type. Must be \"role\" or \"user\"", v)...)
					}

					return diags
				},
			},
			"overwrite_id": {
				ForceNew: true,
				Required: true,
				Type:     schema.TypeString,
			},
			"allow": {
				AtLeastOneOf: []string{"allow", "deny"},
				Optional:     true,
				Type:         schema.TypeInt,
			},
			"deny": {
				AtLeastOneOf: []string{"allow", "deny"},
				Optional:     true,
				Type:         schema.TypeInt,
			},
		},
	}
}

func resourceChannelPermissionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*Context).Client

	channelId := getId(d.Get("channel_id").(string))
	overwriteId := getId(d.Get("overwrite_id").(string))

	err := client.Channel(channelId).UpdatePermissions(overwriteId, &disgord.UpdateChannelPermissionsParams{
		Allow: disgord.PermissionBit(d.Get("allow").(uint64)),
		Deny:  disgord.PermissionBit(d.Get("deny").(uint64)),
		Type:  d.Get("type").(uint),
	})

	if err != nil {
		return diag.Errorf("Failed to update channel permissions %s: %s", channelId.String(), err.Error())
	}

	d.SetId(strconv.Itoa(Hashcode(fmt.Sprintf("%s:%s:%s", channelId, overwriteId, d.Get("type").(string)))))

	return diags
}

func resourceChannelPermissionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*Context).Client

	channelId := getId(d.Get("channel_id").(string))
	overwriteId := getId(d.Get("overwrite_id").(string))

	channel, err := client.Channel(channelId).Get()
	if err != nil {
		return diag.Errorf("Failed to find channel %s: %s", channelId.String(), err.Error())
	}

	for _, x := range channel.PermissionOverwrites {
		if x.Type == disgord.PermissionOverwriteType(d.Get("type").(uint8)) && x.ID == overwriteId {
			d.Set("allow", uint64(x.Allow))
			d.Set("deny", uint64(x.Deny))
			break
		}
	}

	return diags
}

func resourceChannelPermissionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*Context).Client

	channelId := getId(d.Get("channel_id").(string))
	overwriteId := getId(d.Get("overwrite_id").(string))

	err := client.Channel(channelId).UpdatePermissions(overwriteId, &disgord.UpdateChannelPermissionsParams{
		Allow: disgord.PermissionBit(d.Get("allow").(uint64)),
		Deny:  disgord.PermissionBit(d.Get("deny").(uint64)),
		Type:  d.Get("type").(uint),
	})

	if err != nil {
		return diag.Errorf("Failed to update channel permissions %s: %s", channelId.String(), err.Error())
	}

	return diags
}

func resourceChannelPermissionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*Context).Client

	channelId := getId(d.Get("channel_id").(string))
	overwriteId := getId(d.Get("overwrite_id").(string))
	err := client.Channel(channelId).DeletePermission(overwriteId)

	if err != nil {
		return diag.Errorf("Failed to delete channel permissions %s: %s", channelId.String(), err.Error())
	}

	return diags
}
