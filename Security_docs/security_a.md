# Real-Time Threat Detection and Response for InvoiceB2B

### Overview

This document outlines the implementation of a comprehensive real-time threat detection and response system for the InvoiceB2B application. Building upon the existing security measures and monitoring infrastructure, this system will provide continuous surveillance, automated threat detection, and coordinated incident response capabilities to protect the application from emerging security threats.

### Architecture Components

The real-time threat detection and response system will integrate with the existing infrastructure while adding specialized security monitoring components:

#### 1. Security Information and Event Management (SIEM)

**Implementation:**
- Deploy AWS Security Hub as the central SIEM solution
- Integrate with existing CloudWatch logs and metrics
- Configure custom security dashboards in Grafana for visualization

**Data Sources:**
- Application logs (authentication events, API access, admin actions)
- Infrastructure logs (ECS, RDS, ElastiCache, ALB)
- AWS CloudTrail for AWS API activity
- VPC Flow Logs for network traffic analysis
- AWS GuardDuty findings

#### 2. Web Application Firewall (WAF)

**Implementation:**
- Deploy AWS WAF in front of the Application Load Balancer
- Configure core rule sets for common attack patterns (OWASP Top 10)
- Implement custom rules for application-specific threats
- Enable logging to CloudWatch for integration with SIEM

#### 3. Runtime Application Self-Protection (RASP)

**Implementation:**
- Integrate a Go-compatible RASP solution into the application containers
- Configure to monitor and protect against:
    - SQL injection attempts
    - Cross-site scripting (XSS)
    - Command injection
    - Unauthorized file access
    - Abnormal API usage patterns

#### 4. User and Entity Behavior Analytics (UEBA)

**Implementation:**
- Deploy a UEBA solution that integrates with the application's activity logs
- Establish baseline behavior for users and system components
- Configure anomaly detection for:
    - Login patterns (time, location, frequency)
    - Transaction volumes and amounts
    - API usage patterns
    - Admin activity

### Threat Detection Mechanisms

#### 1. Authentication and Authorization Monitoring

**Implementation:**
```go
// Add to auth_middleware.go
func (am *AuthMiddleware) LogAuthEvent(c *fiber.Ctx, event string, success bool, details map[string]interface{}) {
    // Extract IP address, user agent, etc.
    ipAddress := c.IP()
    userAgent := c.Get("User-Agent")
    
    // Create structured log entry
    logEntry := map[string]interface{}{
        "event_type": "auth",
        "sub_type": event,
        "success": success,
        "ip_address": ipAddress,
        "user_agent": userAgent,
        "timestamp": time.Now().UTC(),
        "details": details,
    }
    
    // Log to structured logging system
    logger.Info("auth_event", logEntry)
    
    // Publish to security event stream for real-time analysis
    am.securityEventService.PublishEvent("auth", logEntry)
}
```

**Detection Rules:**
- Multiple failed login attempts from the same IP
- Successful logins from unusual locations or devices
- Brute force attempts against multiple accounts
- Unusual admin privilege usage patterns
- JWT token manipulation attempts

#### 2. API Abuse Detection

**Implementation:**
- Add middleware to track API usage patterns:

```go
// Add to api_security_middleware.go
func (asm *APISecurityMiddleware) MonitorAPIUsage() fiber.Handler {
    return func(c *fiber.Ctx) error {
        startTime := time.Now()
        
        // Store original context for later analysis
        requestBody := string(c.Body())
        path := c.Path()
        method := c.Method()
        
        // Process the request
        err := c.Next()
        
        // Record metrics
        duration := time.Since(startTime)
        statusCode := c.Response().StatusCode()
        
        // Log API usage for security analysis
        asm.securityMetricsService.RecordAPIUsage(
            path, 
            method, 
            statusCode, 
            duration, 
            c.Locals("user_id"), 
            c.IP(),
        )
        
        return err
    }
}
```

**Detection Rules:**
- Abnormal request rates or patterns
- Suspicious parameter values (potential injection attempts)
- Unusual API sequence calls
- Excessive error responses
- Scanning behavior (accessing multiple resources rapidly)

#### 3. Invoice and Financial Fraud Detection

**Implementation:**
- Add transaction monitoring service:

```go
// Add to invoice_service.go
func (is *InvoiceService) CheckForFraudIndicators(invoice *models.Invoice) []string {
    var indicators []string
    
    // Check for known fraud patterns
    if time.Until(invoice.DueDate).Hours() < 24 {
        indicators = append(indicators, "extremely_short_payment_term")
    }
    
    if invoice.Amount > is.thresholdService.GetHighValueThreshold() {
        indicators = append(indicators, "high_value_transaction")
    }
    
    // Check for unusual patterns for this user
    userStats, _ := is.userStatsService.GetInvoiceStats(invoice.UserID)
    if invoice.Amount > userStats.AverageAmount*3 {
        indicators = append(indicators, "amount_significantly_above_user_average")
    }
    
    // Log potential fraud indicators
    if len(indicators) > 0 {
        is.securityEventService.PublishEvent("potential_fraud", map[string]interface{}{
            "invoice_id": invoice.ID,
            "user_id": invoice.UserID,
            "indicators": indicators,
            "timestamp": time.Now().UTC(),
        })
    }
    
    return indicators
}
```

**Detection Rules:**
- Unusual invoice amounts or payment terms
- Multiple invoices with similar characteristics but different debtors
- Rapid succession of invoice submissions
- Invoices with suspicious or inconsistent data
- Pattern matching against known fraud schemes

#### 4. File Upload Protection

**Implementation:**
- Enhance file upload security:

```go
// Add to file_service.go
func (fs *FileService) ScanUploadedFile(fileData []byte, fileName string) (bool, []string) {
    var threats []string
    
    // Check file type and content
    mimeType := http.DetectContentType(fileData)
    if !fs.isAllowedMimeType(mimeType) {
        threats = append(threats, "mime_type_mismatch")
    }
    
    // Scan for malware using integrated scanner
    malwareDetected, malwareType := fs.malwareScanner.Scan(fileData)
    if malwareDetected {
        threats = append(threats, "malware_detected:"+malwareType)
    }
    
    // Check for embedded scripts or macros in PDFs
    if strings.HasSuffix(strings.ToLower(fileName), ".pdf") {
        if fs.containsJavaScript(fileData) {
            threats = append(threats, "pdf_contains_javascript")
        }
    }
    
    // Log security event if threats detected
    if len(threats) > 0 {
        fs.securityEventService.PublishEvent("malicious_file", map[string]interface{}{
            "file_name": fileName,
            "mime_type": mimeType,
            "threats": threats,
            "timestamp": time.Now().UTC(),
        })
    }
    
    return len(threats) == 0, threats
}
```

**Detection Rules:**
- MIME type mismatches
- Malware signatures
- Embedded scripts or macros
- Unusually large files
- Encrypted content in unexpected places

### Automated Response Actions

#### 1. Account Protection

**Implementation:**
```go
// Add to security_response_service.go
func (srs *SecurityResponseService) HandleAuthenticationThreat(userID uint, threatType string, evidence map[string]interface{}) {
    switch threatType {
    case "brute_force_attempt":
        // Temporarily lock account
        srs.userService.LockAccount(userID, "security_lock", 30*time.Minute)
        
        // Notify user
        srs.notificationService.SendSecurityAlert(userID, "unusual_login_activity")
        
        // Log incident
        srs.securityIncidentService.LogIncident("brute_force", userID, evidence)
        
    case "suspicious_location":
        // Require additional verification
        srs.userService.RequireAdditionalVerification(userID)
        
        // Notify user
        srs.notificationService.SendSecurityAlert(userID, "login_from_new_location")
        
    case "token_manipulation":
        // Invalidate all user sessions
        srs.sessionService.InvalidateAllUserSessions(userID)
        
        // Require password reset
        srs.userService.RequirePasswordReset(userID)
        
        // Log high severity incident
        srs.securityIncidentService.LogIncident("token_manipulation", userID, evidence, "high")
    }
}
```

**Response Actions:**
- Temporary account lockout after multiple failed login attempts
- Step-up authentication for suspicious login locations
- Session invalidation for detected token manipulation
- Forced password reset for compromised accounts
- User notifications of suspicious activity

#### 2. API Abuse Mitigation

**Implementation:**
```go
// Add to rate_limiter_service.go
func (rls *RateLimiterService) DynamicRateLimit(c *fiber.Ctx) error {
    // Get identifier (IP or user ID if authenticated)
    identifier := c.IP()
    if userID := c.Locals("user_id"); userID != nil {
        identifier = fmt.Sprintf("user:%v", userID)
    }
    
    // Check current threat level for this identifier
    threatLevel := rls.threatIntelService.GetThreatLevel(identifier)
    
    // Apply rate limits based on threat level
    var limit, window int
    switch threatLevel {
    case "high":
        limit = 10
        window = 60 // 10 requests per minute
    case "medium":
        limit = 30
        window = 60 // 30 requests per minute
    case "low":
        limit = 100
        window = 60 // 100 requests per minute
    default:
        limit = 300
        window = 60 // 300 requests per minute (normal)
    }
    
    // Check if limit exceeded
    count, _ := rls.redisClient.Incr(fmt.Sprintf("ratelimit:%s:%d", identifier, time.Now().Unix()/window))
    if count > int64(limit) {
        // Log rate limit event
        rls.securityEventService.PublishEvent("rate_limit_exceeded", map[string]interface{}{
            "identifier": identifier,
            "threat_level": threatLevel,
            "limit": limit,
            "window": window,
            "timestamp": time.Now().UTC(),
        })
        
        return utils.HandleError(c, fiber.StatusTooManyRequests, "Rate limit exceeded", nil)
    }
    
    return c.Next()
}
```

**Response Actions:**
- Dynamic rate limiting based on threat intelligence
- IP blocking for persistent abusers
- CAPTCHA challenges for suspicious traffic patterns
- Request throttling for endpoints under attack
- Temporary service restrictions for suspicious users

#### 3. Fraud Prevention

**Implementation:**
```go
// Add to invoice_security_service.go
func (iss *InvoiceSecurityService) HandlePotentialFraud(invoiceID uint, indicators []string) {
    // Get invoice details
    invoice, _ := iss.invoiceRepo.FindByID(context.Background(), invoiceID)
    
    // Determine risk level
    riskLevel := iss.calculateRiskLevel(indicators)
    
    switch riskLevel {
    case "high":
        // Automatically reject high-risk invoices
        iss.invoiceService.UpdateStatus(invoiceID, models.InvoiceStatusRejected, "Automated security rejection: High risk indicators detected")
        
        // Flag account for review
        iss.userSecurityService.FlagForReview(invoice.UserID, "multiple_high_risk_invoices")
        
        // Alert fraud team
        iss.alertService.SendFraudAlert("high_risk_invoice", map[string]interface{}{
            "invoice_id": invoiceID,
            "user_id": invoice.UserID,
            "indicators": indicators,
        })
        
    case "medium":
        // Move to manual review queue with priority
        iss.reviewQueueService.AddToQueue(invoiceID, "security_review", "high")
        
        // Add security notes for reviewer
        iss.invoiceService.AddSecurityNote(invoiceID, fmt.Sprintf("Security review required: %s", strings.Join(indicators, ", ")))
        
    case "low":
        // Add security note but allow normal processing
        iss.invoiceService.AddSecurityNote(invoiceID, fmt.Sprintf("Low risk indicators detected: %s", strings.Join(indicators, ", ")))
    }
    
    // Log security event
    iss.securityEventService.LogFraudDetection(invoiceID, riskLevel, indicators)
}
```

**Response Actions:**
- Automatic rejection of high-risk invoices
- Routing suspicious invoices to manual review
- Temporary suspension of financing for accounts with multiple fraud indicators
- Enhanced verification requirements for suspicious transactions
- Notification to finance team for manual review

#### 4. Infrastructure Protection

**Implementation:**
```go
// Add to infrastructure_security_service.go
func (iss *InfrastructureSecurityService) HandleInfrastructureThreat(threatType string, resource string, details map[string]interface{}) {
    switch threatType {
    case "ddos_attack":
        // Enable AWS Shield Advanced protection
        iss.awsService.EnableShieldProtection(resource)
        
        // Scale up resources to handle load
        iss.autoScalingService.IncreaseCapacity(resource, "security_event")
        
        // Notify operations team
        iss.alertService.SendOperationsAlert("ddos_mitigation_active", details)
        
    case "unauthorized_access_attempt":
        // Tighten security group rules
        iss.securityGroupService.ApplyRestrictiveRules(resource)
        
        // Rotate access credentials
        iss.credentialService.RotateCredentials(resource)
        
        // Initiate security audit
        iss.securityAuditService.InitiateAudit(resource, "unauthorized_access")
        
    case "data_exfiltration":
        // Block suspicious outbound connections
        iss.networkControlService.BlockOutboundTraffic(details["destination"].(string))
        
        // Isolate affected resources
        iss.containerService.IsolateContainer(resource)
        
        // Trigger incident response plan
        iss.incidentResponseService.TriggerResponse("data_breach", details)
    }
    
    // Log infrastructure security event
    iss.securityEventService.LogInfrastructureEvent(threatType, resource, details)
}
```

**Response Actions:**
- Dynamic scaling to mitigate DDoS attacks
- Automatic security group rule tightening
- Container isolation for compromised services
- Credential rotation for suspected unauthorized access
- Network traffic blocking for data exfiltration attempts

### Incident Response Workflow

#### 1. Detection Phase

- Security events collected from all monitoring sources
- Events correlated and analyzed by SIEM
- Alerts generated based on detection rules
- Initial severity classification assigned

#### 2. Triage Phase

**Implementation:**
```go
// Add to incident_response_service.go
func (irs *IncidentResponseService) TriageSecurityIncident(incidentID string) {
    // Get incident details
    incident, _ := irs.incidentRepo.FindByID(incidentID)
    
    // Assess impact and urgency
    impact := irs.assessImpact(incident)
    urgency := irs.assessUrgency(incident)
    
    // Calculate priority
    priority := irs.calculatePriority(impact, urgency)
    
    // Update incident record
    irs.incidentRepo.UpdatePriority(incidentID, priority)
    
    // Assign to appropriate team based on type and priority
    team := irs.determineResponseTeam(incident.Type, priority)
    irs.incidentRepo.AssignTeam(incidentID, team)
    
    // Notify team
    irs.notificationService.NotifyTeam(team, incident)
    
    // Initiate response based on incident type and priority
    if priority == "critical" || priority == "high" {
        irs.initiateAutomatedResponse(incident)
    }
    
    // Update incident status
    irs.incidentRepo.UpdateStatus(incidentID, "triaged")
}
```

**Process:**
- Automated assessment of incident severity and impact
- Assignment to appropriate security team
- Initiation of response playbooks
- Notification to relevant stakeholders

#### 3. Containment Phase

**Implementation:**
```go
// Add to incident_containment_service.go
func (ics *IncidentContainmentService) ContainSecurityThreat(incidentID string) {
    // Get incident details
    incident, _ := ics.incidentRepo.FindByID(incidentID)
    
    // Apply containment strategy based on threat type
    switch incident.ThreatType {
    case "account_compromise":
        // Lock affected accounts
        for _, userID := range incident.AffectedUsers {
            ics.userService.LockAccount(userID, "security_incident", 0) // 0 duration means until manually unlocked
        }
        
        // Invalidate all sessions for affected users
        for _, userID := range incident.AffectedUsers {
            ics.sessionService.InvalidateAllUserSessions(userID)
        }
        
    case "api_abuse":
        // Block offending IPs
        for _, ip := range incident.SourceIPs {
            ics.wafService.BlockIP(ip, "incident_response", 24*time.Hour)
        }
        
        // Apply stricter rate limits
        ics.rateLimitService.ApplyEmergencyLimits(incident.AffectedEndpoints)
        
    case "data_breach":
        // Isolate affected services
        for _, service := range incident.AffectedServices {
            ics.containerService.IsolateContainer(service)
        }
        
        // Revoke compromised credentials
        for _, credential := range incident.CompromisedCredentials {
            ics.credentialService.RevokeCredential(credential)
        }
    }
    
    // Update incident status
    ics.incidentRepo.UpdateStatus(incidentID, "contained")
    
    // Log containment actions
    ics.securityEventService.LogContainmentActions(incidentID, incident.ThreatType)
}
```

**Actions:**
- Isolation of affected systems
- Blocking of malicious IP addresses
- Suspension of compromised accounts
- Revocation of affected credentials
- Implementation of emergency access controls

#### 4. Eradication Phase

**Process:**
- Removal of malware or unauthorized access
- Patching of exploited vulnerabilities
- Strengthening of affected security controls
- Verification that threat has been eliminated

#### 5. Recovery Phase

**Process:**
- Restoration of affected systems to normal operation
- Verification of system integrity
- Gradual removal of emergency controls
- Monitoring for any signs of persistent threats

#### 6. Post-Incident Analysis

**Process:**
- Comprehensive review of the incident
- Documentation of lessons learned
- Updates to detection rules and response procedures
- Implementation of preventive measures

### Integration with Existing Infrastructure

#### 1. Monitoring Integration

- Configure Prometheus to collect security metrics
- Create dedicated security dashboards in Grafana
- Set up Alertmanager rules for security alerts
- Integrate with existing logging infrastructure

#### 2. CI/CD Pipeline Integration

- Add security scanning to the existing CodeQL workflow
- Implement dependency vulnerability scanning
- Add container image scanning before deployment
- Perform infrastructure-as-code security validation

#### 3. AWS Services Integration

- Configure AWS Security Hub as the central SIEM
- Enable AWS GuardDuty for threat detection
- Implement AWS WAF for web application protection
- Use AWS Config for compliance monitoring

### Implementation Roadmap

#### Phase 1: Foundation (1-2 months)
- Deploy SIEM solution and integrate with existing logs
- Implement basic WAF protection
- Enhance authentication monitoring
- Develop initial incident response procedures

#### Phase 2: Enhanced Detection (2-3 months)
- Implement RASP solution
- Develop fraud detection capabilities
- Create behavioral analytics baseline
- Expand WAF rules and protection

#### Phase 3: Automated Response (3-4 months)
- Implement automated containment procedures
- Develop dynamic rate limiting
- Create security playbooks for common threats
- Integrate threat intelligence feeds

#### Phase 4: Maturity and Optimization (4-6 months)
- Implement advanced analytics and ML-based detection
- Conduct red team exercises to test effectiveness
- Optimize detection rules and response procedures
- Implement continuous security improvement process

### Conclusion

The implementation of real-time threat detection and response capabilities will significantly enhance the security posture of the InvoiceB2B application. By leveraging the existing monitoring infrastructure and adding specialized security components, the system will be able to detect and respond to threats in real-time, minimizing the potential impact of security incidents.

The phased implementation approach allows for gradual deployment and testing, ensuring that each component is properly integrated and functioning as expected before moving to the next phase. Regular reviews and updates to detection rules and response procedures will ensure that the system remains effective against evolving threats.