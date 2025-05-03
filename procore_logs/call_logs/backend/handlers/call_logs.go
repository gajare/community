package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type callLog struct {
	ID              int    `json:"id"`
	Comments        string `json:"comments"`
	Date            string `json:"date"`
	Datetime        string `json:"datetime"`
	InvolvedCompany string `json:"involved_company"`
	InvolvedName    string `json:"involved_name"`
	TimeHour        int    `json:"time_hour"`
	TimeMinute      int    `json:"time_minute"`
	Severity        string `json:"severity"`
	Location        string `json:"location"`
}

type AuthTokenRequest struct {
	Code string `json:"code"`
}

type AuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func GetAuthToken(c *gin.Context) {
	var req AuthTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is required"})
		return
	}

	// Prepare request to Procore's token endpoint
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", os.Getenv("PROCORE_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("PROCORE_CLIENT_SECRET"))
	data.Set("code", req.Code)
	data.Set("redirect_uri", "urn:ietf:wg:oauth:2.0:oob")

	reqURL := "https://login-sandbox.procore.com/oauth/token"
	client := &http.Client{Timeout: 10 * time.Second}

	request, err := http.NewRequest("POST", reqURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token request"})
		return
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get token from Procore"})
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read token response"})
		return
	}

	if response.StatusCode != http.StatusOK {
		c.JSON(response.StatusCode, gin.H{"error": string(body)})
		return
	}

	var tokenResp AuthTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token response"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": tokenResp.AccessToken,
		"token_type":   tokenResp.TokenType,
		"expires_in":   tokenResp.ExpiresIn,
	})
}

func GetcallLogs(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	projectID := os.Getenv("PROCORE_PROJECT_ID")
	companyID := os.Getenv("PROCORE_COMPANY_ID")

	apiUrl := "https://sandbox.procore.com/rest/v1.0/projects/" + projectID + "/call_logs"

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Header.Set("Authorization", accessToken)
	req.Header.Set("Procore-Company-Id", companyID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

func GetcallLogDetails(c *gin.Context) {

	accessToken := c.GetHeader("Authorization")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	logID := c.Param("id")
	if logID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Log ID is required"})
		return
	}

	projectID := os.Getenv("PROCORE_PROJECT_ID")
	companyID := os.Getenv("PROCORE_COMPANY_ID")
	// apiUrl := "https://sandbox.procore.com/rest/v1.0/projects/117922/call_logs/712"

	apiUrl := "https://sandbox.procore.com/rest/v1.0/projects/" + projectID + "/call_logs/" + logID
	fmt.Println("projectID :", projectID)
	fmt.Println("logID :", logID)

	fmt.Println("api url :", apiUrl)
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Header.Set("Authorization", accessToken)
	req.Header.Set("Procore-Company-Id", companyID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

func GetFilteredCallLogs(c *gin.Context) {
	// Get Authorization header
	accessToken := c.GetHeader("Authorization")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	// Extract query parameters - using start_date and end_date now
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	severity := c.Query("severity")
	company := c.Query("company")

	// Get required environment variables
	projectID := os.Getenv("PROCORE_PROJECT_ID")
	companyID := os.Getenv("PROCORE_COMPANY_ID")

	if projectID == "" || companyID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing required environment variables"})
		return
	}

	// Build Procore API URL with date parameters
	baseURL := fmt.Sprintf("https://sandbox.procore.com/rest/v1.0/projects/%s/call_logs", projectID)

	// Create request to Procore API
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request: " + err.Error()})
		return
	}

	// Add query parameters
	q := req.URL.Query()
	if startDate != "" {
		q.Add("start_date", startDate)
	}
	if endDate != "" {
		q.Add("end_date", endDate)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", accessToken)
	req.Header.Set("Procore-Company-Id", companyID)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to contact Procore API: " + err.Error()})
		return
	}

	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response: " + err.Error()})
		return
	}

	// Parse the response
	var logs []map[string]interface{}
	if err := json.Unmarshal(body, &logs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response: " + err.Error()})
		return
	}

	// Apply additional filters (severity and company) locally since Procore API may not support them
	filteredLogs := make([]map[string]interface{}, 0)

	// Get search query parameter
	searchTerm := c.Query("search")

	fmt.Println("Search term:", searchTerm)

	// ... existing code ...

	// Apply search filter
	if searchTerm != "" {
		searchTerm = strings.ToLower(searchTerm)
		filteredBySearch := make([]map[string]interface{}, 0)

		for _, log := range filteredLogs {
			match := false

			// Check all relevant fields
			fieldsToSearch := []string{
				"involved_name",
				"involved_company",
				"comments",
				"location",
				"severity",
			}

			for _, field := range fieldsToSearch {
				if value, ok := log[field].(string); ok {
					if strings.Contains(strings.ToLower(value), searchTerm) {
						match = true
						break
					}
				}
			}

			// Check numeric/date fields if needed
			if !match {
				if date, ok := log["date"].(string); ok {
					if strings.Contains(date, searchTerm) {
						match = true
					}
				}
			}

			if match {
				filteredBySearch = append(filteredBySearch, log)
			}
		}

		filteredLogs = filteredBySearch
	}

	for _, log := range logs {
		// Apply severity filter
		if severity != "" {
			logSeverity, ok := log["severity"].(string)
			if !ok || !strings.EqualFold(logSeverity, severity) {
				continue
			}
		}

		// Apply company filter
		if company != "" {
			logCompany, ok := log["involved_company"].(string)
			if !ok || !strings.Contains(strings.ToLower(logCompany), strings.ToLower(company)) {
				continue
			}
		}

		filteredLogs = append(filteredLogs, log)
	}

	// Return the filtered logs
	c.JSON(http.StatusOK, filteredLogs)
}

func CreateCallLog(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	var logData callLog
	if err := c.ShouldBindJSON(&logData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectID := os.Getenv("PROCORE_PROJECT_ID")
	companyID := os.Getenv("PROCORE_COMPANY_ID")

	formData := url.Values{}
	formData.Set("call_log[comments]", logData.Comments)
	formData.Set("call_log[date]", logData.Date)
	formData.Set("call_log[datetime]", logData.Datetime)
	formData.Set("call_log[involved_company]", logData.InvolvedCompany)
	formData.Set("call_log[involved_name]", logData.InvolvedName)
	formData.Set("call_log[time_hour]", strconv.Itoa(logData.TimeHour))
	formData.Set("call_log[time_minute]", strconv.Itoa(logData.TimeMinute))
	if logData.Severity != "" {
		formData.Set("call_log[severity]", logData.Severity)
	}
	if logData.Location != "" {
		formData.Set("call_log[location]", logData.Location)
	}

	req, err := http.NewRequest("POST", "https://sandbox.procore.com/rest/v1.0/projects/"+projectID+"/call_logs", bytes.NewBufferString(formData.Encode()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Header.Set("Authorization", accessToken)
	req.Header.Set("Procore-Company-Id", companyID)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

func UpdateCallLog(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	logID := c.Param("id")
	if logID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Log ID is required"})
		return
	}

	var logData callLog
	if err := c.ShouldBindJSON(&logData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectID := os.Getenv("PROCORE_PROJECT_ID")
	companyID := os.Getenv("PROCORE_COMPANY_ID")

	formData := url.Values{}
	if logData.Comments != "" {
		formData.Set("call_log[comments]", logData.Comments)
	}
	if logData.Date != "" {
		formData.Set("call_log[date]", logData.Date)
	}
	if logData.Datetime != "" {
		formData.Set("call_log[datetime]", logData.Datetime)
	}
	if logData.InvolvedCompany != "" {
		formData.Set("call_log[involved_company]", logData.InvolvedCompany)
	}
	if logData.InvolvedName != "" {
		formData.Set("call_log[involved_name]", logData.InvolvedName)
	}
	if logData.TimeHour != 0 {
		formData.Set("call_log[time_hour]", strconv.Itoa(logData.TimeHour))
	}
	if logData.TimeMinute != 0 {
		formData.Set("call_log[time_minute]", strconv.Itoa(logData.TimeMinute))
	}
	if logData.Severity != "" {
		formData.Set("call_log[severity]", logData.Severity)
	}
	if logData.Location != "" {
		formData.Set("call_log[location]", logData.Location)
	}

	req, err := http.NewRequest("PUT", "https://sandbox.procore.com/rest/v1.0/projects/"+projectID+"/call_logs/"+logID, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Header.Set("Authorization", accessToken)
	req.Header.Set("Procore-Company-Id", companyID)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

func DeleteCallLog(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	logID := c.Param("id")
	if logID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Log ID is required"})
		return
	}

	projectID := os.Getenv("PROCORE_PROJECT_ID")
	companyID := os.Getenv("PROCORE_COMPANY_ID")

	req, err := http.NewRequest("DELETE", "https://sandbox.procore.com/rest/v1.0/projects/"+projectID+"/call_logs/"+logID, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Header.Set("Authorization", accessToken)
	req.Header.Set("Procore-Company-Id", companyID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}
