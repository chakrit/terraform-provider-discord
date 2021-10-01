package discord

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/net/context"
)

func resourceDiscordInvite() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInviteCreate,
		ReadContext:   resourceInviteRead,
		DeleteContext: resourceInviteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"channel_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"max_age": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
				Default:  86400,
			},
			"max_uses": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
			},
			"temporary": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
			},
			"unique": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
			},
			"code": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceInviteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*Context).Client

	channelId := getId(d.Get("channel_id").(string))

	builder := client.Channel(channelId).CreateInvite()
	builder.SetMaxAge(d.Get("max_age").(int))
	builder.SetMaxUses(d.Get("max_uses").(int))
	builder.SetTemporary(d.Get("temporary").(bool))
	builder.SetUnique(d.Get("unique").(bool))

	invite, err := builder.Execute()
	if err != nil {
		return diag.Errorf("Failed to create a invite: %s", err.Error())
	}

	d.SetId(invite.Code)
	d.Set("code", invite.Code)

	return diags
}

func resourceInviteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*Context).Client

	invite, err := client.Invite(d.Id()).Get(false)
	if err != nil {
		d.SetId("")
	} else {
		d.Set("code", invite.Code)
	}

	return diags
}

func resourceInviteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*Context).Client

	_, err := client.Invite(d.Id()).Delete()
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
