package solidserver

import (
  "encoding/json"
  "net/url"
  "strconv"
  "strings"
  "fmt"
  "log"
)

// Integer Absolute value
func abs(x int) int {
  if (x < 0) {
    return -x
  }

  return x
}

// Convert hexa IP address string into standard IP address string
// Return an empty string in case of failure
func hexiptoip(hexip string) string {
  a, b, c, d := 0,0,0,0

  count, _ := fmt.Sscanf(hexip, "%02x%02x%02x%02x", &a, &b, &c, &d)

  if (count == 4) {
    return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
  }

  return ""
}

// Convert standard IP address string into hexa IP address string
// Return an empty string in case of failure
func iptohexip(ip string) string {
  ip_dec := strings.Split(ip, ".")

  if (len(ip_dec) == 4) {

    a, _ := strconv.Atoi(ip_dec[0])
    b, _ := strconv.Atoi(ip_dec[1])
    c, _ := strconv.Atoi(ip_dec[2])
    d, _ := strconv.Atoi(ip_dec[3])

    if (0 <= a && a <= 255 && 0 <= b && b <= 255 &&
        0 <= c && c <= 255 && 0 <= d && d <= 255) {
      return fmt.Sprintf("%02x%02x%02x%02x", a, b, c, d)
    }

    return ""
  }

  return ""
}

// Convert standard IP address string into unsigned int32
// Return 0 in case of failure
func iptolong(ip string) uint32 {
  ip_dec := strings.Split(ip, ".")

  if (len(ip_dec) == 4) {
    a, _ := strconv.Atoi(ip_dec[0])
    b, _ := strconv.Atoi(ip_dec[1])
    c, _ := strconv.Atoi(ip_dec[2])
    d, _ := strconv.Atoi(ip_dec[3])

    var iplong uint32 = uint32(a) * 0x1000000
    iplong += uint32(b) * 0x10000
    iplong += uint32(c) * 0x100
    iplong += uint32(d) * 0x1

    return iplong
  }

  return 0
}

// Convert unsigned int32 into standard IP address string
// Return an IP formated string
func longtoip(iplong uint32) string {
  a := (iplong & 0xFF000000) >> 24;
  b := (iplong & 0xFF0000) >> 16;
  c := (iplong & 0xFF00) >> 8;
  d := (iplong & 0xFF);

  if (a < 0) {
    a = a + 0x100;
  }

  return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
}

// Compute the actual size of a CIDR prefix from its length
// Return -1 in case of failure
func prefixlengthtosize(length int) int {
  if (length <= 32) {
    return (1 << (32 - uint32(length)));
  }

  return -1
}

// Build url value object from class parameters
// Return an url.Values{} object
func urlfromclassparams(parameters interface{}) url.Values {
  class_parameters := url.Values{}

  for k, v := range parameters.(map[string]interface{}) {
    class_parameters.Add(k, v.(string))
  }

  return class_parameters
}

// Return the oid of a device from hostdev_name
// Or an empty string in case of failure
func hostdevidbyname(hostdev_name string, meta interface{}) (string, error) {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("WHERE", "hostdev_name='" + strings.ToLower(hostdev_name) + "'")


  // Sending the read request
  http_resp, body, err := s.Request("get", "rest/hostdev_list", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if (http_resp.StatusCode == 200 && len(buf) > 0) {
      if hostdev_id, hostdev_id_exist := buf[0]["hostdev_id"].(string); (hostdev_id_exist) {
        return hostdev_id, nil
      }
    }
  }

  log.Printf("[DEBUG] SOLIDServer - Unable to find device: %s", hostdev_name)

  return "", err
}

// Return an available IP addresses from site_id, block_id and expected subnet_size
// Or an empty table of string in case of failure
func ipaddressfindfree(subnet_id string, meta interface{}) ([]string, error) {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("subnet_id", subnet_id)
  parameters.Add("max_find", "4")

  // Sending the creation request
  http_resp, body, err := s.Request("get", "rpc/ip_find_free_address", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if (http_resp.StatusCode == 200 && len(buf) > 0) {
      addresses := []string{}

      for i := 0; i < len(buf); i++ {
        if addr, addr_exist := buf[i]["hostaddr"].(string); (addr_exist) {
          log.Printf("[DEBUG] SOLIDServer - Suggested IP address: %s", addr)
          addresses = append(addresses, addr)
        }
      }
      return addresses, nil
    }
  }

  log.Printf("[DEBUG] SOLIDServer - Unable to find a free IP address in subnet (oid): %s", subnet_id)

  return []string{}, err
}

// Return an available vlan from specified vlmdomain_name
// Or an empty table strings in case of failure
func vlanidfindfree(vlmdomain_name string, meta interface{}) ([]string, error) {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("limit", "4")
  //https://192.168.1.2:443/rest/vlmvlan_list?WHERE=vlmdomain_name%3D'local' and row_enabled%3D'2'&limit=2
  parameters.Add("WHERE", "vlmdomain_name='" + strings.ToLower(vlmdomain_name) + "' AND row_enabled='2'")

  // Sending the creation request
  http_resp, body, err := s.Request("get", "rest/vlmvlan_list", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if (http_resp.StatusCode == 200 && len(buf) > 0) {
      vnids := []string{}

      for i := 0; i < len(buf); i++ {
        if vnid, vnid_exist := buf[i]["vlmvlan_vlan_id"].(string); (vnid_exist) {
          log.Printf("[DEBUG] SOLIDServer - Suggested vlan ID: %s", vnid)
          vnids = append(vnids, vnid)
        }
      }
      return vnids, nil
    }
  }

  log.Printf("[DEBUG] SOLIDServer - Unable to find a free vlan ID in vlan domain: %s", vlmdomain_name)

  return []string{}, err
}

// Return the oid of a space from site_name
// Or an empty string in case of failure
func ipsiteidbyname(site_name string, meta interface{}) (string, error) {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("WHERE", "site_name='" + strings.ToLower(site_name) + "'")


  // Sending the read request
  http_resp, body, err := s.Request("get", "rest/ip_site_list", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if (http_resp.StatusCode == 200 && len(buf) > 0) {
      if site_id, site_id_exist := buf[0]["site_id"].(string); (site_id_exist) {
        return site_id, nil
      }
    }
  }

  log.Printf("[DEBUG] SOLIDServer - Unable to find IP space: %s", site_name)

  return "", err
}

// Return the oid of a vlan domain from vlmdomain_name
// Or an empty string in case of failure
func vlandomainidbyname(vlmdomain_name string, meta interface{}) (string, error) {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("WHERE", "vlmdomain_name='" + strings.ToLower(vlmdomain_name) + "'")


  // Sending the read request
  http_resp, body, err := s.Request("get", "rest/vlmdomain_name", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if (http_resp.StatusCode == 200 && len(buf) > 0) {
      if vlmdomain_id, vlmdomain_id_exist := buf[0]["vlmdomain_id"].(string); (vlmdomain_id_exist) {
        return vlmdomain_id, nil
      }
    }
  }

  log.Printf("[DEBUG] SOLIDServer - Unable to find vlan domain: %s", vlmdomain_name)

  return "", err
}

// Return the oid of a subnet from site_id, subnet_name and is_terminal property
// Or an empty string in case of failure
func ipsubnetidbyname(site_id string, subnet_name string, terminal bool, meta interface{}) (string, error) {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("WHERE", "site_id='" + site_id + "' AND " + "subnet_name='" + strings.ToLower(subnet_name) + "'")
  if (terminal) {
    parameters.Add("is_terminal", "1")
  } else {
    parameters.Add("is_terminal", "0")
  }

  // Sending the read request
  http_resp, body, err := s.Request("get", "rest/ip_block_subnet_list", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if (http_resp.StatusCode == 200 && len(buf) > 0) {
      if subnet_id, subnet_id_exist := buf[0]["subnet_id"].(string); (subnet_id_exist) {
        return subnet_id, nil
      }
    }
  }

  log.Printf("[DEBUG] SOLIDServer - Unable to find IP subnet: %s", subnet_name)

  return "", err
}

// Return the oid of an address from site_id, ip_address
// Or an empty string in case of failure
func ipaddressidbyip(site_id string, ip_address string, meta interface{}) (string, error) {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("WHERE", "site_id='" + site_id + "' AND " + "ip_addr='" + iptohexip(ip_address) + "'")

  // Sending the read request
  http_resp, body, err := s.Request("get", "rest/ip_address_list", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if (http_resp.StatusCode == 200 && len(buf) > 0) {
      if ip_id, ip_id_exist := buf[0]["ip_id"].(string); (ip_id_exist) {
        return ip_id, nil
      }
    }
  }

  log.Printf("[DEBUG] SOLIDServer - Unable to find IP address: %s", ip_address)

  return "", err
}

// Return the oid of an address from ip_id, ip_name_type, alias_name
// Or an empty string in case of failure
func ipaliasidbyinfo(address_id string, alias_name string, ip_name_type string, meta interface{}) (string, error) {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("ip_id", address_id)
  // Bug - Ticket 18653
  // parameters.Add("WHERE", "ip_name_type='" + ip_name_type + "' AND " + "alias_name='" + alias_name + "'")

  // Sending the read request
  http_resp, body, err := s.Request("get", "rest/ip_alias_list", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Shall be removed once Ticket 18653 is closed
    // Checking the answer
    if (http_resp.StatusCode == 200 && len(buf) > 0) {
      for i := 0; i < len(buf); i++ {
        r_ip_name_id, r_ip_name_id_exist := buf[i]["ip_name_id"].(string)
        r_ip_name_type, r_ip_name_type_exist := buf[i]["ip_name_type"].(string)
        r_alias_name, r_alias_name_exist := buf[i]["alias_name"].(string)

        log.Printf("[DEBUG] SOLIDServer - Comparing '%s' with '%s' looking for IP alias associated with IP address ID %s", alias_name, r_alias_name, address_id)
        log.Printf("[DEBUG] SOLIDServer - Comparing '%s' with '%s' looking for IP alias associated with IP address ID %s", ip_name_type, r_ip_name_type, address_id)

        if (r_ip_name_type_exist && strings.Compare(ip_name_type, r_ip_name_type) == 0 &&
            r_alias_name_exist   && strings.Compare(alias_name, r_alias_name) == 0 &&
            r_ip_name_id_exist) {
          return r_ip_name_id, nil
        }
      }
    }
  }

  // Shall be restored once Ticket 18653 is closed
  // Checking the answer
  //if (http_resp.StatusCode == 200 && len(buf) > 0) {
  //  if ip_name_id, ip_name_id_exist := buf[0]["ip_name_id"].(string); (ip_name_id_exist) {
  //    return ip_name_id
  //  }
  //}

  log.Printf("[DEBUG] SOLIDServer - Unable to find IP alias: %s - %s associated with IP address ID %s", alias_name, ip_name_type, address_id)

  return "", err
}

// Return an available subnet address from site_id, block_id and expected subnet_size
// Or an empty string in case of failure
func ipsubnetfindbysize(site_id string, block_id string, prefix_size int, meta interface{}) ([]string, error) {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("site_id", site_id)
  parameters.Add("block_id", block_id)
  parameters.Add("prefix", strconv.Itoa(prefix_size))
  parameters.Add("max_find", "4")

  // Sending the creation request
  http_resp, body, err := s.Request("get", "rpc/ip_find_free_subnet", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if (http_resp.StatusCode == 200 && len(buf) > 0) {
      subnet_addresses := []string{}

      for i := 0; i < len(buf); i++ {
        if hexaddr, hexaddr_exist := buf[i]["start_ip_addr"].(string); (hexaddr_exist) {
          log.Printf("[DEBUG] SOLIDServer - Suggested subnet address: %s", hexiptoip(hexaddr))
          subnet_addresses = append(subnet_addresses, hexaddr)
        }
      }
      return subnet_addresses, nil
    }
  }

  log.Printf("[DEBUG] SOLIDServer - Unable to find a free IP subnet in space (oid): %s, block (oid): %s, size: ", site_id, block_id, strconv.Itoa(prefix_size))

  return []string{}, err
}