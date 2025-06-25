1. Project Overview
   Develop an invoice financing platform that connects businesses seeking capital with financial institutions. The platform will facilitate the secure upload of invoices, automated data extraction, presentation of funding offers, and real-time tracking of financing progress. An robust admin interface will allow for business onboarding management, financial institution oversight, and activity monitoring.

2. Technology Stack
   Backend Server: Nuxt.js (acting as an API server, handling routes, business logic, and database interactions).

Database: PostgreSQL (for relational data storage).

Message Broker: RabbitMQ (for asynchronous communication and event-driven architecture).

Frontend: Nuxt.js (for the user interface, interacting with the Nuxt.js backend API).

External Workflow Tool: n8n (for invoice data extraction, integrated via RabbitMQ).

3. Core User Roles & Functionality
   3.1. Business User Features
   Secure Authentication: User registration, login, forgot password, and new password flows with email-based OTP/verification.

Invoice Upload:

Ability to upload invoices in PDF, Excel, or common image formats (e.g., JPG, PNG).

Upon successful upload, an event is sent to RabbitMQ (e.g., invoice_uploaded queue) to trigger the n8n workflow for data extraction.

Immediate email confirmation to the business upon successful upload.

Invoice Data Population & Funding Options:

Once n8n extracts data and sends it back (via RabbitMQ, e.g., invoice_extracted queue), the extracted invoice details should be stored and presented to the business.

The system displays a curated list of financial institutions (FIs) that are willing to fund the uploaded invoice, along with their respective loan terms (e.g., discount rate, repayment terms).

Financial Institution Selection:

Businesses can select a preferred financial institution from the displayed options.

Upon selection, a new event is sent to RabbitMQ (e.g., fi_selected queue) to trigger the notification process.

An automated email is sent to the selected financial institution containing:

Business information (from KYC data).

Detailed invoice information.

Relevant contact details for the business.

An email confirmation is sent to the business regarding their selection.

Invoice Tracking:

A dashboard/section where businesses can track the real-time funding status of each uploaded invoice (e.g., "Uploaded," "Processing," "Awaiting FI Selection," "FI Selected," "Funded," "Rejected").

Status updates should ideally be event-driven via RabbitMQ, ensuring near real-time visibility.

For "Funded" invoices, businesses can view the funded details and download the transaction receipt.

KYC & Business Information: A section for businesses to upload and manage their KYC documents and basic business profile information, which will be accessible to selected FIs.

Email Notifications: Implement comprehensive email notifications for all critical actions:

Registration success

Forgot password / New password

OTP delivery

Invoice upload confirmation

Invoice processed/data extracted

FI selection confirmation

Invoice status changes (funded, rejected, etc.)

Invoice Funded (including a link to view details/download receipt).

3.2. Admin User Features
Admin Dashboard: A centralized, secure interface for managing the entire platform.

Activity Tracking:

Comprehensive logs and real-time tracking of all critical user activities (business logins, invoice uploads, FI selections, status changes).

Track when an invoice is funded by an FI and when transaction details are uploaded.

Ability to filter and search activities.

Business Onboarding & Management:

View a list of all registered businesses.

Manual Approval System: Businesses must be manually approved by an admin before they can start uploading invoices or engaging with financial institutions.

Ability to view, edit, and manage business profiles and their uploaded KYC documents.

Option to activate/deactivate business accounts.

Financial Institution Management:

CRUD (Create, Read, Update, Delete) functionality for financial institutions.

For each FI, admins can define:

Name, contact details, logo.

Default loan terms (e.g., minimum/maximum financing amount, interest rates/discount rates, repayment periods).

Terms and Conditions specific to that FI.

Ability to link FIs to specific industries or business types (optional, but good for scalability).

Invoice Management:

View all uploaded invoices, their current status, and associated business/FI.

Ability to manually update invoice statuses if automated processes fail or require override.

Manual trigger for FI notification/email sending.

View transaction information and receipt for funded invoices, and verify funding if needed.

Automated Functions with Manual Override:

Automate notifications and status updates where possible (e.g., based on RabbitMQ events).

Provide manual controls for critical actions, such as manually approving a business, changing an invoice status, or re-sending an email notification.

User Management: Admin can manage other admin accounts (create, edit, delete, assign roles/permissions).

4. Database Schema Considerations (PostgreSQL)
   Design a robust PostgreSQL schema including, but not limited to, the following tables:

users: Stores id, email, password_hash, role (business/admin), is_approved, created_at, updated_at.

businesses: id, user_id (FK to users), name, registration_number, contact_person, phone, address, kyc_status, kyc_documents (JSONB for file paths/metadata), created_at, updated_at.

financial_institutions: id, name, email, contact_person, phone, address, loan_terms (JSONB for structured terms), terms_conditions (TEXT), created_at, updated_at.

invoices: id, business_id (FK to businesses), invoice_number, amount, currency, due_date, uploaded_document_path, extracted_data (JSONB), current_status (enum/text), selected_fi_id (FK to financial_institutions, nullable), created_at, updated_at, funded_amount (NUMERIC, nullable), funded_date (TIMESTAMP, nullable), transaction_id (TEXT, nullable), receipt_path (TEXT, nullable).

invoice_status_history: id, invoice_id (FK), status_change_timestamp, new_status, changed_by (user_id/system), notes.

activity_logs: id, user_id (nullable), activity_type, description, timestamp, ip_address.

otp_codes: id, user_id, code, expires_at, type (e.g., 'registration', 'password_reset').

5. Backend Logic (Nuxt.js Server)
   API Endpoints:

Authentication: /api/auth/register, /api/auth/login, /api/auth/forgot-password, /api/auth/reset-password, /api/auth/verify-otp.

Business: /api/business/profile, /api/business/invoices (upload, list, detail, select FI), /api/business/kyc.

Admin: /api/admin/businesses (list, approve, update), /api/admin/financial-institutions (CRUD), /api/admin/invoices (list, update status), /api/admin/activities.

RabbitMQ Integration:

Producers: Nuxt.js server publishes messages to RabbitMQ on events like invoice upload, FI selection.

Consumers: Nuxt.js server consumes messages from RabbitMQ for events like invoice data extracted from n8n, or status updates triggered by external systems/n8n.

Email Service Integration: Connect with an email service provider (e.g., SendGrid, Nodemailer) to send all transactional emails.

Security: Implement robust authentication (JWT recommended) and authorization (role-based access control) for all API endpoints. Input validation and data sanitization are crucial.

File Storage: Implement a secure way to store uploaded invoices and KYC documents (e.g., local storage, S3, Google Cloud Storage, with file paths stored in PostgreSQL).

6. Frontend (Nuxt.js)
   Develop a clean, intuitive, and responsive user interface for both business and admin users.

Utilize Nuxt.js's capabilities for component-based UI, routing, and state management.

Ensure real-time updates for invoice tracking where applicable (e.g., using WebSockets or frequent polling, potentially integrated with RabbitMQ).

7. RabbitMQ Event Architecture Examples
   invoice_uploaded queue:

Producer: Nuxt.js (when a business uploads an invoice).

Consumer: n8n workflow (triggers data extraction).

invoice_extracted queue:

Producer: n8n workflow (after successful data extraction).

Consumer: Nuxt.js (to update invoice record in DB and present to business).

fi_selected queue:

Producer: Nuxt.js (when a business selects an FI).

Consumer: Nuxt.js (to send email to FI, update invoice status).

invoice_funded_by_fi queue:

Producer: n8n workflow (triggered by FI email/webhook indicating funding, includes transaction data).

Consumer: Nuxt.js (to update invoice status to "Funded", store funded_amount, funded_date, transaction_id, receipt_path in the database, send email to business, notify admin).

send_email queue:

Producer: Nuxt.js (for all types of transactional emails: registration, password reset, notifications).

Consumer: Nuxt.js (dedicated worker to handle email sending asynchronously).

invoice_status_update queue:

Producer: Nuxt.js (admin manual update), n8n (if n8n handles further funding stages).

Consumer: Nuxt.js (to update UI for businesses).

8. Scalability & Error Handling
   Design the system with scalability in mind, especially considering the asynchronous nature provided by RabbitMQ.

Implement comprehensive error handling, logging, and monitoring across all layers (frontend, backend, database, RabbitMQ consumers/producers).

Consider retry mechanisms for failed RabbitMQ messages.

## Setup

Make sure to install dependencies:

```bash
# npm
npm install

# pnpm
pnpm install

# yarn
yarn install

# bun
bun install
```

## Development Server

Start the development server on `http://localhost:3000`:

```bash
# npm
npm run dev

# pnpm
pnpm dev

# yarn
yarn dev

# bun
bun run dev
```

## Production

Build the application for production:

```bash
# npm
npm run build

# pnpm
pnpm build

# yarn
yarn build

# bun
bun run build
```

Locally preview production build:

```bash
# npm
npm run preview

# pnpm
pnpm preview

# yarn
yarn preview

# bun
bun run preview
```

Check out the [deployment documentation](https://nuxt.com/docs/getting-started/deployment) for more information.
