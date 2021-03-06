package solidserver

import (
  "github.com/hashicorp/terraform/terraform"
  "github.com/hashicorp/terraform/helper/schema"
)

func Provider() terraform.ResourceProvider {
  return &schema.Provider{
    Schema: map[string]*schema.Schema{
      "username": &schema.Schema{
        Type:        schema.TypeString,
        Required:    true,
        DefaultFunc: schema.EnvDefaultFunc("SOLIDServer_USERNAME", nil),
        Description: "SOLIDServer API user's ID",
      },
      "password": &schema.Schema{
        Type:        schema.TypeString,
        Required:    true,
        DefaultFunc: schema.EnvDefaultFunc("SOLIDServer_PASSWORD", nil),
        Description: "SOLIDServer API user's password",
      },
      "host": &schema.Schema{
        Type:        schema.TypeString,
        Required:    true,
        DefaultFunc: schema.EnvDefaultFunc("SOLIDServer_HOST", nil),
        Description: "SOLIDServer API hostname or IP address",
      },
      "sslverify": &schema.Schema{
        Type:        schema.TypeBool,
        Required:    false,
        Optional:    true,
        DefaultFunc: schema.EnvDefaultFunc("SOLIDServer_SSLVERIFY", true),
        Description: "Enable/Disable ssl verify (Default : enabled)",
      },
      "additional_trust_certs_file": &schema.Schema{
        Type:        schema.TypeString,
        Required:    false,
        Optional:    true,
        DefaultFunc: schema.EnvDefaultFunc("SOLIDServer_ADDITIONALTRUSTCERTSFILE", nil),
        Description: "PEM formatted file with additional certificates to trust for TLS connection",
      },
    },

    ResourcesMap: map[string]*schema.Resource{
      "solidserver_ip_subnet": resourceipsubnet(),
      "solidserver_ip_address": resourceipaddress(),
      "solidserver_ip_alias": resourceipalias(),
      "solidserver_device": resourcedevice(),
      "solidserver_vlan": resourcevlan(),
      "solidserver_dns_zone": resourcednszone(),
      "solidserver_dns_rr": resourcednsrr(),
    },

    ConfigureFunc: ProviderConfigure,
  }
}

func ProviderConfigure(d *schema.ResourceData) (interface{}, error) {
  s := NewSOLIDserver(
    d.Get("host").(string),
    d.Get("username").(string),
    d.Get("password").(string),
    d.Get("sslverify").(bool),
    d.Get("additional_trust_certs_file").(string),
  )

  return s, nil
}
