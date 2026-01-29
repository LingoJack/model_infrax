package tool

import (
	"net"
	"net/http"
	"strings"
)

// IPOption 定义获取客户端IP的选项（使用位运算）
//
// 通过位运算可以灵活组合多个选项，例如：
//   - IPOptionTrustProxy | IPOptionAllowPrivate
//   - IPOptionStrict | IPOptionFirstOnly
//
// 使用示例:
//
//	opts := IPOptionTrustProxy | IPOptionAllowPrivate
//	ip := GetClientIP(r, opts)
type IPOption uint32

const (
	// IPOptionNone 默认选项，不做任何特殊配置
	// 默认会尝试从所有常见的代理头部获取IP
	IPOptionNone IPOption = 0

	// IPOptionTrustProxy 信任代理头部（X-Forwarded-For, X-Real-IP等）
	// 适用于应用部署在反向代理（如Nginx、负载均衡器）后面的场景
	// 如果不设置此选项，将只使用 RemoteAddr
	IPOptionTrustProxy IPOption = 1 << 0 // 1

	// IPOptionFirstOnly 仅获取X-Forwarded-For中的第一个IP
	// 适用于只需要原始客户端IP的场景
	// 如果不设置，会返回最后一个非私有IP（更可靠）
	IPOptionFirstOnly IPOption = 1 << 1 // 2

	// IPOptionAllowPrivate 允许返回私有IP地址
	// 适用于内网环境或开发测试场景
	// 如果不设置，会跳过私有IP地址
	IPOptionAllowPrivate IPOption = 1 << 2 // 4

	// IPOptionStrict 严格模式，只从X-Real-IP获取
	// 适用于已知代理服务器会设置X-Real-IP的场景
	// 通常用于Nginx配置了 proxy_set_header X-Real-IP $remote_addr 的情况
	IPOptionStrict IPOption = 1 << 3 // 8

	// IPOptionIncludePort 返回结果包含端口号
	// 适用于需要记录完整连接信息的场景
	// 格式为 IP:Port
	IPOptionIncludePort IPOption = 1 << 4 // 16
)

// 常见的代理头部名称
const (
	// HeaderXForwardedFor 标准的代理头部，格式为 "client, proxy1, proxy2"
	HeaderXForwardedFor = "X-Forwarded-For"

	// HeaderXRealIP Nginx常用的头部，通常包含原始客户端IP
	HeaderXRealIP = "X-Real-IP"

	// HeaderXClientIP 某些代理使用的头部
	HeaderXClientIP = "X-Client-IP"

	// HeaderCFConnectingIP Cloudflare使用的头部
	HeaderCFConnectingIP = "CF-Connecting-IP"

	// HeaderTrueClientIP Akamai和Cloudflare使用的头部
	HeaderTrueClientIP = "True-Client-IP"

	// HeaderXForwardedHost 转发的原始Host
	HeaderXForwardedHost = "X-Forwarded-Host"
)

// Has 检查是否包含指定选项
//
// 使用示例:
//
//	opts := IPOptionTrustProxy | IPOptionAllowPrivate
//	if opts.Has(IPOptionTrustProxy) {
//	    fmt.Println("信任代理头部")
//	}
func (o IPOption) Has(option IPOption) bool {
	return o&option != 0
}

// GetClientIP 从HTTP请求中获取客户端IP地址
//
// 该函数按照以下优先级尝试获取客户端的真实IP地址：
//  1. X-Real-IP（严格模式下优先）
//  2. CF-Connecting-IP（Cloudflare）
//  3. True-Client-IP（Akamai/Cloudflare）
//  4. X-Forwarded-For（标准代理头部）
//  5. X-Client-IP
//  6. RemoteAddr（直接连接地址）
//
// 参数:
//   - r: HTTP请求对象（*http.Request）
//   - options: 配置选项，可以通过位运算组合多个选项（可选参数，默认为IPOptionNone）
//
// 返回值:
//   - ip: 客户端IP地址字符串，如果无法获取则返回空字符串
//
// 使用示例:
//
//	// 基本使用（仅使用RemoteAddr，最安全）
//	ip := GetClientIP(r)
//
//	// 信任代理头部（应用在反向代理后面时使用）
//	ip := GetClientIP(r, IPOptionTrustProxy)
//
//	// 信任代理 + 允许私有IP（内网/开发环境）
//	ip := GetClientIP(r, IPOptionTrustProxy|IPOptionAllowPrivate)
//
//	// 严格模式，只从X-Real-IP获取
//	ip := GetClientIP(r, IPOptionTrustProxy|IPOptionStrict)
//
//	// 获取包含端口的完整地址
//	ip := GetClientIP(r, IPOptionIncludePort)
//
// 输出示例:
//
//	"192.168.1.100"      // 普通IP
//	"2001:db8::1"        // IPv6地址
//	"192.168.1.100:8080" // 带端口（使用IPOptionIncludePort时）
//
// 安全注意事项:
//   - 如果应用直接面向公网，不要使用IPOptionTrustProxy，因为客户端可以伪造这些头部
//   - 只有在确认有可信的反向代理时才使用IPOptionTrustProxy
//   - 建议在Nginx等代理中配置覆盖X-Forwarded-For，防止伪造
//
// Nginx配置示例:
//
//	proxy_set_header X-Real-IP $remote_addr;
//	proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
func GetClientIP(r *http.Request, options ...IPOption) (ip string) {
	if r == nil {
		return ""
	}

	// 合并所有选项
	var opt IPOption = IPOptionNone
	for _, o := range options {
		opt |= o
	}

	// 如果不信任代理头部，直接返回RemoteAddr
	if !opt.Has(IPOptionTrustProxy) {
		return extractIP(r.RemoteAddr, opt)
	}

	// 严格模式：只从X-Real-IP获取
	if opt.Has(IPOptionStrict) {
		if realIP := r.Header.Get(HeaderXRealIP); realIP != "" {
			if validIP := validateAndExtractIP(realIP, opt); validIP != "" {
				return validIP
			}
		}
		// 严格模式下，如果X-Real-IP无效，回退到RemoteAddr
		return extractIP(r.RemoteAddr, opt)
	}

	// 按优先级尝试获取IP

	// 1. X-Real-IP（通常由Nginx设置的可信IP）
	if realIP := r.Header.Get(HeaderXRealIP); realIP != "" {
		if validIP := validateAndExtractIP(realIP, opt); validIP != "" {
			return validIP
		}
	}

	// 2. CF-Connecting-IP（Cloudflare）
	if cfIP := r.Header.Get(HeaderCFConnectingIP); cfIP != "" {
		if validIP := validateAndExtractIP(cfIP, opt); validIP != "" {
			return validIP
		}
	}

	// 3. True-Client-IP（Akamai/Cloudflare）
	if trueIP := r.Header.Get(HeaderTrueClientIP); trueIP != "" {
		if validIP := validateAndExtractIP(trueIP, opt); validIP != "" {
			return validIP
		}
	}

	// 4. X-Forwarded-For（可能包含多个IP）
	if xff := r.Header.Get(HeaderXForwardedFor); xff != "" {
		if validIP := parseXForwardedFor(xff, opt); validIP != "" {
			return validIP
		}
	}

	// 5. X-Client-IP
	if clientIP := r.Header.Get(HeaderXClientIP); clientIP != "" {
		if validIP := validateAndExtractIP(clientIP, opt); validIP != "" {
			return validIP
		}
	}

	// 6. 最后使用RemoteAddr
	return extractIP(r.RemoteAddr, opt)
}

// GetClientIPSimple 简化版获取客户端IP
//
// 这是GetClientIP的简化版本，默认信任代理头部。
// 适用于大多数部署在反向代理后面的应用场景。
//
// 参数:
//   - r: HTTP请求对象
//
// 返回值:
//   - ip: 客户端IP地址
//
// 使用示例:
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//	    clientIP := GetClientIPSimple(r)
//	    log.Printf("客户端IP: %s", clientIP)
//	}
//
// 注意事项:
//   - 此函数默认信任代理头部，不适用于直接面向公网的应用
//   - 如需更严格的控制，请使用 GetClientIP 并指定选项
func GetClientIPSimple(r *http.Request) string {
	return GetClientIP(r, IPOptionTrustProxy)
}

// MustGetClientIP 获取客户端IP，如果无法获取则返回默认值
//
// 参数:
//   - r: HTTP请求对象
//   - defaultIP: 无法获取时返回的默认IP
//   - options: 配置选项（可选）
//
// 返回值:
//   - ip: 客户端IP地址，如果无法获取则返回defaultIP
//
// 使用示例:
//
//	// 无法获取IP时返回 "unknown"
//	ip := MustGetClientIP(r, "unknown", IPOptionTrustProxy)
//
//	// 无法获取IP时返回 "0.0.0.0"
//	ip := MustGetClientIP(r, "0.0.0.0")
func MustGetClientIP(r *http.Request, defaultIP string, options ...IPOption) string {
	ip := GetClientIP(r, options...)
	if ip == "" {
		return defaultIP
	}
	return ip
}

// parseXForwardedFor 解析X-Forwarded-For头部
//
// X-Forwarded-For格式: "client, proxy1, proxy2"
// 第一个IP是原始客户端IP，后续是经过的代理服务器IP
//
// 参数:
//   - xff: X-Forwarded-For头部值
//   - opt: 选项配置
//
// 返回值:
//   - ip: 有效的IP地址
func parseXForwardedFor(xff string, opt IPOption) string {
	// 按逗号分割IP列表
	ips := strings.Split(xff, ",")
	if len(ips) == 0 {
		return ""
	}

	// 如果只需要第一个IP
	if opt.Has(IPOptionFirstOnly) {
		firstIP := strings.TrimSpace(ips[0])
		return validateAndExtractIP(firstIP, opt)
	}

	// 从后向前遍历，找到第一个有效的非私有IP
	// 这种方式可以防止客户端伪造X-Forwarded-For头部
	for i := len(ips) - 1; i >= 0; i-- {
		ip := strings.TrimSpace(ips[i])
		if validIP := validateAndExtractIP(ip, opt); validIP != "" {
			// 如果允许私有IP，直接返回
			if opt.Has(IPOptionAllowPrivate) {
				return validIP
			}
			// 否则检查是否为私有IP
			if !isPrivateIP(validIP) {
				return validIP
			}
		}
	}

	// 如果所有IP都是私有的，且允许私有IP，返回第一个有效IP
	if opt.Has(IPOptionAllowPrivate) {
		firstIP := strings.TrimSpace(ips[0])
		return validateAndExtractIP(firstIP, opt)
	}

	return ""
}

// validateAndExtractIP 验证并提取IP地址
//
// 参数:
//   - ipStr: IP字符串（可能包含端口）
//   - opt: 选项配置
//
// 返回值:
//   - ip: 有效的IP地址，无效则返回空字符串
func validateAndExtractIP(ipStr string, opt IPOption) string {
	ipStr = strings.TrimSpace(ipStr)
	if ipStr == "" {
		return ""
	}

	// 尝试解析为IP
	ip := net.ParseIP(ipStr)
	if ip != nil {
		// 检查私有IP
		if !opt.Has(IPOptionAllowPrivate) && isPrivateIP(ipStr) {
			return ""
		}
		return ipStr
	}

	// 可能包含端口，尝试分离
	host, _, err := net.SplitHostPort(ipStr)
	if err != nil {
		return ""
	}

	ip = net.ParseIP(host)
	if ip == nil {
		return ""
	}

	// 检查私有IP
	if !opt.Has(IPOptionAllowPrivate) && isPrivateIP(host) {
		return ""
	}

	return host
}

// extractIP 从地址字符串中提取IP
//
// 处理可能带端口的地址格式（如 "192.168.1.1:8080" 或 "[::1]:8080"）
//
// 参数:
//   - addr: 地址字符串
//   - opt: 选项配置
//
// 返回值:
//   - ip: IP地址字符串
func extractIP(addr string, opt IPOption) string {
	if addr == "" {
		return ""
	}

	// 如果需要包含端口，直接返回
	if opt.Has(IPOptionIncludePort) {
		return addr
	}

	// 尝试分离主机和端口
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		// 可能是纯IP地址（没有端口）
		if net.ParseIP(addr) != nil {
			return addr
		}
		return ""
	}

	return host
}

// isPrivateIP 检查IP是否为私有地址
//
// 私有IP地址范围：
//   - 10.0.0.0/8
//   - 172.16.0.0/12
//   - 192.168.0.0/16
//   - fc00::/7 (IPv6)
//   - 127.0.0.0/8 (localhost)
//   - ::1 (IPv6 localhost)
//
// 参数:
//   - ipStr: IP地址字符串
//
// 返回值:
//   - bool: 是否为私有IP
//
// 使用示例:
//
//	if isPrivateIP("192.168.1.1") {
//	    fmt.Println("这是私有IP")
//	}
func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// 检查是否为回环地址
	if ip.IsLoopback() {
		return true
	}

	// 检查是否为私有地址
	if ip.IsPrivate() {
		return true
	}

	// 检查是否为链路本地地址
	if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	return false
}

// IsValidIP 检查字符串是否为有效的IP地址
//
// 支持IPv4和IPv6地址格式
//
// 参数:
//   - ipStr: 待检查的IP字符串
//
// 返回值:
//   - bool: 是否为有效IP地址
//
// 使用示例:
//
//	if IsValidIP("192.168.1.1") {
//	    fmt.Println("有效的IP地址")
//	}
//
//	if IsValidIP("2001:db8::1") {
//	    fmt.Println("有效的IPv6地址")
//	}
//
//	if !IsValidIP("invalid") {
//	    fmt.Println("无效的IP地址")
//	}
func IsValidIP(ipStr string) bool {
	return net.ParseIP(ipStr) != nil
}

// IsIPv4 检查是否为IPv4地址
//
// 参数:
//   - ipStr: IP地址字符串
//
// 返回值:
//   - bool: 是否为IPv4地址
//
// 使用示例:
//
//	if IsIPv4("192.168.1.1") {
//	    fmt.Println("这是IPv4地址")
//	}
func IsIPv4(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return ip.To4() != nil
}

// IsIPv6 检查是否为IPv6地址
//
// 参数:
//   - ipStr: IP地址字符串
//
// 返回值:
//   - bool: 是否为IPv6地址
//
// 使用示例:
//
//	if IsIPv6("2001:db8::1") {
//	    fmt.Println("这是IPv6地址")
//	}
func IsIPv6(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return ip.To4() == nil
}

// GetDefaultIPOptions 获取推荐的默认选项组合
//
// 该函数返回一组常用的选项组合，适用于大多数部署在反向代理后面的场景。
//
// 返回值:
//   - options: 推荐的默认选项组合（信任代理）
//
// 使用示例:
//
//	opts := GetDefaultIPOptions()
//	ip := GetClientIP(r, opts)
func GetDefaultIPOptions() IPOption {
	return IPOptionTrustProxy
}

// GetStrictIPOptions 获取严格模式的选项组合
//
// 该函数返回严格模式的选项组合，只从X-Real-IP获取IP。
// 适用于Nginx配置了 proxy_set_header X-Real-IP $remote_addr 的场景。
//
// 返回值:
//   - options: 严格模式的选项组合
//
// 使用示例:
//
//	opts := GetStrictIPOptions()
//	ip := GetClientIP(r, opts)
func GetStrictIPOptions() IPOption {
	return IPOptionTrustProxy | IPOptionStrict
}

// GetInternalIPOptions 获取内网环境的选项组合
//
// 该函数返回适用于内网/开发环境的选项组合，允许返回私有IP。
//
// 返回值:
//   - options: 内网环境的选项组合
//
// 使用示例:
//
//	opts := GetInternalIPOptions()
//	ip := GetClientIP(r, opts)
func GetInternalIPOptions() IPOption {
	return IPOptionTrustProxy | IPOptionAllowPrivate
}