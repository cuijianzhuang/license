package service

import (
	"archive/zip"
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"license/config"
	"license/crypto"
	"license/gitlab/entity"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	gorsa "github.com/Lyafei/go-rsa"
	"github.com/gin-gonic/gin"
)

// addOn is a map of add-on features with their respective limits
var addOn = map[string]int{
	"GitLab_Auditor_User": 10000,
	"GitLab_FileLocks":    10000,
	"GitLab_Geo":          10000,
}

// features is a list of features available in the license
var features = []string{
	"admin_audit_log",
	"amazon_q",
	"auditor_user",
	"custom_file_templates",
	"custom_project_templates",
	"db_load_balancing",
	"default_branch_protection_restriction_in_groups",
	"elastic_search",
	"enterprise_templates",
	"extended_audit_events",
	"external_authorization_service_api_management",
	"geo",
	"git_abuse_rate_limit",
	"instance_level_scim",
	"integrations_allow_list",
	"ldap_group_sync",
	"ldap_group_sync_filter",
	"multiple_ldap_servers",
	"object_storage",
	"pages_size_limit",
	"password_complexity",
	"project_aliases",
	"repository_size_limit",
	"required_ci_templates",
	"runner_maintenance_note",
	"runner_performance_insights",
	"runner_upgrade_management",
	"seat_link",
	"seat_usage_quotas",
	"pipelines_usage_quotas",
	"transfer_usage_quotas",
	"product_analytics_usage_quotas",
	"zoekt_code_search",
	"disable_private_profiles",
	"observability_alerts",
	"audit_events",
	"blocked_issues",
	"blocked_work_items",
	"board_iteration_lists",
	"code_owners",
	"code_review_analytics",
	"full_codequality_report",
	"group_activity_analytics",
	"group_bulk_edit",
	"issuable_default_templates",
	"issue_weights",
	"iterations",
	"merge_request_approvers",
	"milestone_charts",
	"multiple_issue_assignees",
	"multiple_merge_request_assignees",
	"multiple_merge_request_reviewers",
	"project_merge_request_analytics",
	"protected_refs_for_users",
	"push_rules",
	"resource_access_token",
	"seat_control",
	"wip_limits",
	"description_diffs",
	"send_emails_from_admin_area",
	"maintenance_mode",
	"scoped_issue_board",
	"contribution_analytics",
	"group_webhooks",
	"member_lock",
	"repository_mirrors",
	"ai_chat",
	"adjourned_deletion_for_projects_and_groups",
	"agent_managed_resources",
	"blocking_merge_requests",
	"board_assignee_lists",
	"board_milestone_lists",
	"ci_secrets_management",
	"ci_pipeline_cancellation_restrictions",
	"cluster_agents_ci_impersonation",
	"cluster_agents_user_impersonation",
	"cluster_deployments",
	"code_owner_approval_required",
	"code_suggestions",
	"commit_committer_check",
	"commit_committer_name_check",
	"compliance_framework",
	"custom_compliance_frameworks",
	"custom_fields",
	"cycle_analytics_for_groups",
	"cycle_analytics_for_projects",
	"default_project_deletion_protection",
	"delete_unconfirmed_users",
	"dependency_proxy_for_packages",
	"disable_extensions_marketplace_for_enterprise_users",
	"disable_name_update_for_users",
	"disable_personal_access_tokens",
	"domain_verification",
	"epic_colors",
	"epics",
	"feature_flags_code_references",
	"file_locks",
	"generic_alert_fingerprinting",
	"git_two_factor_enforcement",
	"group_allowed_email_domains",
	"group_coverage_reports",
	"group_forking_protection",
	"group_level_compliance_dashboard",
	"group_milestone_project_releases",
	"group_project_templates",
	"group_repository_analytics",
	"group_saml",
	"group_scoped_ci_variables",
	"ide_schema_config",
	"incident_metric_upload",
	"jira_issues_integration",
	"linked_items_epics",
	"merge_request_performance_metrics",
	"admin_merge_request_approvers_rules",
	"merge_trains",
	"metrics_reports",
	"multiple_alert_http_integrations",
	"multiple_approval_rules",
	"multiple_group_issue_boards",
	"microsoft_group_sync",
	"operations_dashboard",
	"package_forwarding",
	"packages_virtual_registry",
	"pages_multiple_versions",
	"productivity_analytics",
	"protected_environments",
	"reject_non_dco_commits",
	"reject_unsigned_commits",
	"related_epics",
	"remote_development",
	"saml_group_sync",
	"service_accounts",
	"scoped_labels",
	"smartcard_auth",
	"ssh_certificates",
	"swimlanes",
	"target_branch_rules",
	"troubleshoot_job",
	"type_of_work_analytics",
	"minimal_access_role",
	"unprotection_restrictions",
	"ci_project_subscriptions",
	"incident_timeline_view",
	"oncall_schedules",
	"escalation_policies",
	"zentao_issues_integration",
	"coverage_check_approval_rule",
	"issuable_resource_links",
	"group_protected_branches",
	"group_level_merge_checks_setting",
	"oidc_client_groups_claim",
	"disable_deleting_account_for_users",
	"group_saved_replies",
	"requested_changes_block_merge_request",
	"project_saved_replies",
	"default_roles_assignees",
	"ci_component_usages_in_projects",
	"branch_rule_squash_options",
	"work_item_status",
	"glab_ask_git_command",
	"generate_commit_message",
	"summarize_new_merge_request",
	"summarize_review",
	"generate_description",
	"summarize_comments",
	"review_merge_request",
	"group_ip_restriction",
	"issues_analytics",
	"group_wikis",
	"email_additional_text",
	"custom_file_templates_for_namespace",
	"incident_sla",
	"export_user_permissions",
	"cross_project_pipelines",
	"feature_flags_related_issues",
	"merge_pipelines",
	"ci_cd_projects",
	"github_integration",
	"ai_agents",
	"ai_config_chat",
	"ai_features",
	"ai_review_mr",
	"ai_workflows",
	"api_discovery",
	"api_fuzzing",
	"auto_rollback",
	"cluster_receptive_agents",
	"cluster_image_scanning",
	"external_status_checks",
	"combined_project_analytics_dashboards",
	"compliance_pipeline_configuration",
	"container_scanning",
	"credentials_inventory",
	"custom_roles",
	"dast",
	"dependency_scanning",
	"dora4_analytics",
	"description_composer",
	"environment_alerts",
	"evaluate_group_level_compliance_pipeline",
	"explain_code",
	"external_audit_events",
	"experimental_features",
	"generate_test_file",
	"ai_generate_cube_query",
	"group_ci_cd_analytics",
	"group_level_compliance_adherence_report",
	"group_level_compliance_violations_report",
	"project_level_compliance_dashboard",
	"project_level_compliance_adherence_report",
	"project_level_compliance_violations_report",
	"group_level_analytics_dashboard",
	"incident_management",
	"inline_codequality",
	"insights",
	"issuable_health_status",
	"issues_completed_analytics",
	"jira_vulnerabilities_integration",
	"jira_issue_association_enforcement",
	"kubernetes_cluster_vulnerabilities",
	"license_scanning",
	"okrs",
	"personal_access_token_expiration_policy",
	"secret_push_protection",
	"product_analytics",
	"project_quality_summary",
	"project_level_analytics_dashboard",
	"quality_management",
	"release_evidence_test_artifacts",
	"report_approver_rules",
	"requirements",
	"runner_performance_insights_for_namespace",
	"runner_upgrade_management_for_namespace",
	"sast",
	"sast_advanced",
	"sast_iac",
	"sast_custom_rulesets",
	"sast_fp_reduction",
	"secret_detection",
	"security_configuration_in_ui",
	"security_dashboard",
	"security_inventory",
	"security_on_demand_scans",
	"security_orchestration_policies",
	"security_training",
	"ssh_key_expiration_policy",
	"summarize_mr_changes",
	"stale_runner_cleanup_for_namespace",
	"status_page",
	"suggested_reviewers",
	"subepics",
	"observability",
	"unique_project_download_limit",
	"vulnerability_finding_signatures",
	"container_scanning_for_registry",
	"security_exclusions",
	"security_scans_api",
	"measure_comment_temperature",
	"coverage_fuzzing",
	"devops_adoption",
	"group_level_devops_adoption",
	"instance_level_devops_adoption",
}

// KeyManager manages RSA keys with lazy loading and caching
type KeyManager struct {
	privateKey []byte
	publicKey  []byte
	once       sync.Once
	mutex      sync.RWMutex
	err        error
}

var keyManager = &KeyManager{}

// getKeys returns cached keys with lazy loading
func (km *KeyManager) getKeys() ([]byte, []byte, error) {
	km.once.Do(func() {
		if publicBytes, err := os.ReadFile(config.GetConfig().DataDir + "/.license_encryption_key.pub"); err == nil {
			km.publicKey = publicBytes
		} else {
			km.err = fmt.Errorf("failed to read public key: %v", err)
			return
		}
		if privateBytes, err := os.ReadFile(config.GetConfig().DataDir + "/.license_decryption_key.pri"); err == nil {
			km.privateKey = privateBytes
		} else {
			km.err = fmt.Errorf("failed to read private key: %v", err)
			return
		}
	})

	km.mutex.RLock()
	defer km.mutex.RUnlock()
	return km.privateKey, km.publicKey, km.err
}

// LoadKeys initializes the key manager (for backward compatibility)
func LoadKeys() error {
	_, _, err := keyManager.getKeys()
	return err
}

// Buffer pools for efficient memory management
var (
	jsonBufferPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 4096))
		},
	}
	ioBufferPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 32*1024) // 32KB buffer
		},
	}
	aesKeyPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 16)
		},
	}
	ivPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, aes.BlockSize)
		},
	}
)

// LicenseCache implements a simple TTL cache for license generation
type LicenseCache struct {
	cache      *sync.Map
	ttl        time.Duration
	maxEntries int64
	entries    int64
}

type CacheEntry struct {
	data      []byte
	timestamp time.Time
}

var licenseCache = &LicenseCache{
	cache:      &sync.Map{},
	ttl:        15 * time.Minute,
	maxEntries: 1000,
}

func (lc *LicenseCache) Get(key string) ([]byte, bool) {
	if value, ok := lc.cache.Load(key); ok {
		entry := value.(CacheEntry)
		if time.Since(entry.timestamp) < lc.ttl {
			return entry.data, true
		}
		lc.cache.Delete(key)
		atomic.AddInt64(&lc.entries, -1)
	}
	return nil, false
}

func (lc *LicenseCache) Set(key string, data []byte) {
	if atomic.LoadInt64(&lc.entries) >= lc.maxEntries {
		return // Simple eviction
	}

	cachedData := make([]byte, len(data))
	copy(cachedData, data)

	lc.cache.Store(key, CacheEntry{
		data:      cachedData,
		timestamp: time.Now(),
	})
	atomic.AddInt64(&lc.entries, 1)
}

// createLicenseJson creates a JSON representation of the license with optimized buffer usage
func createLicenseJson(licenseInfo entity.LicenseInfo, expireTime string) ([]byte, error) {

	var expirationDate time.Time
	var err error
	if len(expireTime) == 0 {
		// Default expiration time is 2 years
		expirationDate = time.Now().AddDate(2, 0, 0)
	} else {
		expirationDate, err = time.Parse(time.DateTime, expireTime)
		if err != nil {
			log.Printf("Failed to parse expiration time: %v", err)
			return nil, err
		}
	}

	license := entity.License{
		Version:                      1,
		License:                      licenseInfo,
		IssuedAt:                     entity.CustomTime{Time: time.Now()},
		StartsAt:                     entity.CustomTime{Time: time.Now()},
		ExpiresAt:                    entity.CustomTime{Time: expirationDate},
		NotifyAdminsAt:               entity.CustomTime{Time: expirationDate},
		NotifyUsersAt:                entity.CustomTime{Time: expirationDate},
		BlockChangesAt:               entity.CustomTime{Time: expirationDate},
		RestrictedUserCount:          10000,
		ActiveUserCount:              10000,
		Plan:                         "ultimate",
		Trial:                        false,
		AddOn:                        addOn,
		Features:                     features,
		CloudLicensingEnabled:        false,
		OfflineCloudLicensingEnabled: false,
		AutoRenewEnabled:             false,
		SeatReconciliationEnabled:    false,
		OperationalMetricsEnabled:    false,
		GeneratedFromCustomersDot:    false,
		Restrictions: entity.Restriction{
			RestrictedUserCount: 10000,
			ActiveUserCount:     10000,
			Plan:                "ultimate",
			Trial:               false,
			ExpiresAt:           entity.CustomTime{Time: expirationDate},
			AddOn:               addOn,
			Features:            features,
		},
	}

	buf := jsonBufferPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		jsonBufferPool.Put(buf)
	}()

	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(license); err != nil {
		return nil, err
	}

	// Remove trailing newline added by encoder
	data := buf.Bytes()
	if len(data) > 0 && data[len(data)-1] == '\n' {
		data = data[:len(data)-1]
	}

	result := make([]byte, len(data))
	copy(result, data)
	return result, nil
}

// generateRandomIV generates a random initialization vector (IV)
func generateRandomIV() ([]byte, error) {
	iv := make([]byte, aes.BlockSize) // AES block size is fixed at 16 bytes
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	return iv, nil
}

// Encrypt wraps the Encrypt method, using AES-CBC encryption and PKCS7 padding
func Encrypt(data, key, iv []byte) ([]byte, error) {
	aesTool := crypto.AesCbcPkcs7{Key: key, Iv: iv}
	enc, err := aesTool.Encrypt(data)
	if err != nil {
		log.Println("Encrypt error:", err)
		return nil, err
	}
	return enc, err
}

// Uses RSA private key to "encrypt" data with cached keys
func encryptWithPrivateKey(data string) (string, error) {
	privateKey, _, err := keyManager.getKeys()
	if err != nil {
		return "", err
	}
	encrypt, err := gorsa.PriKeyEncrypt(data, string(privateKey))
	if err != nil {
		log.Printf("Failed to encrypt data with RSA private key: %v", err)
		return "", err
	}
	return encrypt, nil
}

// encryptLicense encrypts license data using AES and RSA with pooled resources
func encryptLicense(data []byte) (string, error) {
	// Get pooled AES key and IV
	key := aesKeyPool.Get().([]byte)
	defer aesKeyPool.Put(key)

	iv := ivPool.Get().([]byte)
	defer ivPool.Put(iv)

	// Generate fresh random data
	if _, err := rand.Read(key); err != nil {
		log.Printf("Failed to generate AES key: %v", err)
		return "", err
	}

	if _, err := rand.Read(iv); err != nil {
		log.Printf("Failed to generate AES IV: %v", err)
		return "", err
	}

	encryptedData, err := Encrypt(data, key, iv)
	if err != nil {
		log.Printf("Failed to encrypt data: %v", err)
		return "", err
	}

	// Note: RSA encryption is typically done with a public key, but technically can be done with a private key (although not recommended)
	encryptedKey, err := encryptWithPrivateKey(string(key))
	if err != nil {
		return "", err
	}

	// Encode encrypted data as Base64
	encryptedDataStr := base64.StdEncoding.EncodeToString(encryptedData)
	ivStr := base64.StdEncoding.EncodeToString(iv)

	// Package as JSON format
	result := map[string]string{
		"data": encryptedDataStr,
		"key":  encryptedKey,
		"iv":   ivStr,
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Printf("Failed to package JSON data: %v", err)
		return "", err
	}

	// Encode JSON as Base64
	encodedFinal := base64.StdEncoding.EncodeToString(jsonData)
	return encodedFinal, nil
}

// Generate generates a license and sends it via HTTP response
func Generate(ctx *gin.Context, licenseInfo entity.LicenseInfo, expireTime string) {
	createLicense(ctx, licenseInfo, expireTime)
}

// createLicense creates and sends a license with caching support
func createLicense(ctx *gin.Context, licenseInfo entity.LicenseInfo, expireTime string) {
	// Create cache key
	cacheKey := fmt.Sprintf("%s:%s:%s:%s", licenseInfo.Name, licenseInfo.Email, licenseInfo.Company, expireTime)

	// Check cache first
	if cachedLicense, found := licenseCache.Get(cacheKey); found {
		ctx.Header("Content-Disposition", "attachment; filename=license.zip")
		ctx.Header("Content-Type", "application/zip")
		ctx.Header("X-Cache", "HIT")
		ctx.Data(http.StatusOK, "application/zip", cachedLicense)
		return
	}

	// Create license JSON data
	licenseJson, err := createLicenseJson(licenseInfo, expireTime)
	if err != nil {
		log.Printf("Failed to create license JSON: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Encrypt the license data
	encryptedLicense, err := encryptLicense(licenseJson)
	if err != nil {
		log.Printf("Failed to encrypt license: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Create ZIP file in memory first for caching
	buf := &bytes.Buffer{}
	zipWriter := zip.NewWriter(buf)

	// Add public key file to ZIP
	if err := addFileToZipOptimized(zipWriter, config.GetConfig().DataDir+"/.license_encryption_key.pub", "license/.license_encryption_key.pub"); err != nil {
		log.Printf("Failed to add public key to ZIP: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Add encrypted license data to ZIP
	if err := addLicenseToZip(zipWriter, encryptedLicense, "license/license.gitlab-license"); err != nil {
		log.Printf("Failed to add license to ZIP: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if err := zipWriter.Close(); err != nil {
		log.Printf("Failed to close ZIP writer: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	zipData := buf.Bytes()

	// Cache the result
	licenseCache.Set(cacheKey, zipData)

	// Send response
	ctx.Header("Content-Disposition", "attachment; filename=license.zip")
	ctx.Header("Content-Type", "application/zip")
	ctx.Header("X-Cache", "MISS")
	ctx.Data(http.StatusOK, "application/zip", zipData)
}

// exportZipStream creates and sends a ZIP file containing the encrypted license and public key file
// This function is kept for backward compatibility but is no longer used in the optimized flow
func exportZipStream(ctx *gin.Context, encryptedLicense string) error {
	// Set response headers for file download
	ctx.Status(http.StatusOK) // Explicitly set status code to 200 OK
	ctx.Header("Content-Disposition", "attachment; filename=license.zip")
	ctx.Header("Content-Type", "application/zip")

	zipWriter := zip.NewWriter(ctx.Writer)
	defer func(zipWriter *zip.Writer) {
		err := zipWriter.Close()
		if err != nil {
			log.Printf("Failed to close ZIP writer: %v", err)
		}
	}(zipWriter)

	// Add public key file to ZIP
	if err := addFileToZipOptimized(zipWriter, config.GetConfig().DataDir+"/.license_encryption_key.pub", "license/.license_encryption_key.pub"); err != nil {
		return err
	}

	// Add encrypted license data to ZIP
	if err := addLicenseToZip(zipWriter, encryptedLicense, "license/license.gitlab-license"); err != nil {
		return err
	}

	return nil
}

// addFileToZipOptimized reads a file from the filesystem and adds it to the ZIP with buffered I/O
func addFileToZipOptimized(zipWriter *zip.Writer, filePath, zipPath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}
	header.Name = zipPath
	header.Method = zip.Deflate
	header.Modified = fileInfo.ModTime()

	zipFile, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Use pooled buffer for efficient copying
	buffer := ioBufferPool.Get().([]byte)
	defer ioBufferPool.Put(buffer)

	_, err = io.CopyBuffer(zipFile, file, buffer)
	return err
}

// addLicenseToZip directly writes string data to a ZIP entry
func addLicenseToZip(zipWriter *zip.Writer, data, zipPath string) error {
	// Create a new zip.FileHeader, set filename and modification time
	header := &zip.FileHeader{
		Name:     zipPath,
		Method:   zip.Deflate, // Use compression to reduce file size
		Modified: time.Now(),  // Set current time as file modification time
	}

	// Create ZIP entry
	zipFile, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Write data to ZIP entry
	_, err = zipFile.Write([]byte(data))
	return err
}
