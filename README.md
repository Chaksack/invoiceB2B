# Invoice Financing Application: Plan & Architecture

**Version:** 1.0
**Date:** May 21, 2025

## 1. Overview

This document outlines the plan and architecture for an Invoice Financing Application. The platform will allow businesses (users) to upload invoices for financing. Administrators will review these invoices, manage user KYC, and process financing. The system will include features for user registration, 2FA, invoice processing (PDF/CSV to JSON), email notifications, fund disbursement tracking, repayment handling, and activity logging.

## 2. Core Features

### 2.1. User Features:

* **Registration:**
   * Input: Email, First Name, Last Name, Company Name, Password.
   * Action: Create user account, send registration confirmation email.
* **Login:**
   * Input: Email, Password.
   * Action: Authenticate user.
   * **Two-Factor Authentication (2FA):** Required after successful password authentication (e.g., via authenticator app or email OTP).
* **KYC (Know Your Customer):**
   * User submits required KYC documents/information.
   * Status: Pending, Approved, Rejected.
   * Users cannot upload invoices until KYC is approved.
* **Invoice Management:**
   * **Upload:**
      * Formats: PDF, CSV.
      * Action: Convert uploaded file to JSON format for internal processing. Store original file.
      * Validation: Basic checks on file format and content.
   * **View Invoices:** List of uploaded invoices with their status.
   * **Track Status:** Real-time updates on invoice processing (e.g., Pending Review, Approved, Disbursed, Repaid, Rejected).
* **Repayment:**
   * User initiates repayment for a financed invoice.
   * Record repayment details.
* **Notifications:**
   * In-app notifications.
   * Email notifications for:
      * Registration confirmation.
      * KYC status updates (Submitted, Approved, Rejected).
      * Invoice status updates (Submitted, Approved, Disbursed, Repaid, Rejected).
      * Disbursement receipt.
      * Repayment confirmation.
* **Receipts:**
   * View and download disbursement receipts (PDF).

### 2.2. Admin Features:

* **Dashboard:** Overview of pending KYCs, pending invoices, total financed amount, etc.
* **User Management:**
   * View registered users.
   * **KYC Approval:** Review submitted KYC information and approve/reject users.
* **Invoice Management:**
   * View all uploaded invoices.
   * Inspect invoice details (original file and extracted JSON data).
   * **Financing Decision:** Approve or reject invoices for financing.
   * **Mark as Disbursed:** Update invoice status once funds are confirmed as transferred.
   * **Mark as Repaid:** Confirm user repayments.
* **Staff Management:**
   * Create staff accounts with specific roles/permissions (e.g., KYC reviewer, finance approver).
   * Manage staff accounts.
* **Activity Logging:**
   * Track key actions:
      * Who approved/rejected KYC.
      * Who approved/rejected invoices.
      * Who marked invoices as disbursed.
      * Staff login history.
* **Notifications:**
   * Alerts for new KYC submissions or invoice uploads.

## 3. Technology Stack

* **Backend Framework:** GoFiber (High-performance Go web framework)
* **ORM:** GORM (Developer-friendly ORM for Go)
* **Database:** PostgreSQL (Reliable relational database)
* **Caching/Session/2FA:** Redis (In-memory data store)
* **Message Queue:** RabbitMQ (For asynchronous tasks like email sending, PDF/CSV processing)
* **Containerization:** Docker & Docker Compose (For development and deployment consistency)
* **Authentication:** JWT (JSON Web Tokens) for API authentication.
* **Session Management:** Secure, HTTP-only cookies with JWT.
* **PDF/CSV Processing:** Go libraries (e.g., `unidoc/unipdf` for PDF, standard `encoding/csv` for CSV).
* **Email Sending:** Go SMTP library or a third-party email service API (e.g., SendGrid, Mailgun).

## 4. System Architecture (High-Level)

+-------------------+      +-------------------+      +---------------------+|      User UI      |----->|    GoFiber API    |<---->|     PostgreSQL      || (Web Application) |      |    (Backend)      |      |      (Database)     |+-------------------+      +-------------------+      +----------+----------+|        ^                            ||        |                            ||        v                            ||  +-----------+                      ||  |  RabbitMQ | (Async Tasks)        ||  +-----------+                      ||     /  |  \                         ||    /   |   \                        |v   v    v    v                       |(Email | PDF | CSV | Other Workers)         |Service| Proc| Proc|                        ||        ^                            ||        |                            ||        v                            ||  +-----------+                      |+->|   Redis   |<---------------------++-----------+(Cache, Sessions, 2FA OTPs)
**Components:**

1.  **User Interface (Web Application):** A frontend application (e.g., built with React, Vue, or server-side Go templates) that users and admins interact with.
2.  **GoFiber API Backend:**
   * Handles all business logic, API requests, and responses.
   * Manages authentication (JWT) and authorization.
   * Interacts with PostgreSQL via GORM.
   * Uses Redis for caching, session storage, and 2FA token management.
   * Publishes messages to RabbitMQ for asynchronous tasks.
3.  **PostgreSQL Database:**
   * Stores persistent data: users, invoices, KYC details, transactions, staff accounts, activity logs.
4.  **Redis:**
   * **Caching:** Frequently accessed data (e.g., user profiles, configuration).
   * **Session Management:** Store user session information linked to secure cookies.
   * **2FA:** Store temporary OTPs for two-factor authentication.
5.  **RabbitMQ:**
   * **Message Broker:** Decouples services and handles asynchronous operations.
   * **Workers:** Separate Go services that consume messages from RabbitMQ:
      * **Email Worker:** Sends emails (registration, notifications, receipts).
      * **PDF Processing Worker:** Parses uploaded PDF invoices, extracts data to JSON.
      * **CSV Processing Worker:** Parses uploaded CSV invoices, extracts data to JSON.
      * Other background tasks as needed.
6.  **Docker:**
   * Containerizes each service (GoFiber API, PostgreSQL, Redis, RabbitMQ, Workers) for consistent environments and easier deployment.
   * `docker-compose.yml` will define and manage the multi-container application.

## 5. Database Schema (High-Level)

* **`users`**
   * `id` (PK, UUID/Serial)
   * `email` (UNIQUE, VARCHAR)
   * `first_name` (VARCHAR)
   * `last_name` (VARCHAR)
   * `company_name` (VARCHAR)
   * `password_hash` (VARCHAR)
   * `kyc_id` (FK to `kyc_details`)
   * `is_active` (BOOLEAN, default: true)
   * `is_admin` (BOOLEAN, default: false) - *Consider a separate `staff` table or roles system for more granularity.*
   * `two_fa_secret` (VARCHAR, nullable) - For TOTP
   * `two_fa_enabled` (BOOLEAN, default: false)
   * `created_at`, `updated_at`

* **`kyc_details`**
   * `id` (PK, UUID/Serial)
   * `user_id` (FK to `users`, UNIQUE)
   * `status` (ENUM: 'pending', 'approved', 'rejected', 'resubmit_required')
   * `submitted_at` (TIMESTAMP)
   * `reviewed_by` (FK to `staff`, nullable)
   * `reviewed_at` (TIMESTAMP, nullable)
   * `rejection_reason` (TEXT, nullable)
   * `documents_info` (JSONB) - Store metadata about uploaded documents (paths, types)
   * `created_at`, `updated_at`

* **`invoices`**
   * `id` (PK, UUID/Serial)
   * `user_id` (FK to `users`)
   * `invoice_number` (VARCHAR)
   * `issuer_name` (VARCHAR) - Extracted
   * `issuer_bank_account` (VARCHAR) - Extracted
   * `issuer_bank_name` (VARCHAR) - Extracted
   * `debtor_name` (VARCHAR) - Extracted
   * `amount` (DECIMAL)
   * `currency` (VARCHAR(3))
   * `due_date` (DATE)
   * `status` (ENUM: 'pending_review', 'approved', 'rejected', 'disbursed', 'repayment_pending', 'repaid')
   * `original_file_path` (VARCHAR)
   * `json_data` (JSONB) - Extracted data
   * `uploaded_at` (TIMESTAMP)
   * `approved_by` (FK to `staff`, nullable)
   * `approved_at` (TIMESTAMP, nullable)
   * `disbursed_by` (FK to `staff`, nullable)
   * `disbursed_at` (TIMESTAMP, nullable)
   * `financing_fee_percentage` (DECIMAL, nullable)
   * `financed_amount` (DECIMAL, nullable)
   * `disbursement_receipt_path` (VARCHAR, nullable)
   * `created_at`, `updated_at`

* **`transactions`** (For tracking fund movements, repayments)
   * `id` (PK, UUID/Serial)
   * `invoice_id` (FK to `invoices`)
   * `type` (ENUM: 'disbursement', 'repayment')
   * `amount` (DECIMAL)
   * `transaction_date` (TIMESTAMP)
   * `reference_id` (VARCHAR, nullable) - e.g., bank transaction ID
   * `created_at`, `updated_at`

* **`staff`** (Admin and other privileged users)
   * `id` (PK, UUID/Serial)
   * `email` (UNIQUE, VARCHAR)
   * `first_name` (VARCHAR)
   * `last_name` (VARCHAR)
   * `password_hash` (VARCHAR)
   * `role` (VARCHAR, e.g., 'admin', 'kyc_reviewer', 'finance_manager')
   * `is_active` (BOOLEAN, default: true)
   * `last_login_at` (TIMESTAMP, nullable)
   * `created_at`, `updated_at`

* **`activity_logs`**
   * `id` (PK, UUID/Serial)
   * `staff_id` (FK to `staff`, nullable - system actions might not have a `staff_id`)
   * `user_id` (FK to `users`, nullable - for user-related actions)
   * `action` (VARCHAR, e.g., 'USER_REGISTERED', 'KYC_APPROVED', 'INVOICE_UPLOADED', 'INVOICE_FINANCED', 'FUNDS_DISBURSED', 'STAFF_LOGIN')
   * `details` (JSONB, nullable) - Additional context
   * `ip_address` (VARCHAR, nullable)
   * `timestamp` (TIMESTAMP, default: NOW())

* **`sessions`** (If using Redis for sessions, this might not be a DB table but a Redis structure)
   * `session_id` (PK, VARCHAR)
   * `user_id` (FK to `users` or `staff`)
   * `data` (JSONB or TEXT)
   * `expires_at` (TIMESTAMP)

* **`two_fa_tokens`** (For storing temporary 2FA codes in Redis, structure would be key-value like `user_id:otp`)

## 6. Key Workflows & Processes

### 6.1. User Registration & KYC:
1.  User submits registration form.
2.  API validates data, creates `users` record (inactive or KYC pending).
3.  API publishes 'user_registered' event to RabbitMQ.
4.  Email Worker consumes event, sends confirmation email.
5.  User logs in, prompted to complete KYC.
6.  User submits KYC info/documents.
7.  API creates `kyc_details` record (status: 'pending'), stores documents (e.g., S3 or local secure storage, path in DB).
8.  Admin notified of new KYC submission.
9.  Admin reviews KYC via admin panel.
10. Admin approves/rejects KYC. API updates `kyc_details` status.
11. API publishes 'kyc_status_updated' event to RabbitMQ.
12. Email Worker consumes event, sends KYC status email to user.
13. If approved, user account is fully activated.

### 6.2. User Login with 2FA:
1.  User submits email/password.
2.  API validates credentials.
3.  If valid and 2FA is enabled for the user:
   * Generate OTP (Time-based One-Time Password - TOTP, or send via email/SMS if preferred).
   * Store OTP temporarily in Redis with an expiry (e.g., `user_id:otp_value`).
   * Prompt user for OTP.
4.  User submits OTP.
5.  API validates OTP against Redis store.
6.  If valid:
   * Generate JWT.
   * Set secure, HTTP-only cookie containing JWT.
   * Create session entry in Redis (if applicable).
   * Grant access.
7.  If invalid OTP, deny access.

### 6.3. Invoice Upload & Processing:
1.  User (KYC approved) uploads PDF/CSV invoice.
2.  API receives file, performs basic validation.
3.  API stores original file securely (e.g., S3 or local volume).
4.  API creates `invoices` record (status: 'pending_review', `original_file_path` set).
5.  API publishes 'invoice_uploaded' event to RabbitMQ with file path and invoice ID.
6.  **PDF/CSV Worker:**
   * Consumes message.
   * Retrieves file.
   * Parses content (using appropriate libraries).
   * Extracts key fields (invoice number, amounts, dates, bank details).
   * Updates the `invoices` record with extracted `json_data`.
   * Handles parsing errors and updates invoice status if unparseable.
7.  Admin notified of new invoice for review.

### 6.4. Invoice Financing Workflow:
1.  **Admin Review:**
   * Admin views pending invoices in the admin panel.
   * Inspects original file and extracted JSON data.
2.  **Admin Decision (Approve/Reject):**
   * Admin approves or rejects the invoice.
   * API updates `invoices.status`, `approved_by`, `approved_at`.
   * API publishes 'invoice_status_updated' event (for 'approved' or 'rejected').
   * Email Worker sends notification to user.
   * Activity log created.
3.  **Funds Disbursement (Manual/Simulated):**
   * If approved, finance team processes fund transfer externally (initially).
   * Admin marks invoice as 'disbursed' in the system.
   * API updates `invoices.status`, `disbursed_by`, `disbursed_at`, `financed_amount`.
   * **Receipt Generation:** System generates a disbursement receipt (PDF).
      * Store receipt path in `invoices.disbursement_receipt_path`.
   * API publishes 'invoice_disbursed' event.
   * Email Worker sends 'disbursed' notification and receipt to user.
   * Activity log created.
4.  **User Repayment:**
   * User initiates repayment through the platform (details of payment method TBD).
   * API records repayment attempt.
   * Admin/System verifies repayment.
   * Admin marks invoice as 'repaid'.
   * API updates `invoices.status`.
   * API publishes 'invoice_repaid' event.
   * Email Worker sends repayment confirmation to user.
   * Activity log created.

## 7. API Endpoints (Illustrative - GoFiber)

* **Auth:**
   * `POST /api/v1/auth/register`
   * `POST /api/v1/auth/login`
   * `POST /api/v1/auth/login/2fa/verify`
   * `POST /api/v1/auth/login/2fa/setup` (to get secret for authenticator app)
   * `POST /api/v1/auth/logout`
   * `POST /api/v1/auth/refresh-token`
* **User:**
   * `GET /api/v1/user/profile`
   * `PUT /api/v1/user/profile`
   * `POST /api/v1/user/kyc` (submit KYC)
   * `GET /api/v1/user/kyc` (get KYC status)
* **Invoices (User):**
   * `POST /api/v1/invoices` (upload)
   * `GET /api/v1/invoices` (list user's invoices)
   * `GET /api/v1/invoices/:id` (get specific invoice)
   * `GET /api/v1/invoices/:id/receipt` (download receipt)
   * `POST /api/v1/invoices/:id/repay`
* **Admin - Users:**
   * `GET /api/v1/admin/users`
   * `GET /api/v1/admin/users/:id/kyc`
   * `PUT /api/v1/admin/users/:id/kyc/approve`
   * `PUT /api/v1/admin/users/:id/kyc/reject`
* **Admin - Invoices:**
   * `GET /api/v1/admin/invoices` (list all invoices, with filters)
   * `GET /api/v1/admin/invoices/:id`
   * `PUT /api/v1/admin/invoices/:id/approve`
   * `PUT /api/v1/admin/invoices/:id/reject`
   * `PUT /api/v1/admin/invoices/:id/disburse`
   * `PUT /api/v1/admin/invoices/:id/confirm-repayment`
* **Admin - Staff:**
   * `POST /api/v1/admin/staff`
   * `GET /api/v1/admin/staff`
   * `PUT /api/v1/admin/staff/:id`
   * `DELETE /api/v1/admin/staff/:id`
* **Admin - Activity Logs:**
   * `GET /api/v1/admin/activity-logs`

## 8. Security Considerations

* **Password Hashing:** Use a strong hashing algorithm (e.g., bcrypt or Argon2).
* **JWT Security:**
   * Short-lived access tokens, longer-lived refresh tokens.
   * Store JWT in secure, HTTP-only cookies.
   * HTTPS enforced for all communication.
* **2FA:** Protects against compromised credentials.
* **Input Validation:** Rigorous validation on all incoming data (API and file uploads).
* **SQL Injection Prevention:** Use GORM which handles prepared statements.
* **XSS Prevention:** Proper output encoding in the frontend. GoFiber templates also offer protection.
* **CSRF Protection:** Use CSRF tokens if session cookies are primary auth for state-changing requests not via JWT bearer tokens. GoFiber has middleware for this.
* **Rate Limiting:** Protect against brute-force attacks on login and other sensitive endpoints.
* **Secure File Uploads:**
   * Validate file types and sizes.
   * Scan files for malware if possible (external service).
   * Store uploaded files in a non-web-accessible location or secure cloud storage (e.g., S3 with restricted access).
* **Permissions & Roles:** Granular access control for admin/staff functionalities.
* **Data Encryption:**
   * **In Transit:** TLS/SSL (HTTPS).
   * **At Rest:** Consider encrypting sensitive fields in the database (e.g., bank account details, KYC document info) using application-level encryption or database features.
* **Regular Security Audits & Updates:** Keep dependencies updated.

## 9. Email Notifications Summary

1.  **User Registration:** Welcome email with confirmation link/info.
2.  **KYC Submitted:** Confirmation that KYC info was received.
3.  **KYC Approved:** Notification of approval, user can now use full features.
4.  **KYC Rejected/Resubmit:** Notification with reason and next steps.
5.  **Invoice Uploaded:** Confirmation that invoice was received and is being processed.
6.  **Invoice Approved for Financing:** Notification that invoice is approved.
7.  **Invoice Rejected for Financing:** Notification with reason.
8.  **Invoice Funds Disbursed:** Notification that funds have been sent, with attached receipt.
9.  **Invoice Repayment Due Reminder:** (Optional, future feature)
10. **Invoice Repayment Confirmed:** Notification that repayment was successful.
11. **Password Reset Request:** Email with reset link.
12. **2FA Code (if email OTP is used):** Email with OTP.

## 10. Project Structure (GoFiber - Example)

/invoice-financing-app|-- /cmd|   |-- /api                 # Main API application|   |   |-- main.go|   |-- /pdfworker           # PDF processing worker|   |   |-- main.go|   |-- /csvworker           # CSV processing worker|   |   |-- main.go|   |-- /emailworker         # Email sending worker|   |   |-- main.go|-- /internal|   |-- /auth                # Authentication logic, JWT, 2FA|   |-- /config              # Configuration loading|   |-- /database            # Database connection, GORM setup|   |-- /handlers            # HTTP handlers (Fiber controllers)|   |   |-- auth_handler.go|   |   |-- user_handler.go|   |   |-- invoice_handler.go|   |   |-- admin_handler.go|   |-- /middleware          # Custom Fiber middleware|   |-- /models              # GORM models (User, Invoice, etc.)|   |-- /repositories        # Data access logic (interacts with GORM)|   |-- /services            # Business logic services|   |   |-- user_service.go|   |   |-- invoice_service.go|   |   |-- kyc_service.go|   |   |-- notification_service.go (interfaces with RabbitMQ)|   |-- /utils               # Utility functions (e.g., password hashing, file processing)|   |-- /workers             # Common worker logic, RabbitMQ interaction|-- /pkg                     # Shared libraries (if any, less common in typical Go structure)|-- /migrations              # Database migration files|-- /templates               # HTML templates (if using server-side rendering)|-- /uploads                 # Temporary storage for uploads (ensure proper security) - better to use S3|-- go.mod|-- go.sum|-- Dockerfile               # For the API|-- Dockerfile.pdfworker|-- Dockerfile.csvworker|-- Dockerfile.emailworker|-- docker-compose.yml|-- .env.example             # Environment variables template|-- README.md
## 11. Next Steps & Considerations

* **Frontend Development:** Choose a frontend framework or decide on server-side templates.
* **Detailed API Design:** Flesh out request/response DTOs for each endpoint.
* **Payment Gateway Integration:** For actual fund transfers and repayments (beyond the scope of this initial plan).
* **Cloud Deployment Strategy:** AWS, GCP, Azure, etc.
* **Logging and Monitoring:** Centralized logging (ELK stack, Grafana Loki) and application performance monitoring (Prometheus, Grafana).
* **Testing:** Unit tests, integration tests, end-to-end tests.

This plan provides a comprehensive starting point. Each section can be expanded further as the project progresses.
