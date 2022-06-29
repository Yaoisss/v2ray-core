//go:build !confonly
// +build !confonly

package dns

import (
	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/common/strmatcher"
	"github.com/v2fly/v2ray-core/v5/common/uuid"
)

/*查询域名方式
1.Full		0
2.Domain	1
3.Substr	2
4.Regex		3
*/
var typeMap = map[DomainMatchingType]strmatcher.Type{
	DomainMatchingType_Full:      strmatcher.Full,
	DomainMatchingType_Subdomain: strmatcher.Domain,
	DomainMatchingType_Keyword:   strmatcher.Substr,
	DomainMatchingType_Regex:     strmatcher.Regex,
}

// References:
// https://www.iana.org/assignments/special-use-domain-names/special-use-domain-names.xhtml
// https://unix.stackexchange.com/questions/92441/whats-the-difference-between-local-home-and-lan
/*
	使用grpc的数据结构，读取本地特殊的地址信息（特殊地址的备忘录[RFC6761]）

*/
var localTLDsAndDotlessDomains = []*NameServer_PriorityDomain{
	{Type: DomainMatchingType_Regex, Domain: "^[^.]+$"}, // This will only match domains without any dot
	{Type: DomainMatchingType_Subdomain, Domain: "local"},
	{Type: DomainMatchingType_Subdomain, Domain: "localdomain"},
	{Type: DomainMatchingType_Subdomain, Domain: "localhost"},
	{Type: DomainMatchingType_Subdomain, Domain: "lan"},
	{Type: DomainMatchingType_Subdomain, Domain: "home.arpa"},
	{Type: DomainMatchingType_Subdomain, Domain: "example"},
	{Type: DomainMatchingType_Subdomain, Domain: "invalid"},
	{Type: DomainMatchingType_Subdomain, Domain: "test"},
}

/*
	对本地TLD和无点域规则的匹配，后面得到localTLDsAndDotlessDomains的大小
*/
var localTLDsAndDotlessDomainsRule = &NameServer_OriginalRule{
	Rule: "geosite:private",
	Size: uint32(len(localTLDsAndDotlessDomains)),
}

/*
	字符串匹配器,第一判断类型是不是0、1、2、3，第二判断-确定字符串与模式匹配的接口New->Matcher
*/
func toStrMatcher(t DomainMatchingType, domain string) (strmatcher.Matcher, error) {
	strMType, f := typeMap[t]
	if !f {
		return nil, newError("unknown mapping type", t).AtWarning()
	}
	matcher, err := strMType.New(domain)
	if err != nil {
		return nil, newError("failed to create str matcher").Base(err)
	}
	return matcher, nil
}

/*

 */
func toNetIP(addrs []net.Address) ([]net.IP, error) {
	ips := make([]net.IP, 0, len(addrs))
	for _, addr := range addrs {
		if addr.Family().IsIP() {
			ips = append(ips, addr.IP())
		} else {
			return nil, newError("Failed to convert address", addr, "to Net IP.").AtWarning()
		}
	}
	return ips, nil
}

/*
generate Random Tag
生成随机标记
*/
func generateRandomTag() string {
	id := uuid.New()
	return "v2ray.system." + id.String()
}
