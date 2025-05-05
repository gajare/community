# community
Procore community extensions
## Accident Log Viewer – Authenticated Frontend with Go Backend

### Overview

This project provides a **user-friendly frontend UI** and a **Go backend API** for accessing various types of logs, including daily accident logs. It is designed to work with **Procore’s open-source APIs** and can be extended to support any custom logging requirements based on customer needs.

The frontend allows users to log in using credentials, retrieves an authentication token, and then uses that token to securely fetch and display logs. All logs are organized in a **single-page interface**, giving users a clean and centralized view of the data they need.

---

### Features

#### ✅ Frontend (React + Tailwind CSS or any preferred framework)
- Clean and modern single-page UI.
- Login form to enter credentials and get an **authentication token**.
- Display accident logs (and other types of logs) in a user-friendly table.
- Handles API responses gracefully with success/error indicators.
- Easily extendable for additional log types as per customer requirement.

#### ⚙️ Backend (Go/Golang)
- Written in Go using standard net/http and Gorilla Mux.
- RESTful APIs for:
  - **Authentication** – Validate user credentials and return JWT token.
  - **Log Retrieval** – Secure endpoints to fetch accident logs.
- Token-based authentication to protect sensitive endpoints.
- Modular design, making it easy to plug in additional APIs for different types of logs.

---

### Architecture

```text
+------------+       +----------------------+       +-------------------+
|            |       |                      |       |                   |
|  Frontend  +------>+  Go Backend (API)    +------>+  Procore API /    |
| (SPA UI)   |       |                      |       |  Custom Log Source|
+------------+       +----------------------+       +-------------------+
