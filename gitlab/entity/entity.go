package entity

import (
	"time"
)

type CustomTime struct {
	time.Time
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	// Pre-allocate buffer with known size for better performance
	buf := make([]byte, 0, 12) // "2006-01-02" + quotes = 12 chars
	buf = append(buf, '"')
	buf = ct.AppendFormat(buf, "2006-01-02")
	buf = append(buf, '"')
	return buf, nil
}

// LicenseInfo struct corresponds to the LicenseInfo class in Java, used to store license information
type LicenseInfo struct {
	Name    string `json:"Name"`
	Company string `json:"Company"`
	Email   string `json:"Email"`
}

// Restriction struct corresponds to the Restriction class in Java, used to store restriction information
type Restriction struct {
	RestrictedUserCount int            `json:"restricted_user_count"`
	ActiveUserCount     int            `json:"active_user_count"`
	Plan                string         `json:"plan"`
	Trial               bool           `json:"trial"`
	ExpiresAt           CustomTime     `json:"expires_at"`
	AddOn               map[string]int `json:"add_ons"`
	Features            []string       `json:"features"`
}

// License represents a license, corresponding to the License class in Java
type License struct {
	// 1. 基础元数据
	Version int `json:"version"`
	// 2. 许可证持有者信息
	License LicenseInfo `json:"licensee"`
	// 3. 时间属性
	IssuedAt       CustomTime `json:"issued_at"`
	StartsAt       CustomTime `json:"starts_at"`
	ExpiresAt      CustomTime `json:"expires_at"`
	NotifyAdminsAt CustomTime `json:"notify_admins_at"`
	NotifyUsersAt  CustomTime `json:"notify_users_at"`
	BlockChangesAt CustomTime `json:"block_changes_at"`
	// 4. 用户数量控制
	RestrictedUserCount int `json:"restricted_user_count"`
	ActiveUserCount     int `json:"active_user_count"`
	// 5. 计划和类型
	Plan  string `json:"plan"`
	Trial bool   `json:"trial"`
	// 6. 附加功能
	AddOn map[string]int `json:"add_ons"`
	// 7. 可用功能列表
	Features []string `json:"features"`
	// 8. 云许可相关
	CloudLicensingEnabled        bool `json:"cloud_licensing_enabled"`
	OfflineCloudLicensingEnabled bool `json:"offline_cloud_licensing_enabled"`
	// 9. 其他控制
	AutoRenewEnabled          bool `json:"auto_renew_enabled"`
	SeatReconciliationEnabled bool `json:"seat_reconciliation_enabled"`
	OperationalMetricsEnabled bool `json:"operational_metrics_enabled"`
	GeneratedFromCustomersDot bool `json:"generated_from_customers_dot"`
	// 10. 其他设置
	Restrictions Restriction `json:"restrictions"`
}
