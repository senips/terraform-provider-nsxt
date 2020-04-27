/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/go-vmware-nsxt/manager"
)

func dataSourceLogicalSwitchInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLogicalSwitchInfoRead,

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:        schema.TypeString,
				Description: "The display name of this resource",
				Optional:    true,
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceLogicalSwitchInfoRead(d *schema.ResourceData, m interface{}) error {
	// Read a transport zone by name or id
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return dataSourceNotSupportedError()
	}

	objName := d.Get("display_name").(string)
	var obj manager.LogicalSwitch
	if objName == "" {
		return fmt.Errorf("Error obtaining logical switch name during read")
	} else {
		// Get by full name
		// TODO use localVarOptionals for paging
		localVarOptionals := make(map[string]interface{})

		objList, _, err := nsxClient.LogicalSwitchingApi.ListLogicalSwitches(nsxClient.Context, localVarOptionals)
		if err != nil {
			return fmt.Errorf("Error while reading logical switches: %v", err)
		}
		// go over the list to find the correct one
		found := false
		for _, objInList := range objList.Results {
			if objInList.DisplayName == objName {
				if found {
					return fmt.Errorf("Found multiple logical switches with name '%s'", objName)
				}
				obj = objInList
				found = true
			}
		}
		if !found {
			return fmt.Errorf("Logical Switch with name '%s' was not found", objName)
		}
	}

	d.SetId(obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)

	return nil
}
