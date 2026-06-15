package gitlab

import (
	"time"
)

// CustomTime serializes as a date string ("2006-01-02").
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

// CustomDateTime serializes as an ISO-8601 UTC timestamp ("2006-01-02T15:04:05Z").
// Used by cloud-licensing fields such as activated_at / last_synced_at / next_sync_at.
type CustomDateTime struct {
	time.Time
}

func (ct CustomDateTime) MarshalJSON() ([]byte, error) {
	buf := make([]byte, 0, 22) // "2006-01-02T15:04:05Z" + quotes
	buf = append(buf, '"')
	buf = ct.UTC().AppendFormat(buf, "2006-01-02T15:04:05Z")
	buf = append(buf, '"')
	return buf, nil
}

// LicenseInfo struct corresponds to the LicenseInfo class in Java, used to store license information
type LicenseInfo struct {
	Name    string `json:"Name"`
	Company string `json:"Company"`
	Email   string `json:"Email"`
}

// AddOnPurchase represents a single purchase entry under an add-on product
// in the GitLab cloud-licensing payload (restrictions.add_on_products[name][]).
type AddOnPurchase struct {
	Quantity    int        `json:"quantity"`
	StartedOn   CustomTime `json:"started_on"`
	ExpiresOn   CustomTime `json:"expires_on"`
	PurchaseXID string     `json:"purchase_xid"`
	Trial       bool       `json:"trial"`
}

// Restriction stores GitLab cloud-licensing restriction information.
type Restriction struct {
	ID                      string                     `json:"id"`
	Plan                    string                     `json:"plan"`
	ActiveUserCount         int                        `json:"active_user_count"`
	SubscriptionID          string                     `json:"subscription_id"`
	SubscriptionName        string                     `json:"subscription_name"`
	ReconciliationCompleted bool                       `json:"reconciliation_completed"`
	AddOnProducts           map[string][]AddOnPurchase `json:"add_on_products"`
}

// License represents a GitLab cloud / offline-cloud license payload.
type License struct {
	// 1. 基础元数据
	Version int `json:"version"`
	// 2. 许可证持有者信息
	License LicenseInfo `json:"licensee"`
	// 3. 时间属性（按日精度）
	IssuedAt       CustomTime `json:"issued_at"`
	ExpiresAt      CustomTime `json:"expires_at"`
	NotifyAdminsAt CustomTime `json:"notify_admins_at"`
	NotifyUsersAt  CustomTime `json:"notify_users_at"`
	BlockChangesAt CustomTime `json:"block_changes_at"`
	// 4. 云激活相关时间属性（按秒精度）
	ActivatedAt  CustomDateTime `json:"activated_at"`
	LastSyncedAt CustomDateTime `json:"last_synced_at"`
	NextSyncAt   CustomDateTime `json:"next_sync_at"`
	// 5. 云许可开关
	CloudLicensingEnabled        bool `json:"cloud_licensing_enabled"`
	OfflineCloudLicensingEnabled bool `json:"offline_cloud_licensing_enabled"`
	// 6. 其他控制
	AutoRenewEnabled          bool `json:"auto_renew_enabled"`
	SeatReconciliationEnabled bool `json:"seat_reconciliation_enabled"`
	OperationalMetricsEnabled bool `json:"operational_metrics_enabled"`
	GeneratedFromCustomersDot bool `json:"generated_from_customers_dot"`
	GeneratedFromCancellation bool `json:"generated_from_cancellation"`
	TemporaryExtension        bool `json:"temporary_extension"`
	ContractOveragesAllowed   bool `json:"contract_overages_allowed"`
	// 7. 限制 / Add-on
	Restrictions Restriction `json:"restrictions"`
}
