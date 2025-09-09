# CSRF Protection Implementation Guide

This document provides guidance on how to integrate with the CSRF protection mechanism implemented in the InvoiceB2B API.

## Overview

Cross-Site Request Forgery (CSRF) is an attack that forces authenticated users to execute unwanted actions on a web application in which they're currently authenticated. The InvoiceB2B API implements token-based CSRF protection using the double-submit cookie pattern to prevent these attacks.

## How It Works

1. When a user accesses any API endpoint, the server generates a CSRF token and sends it to the client in two ways:
   - As an HTTP-only cookie named `X-CSRF-Token`
   - As a response header also named `X-CSRF-Token`

2. For state-changing operations (POST, PUT, DELETE), the client must include the token from the response header in a request header with the same name (`X-CSRF-Token`).

3. The server validates that the token in the request header matches the token in the cookie, and that the signature is valid.

## Frontend Integration

### JavaScript/AJAX Requests

For JavaScript applications making AJAX requests, you need to:

1. Extract the CSRF token from the response headers when making any request to the API
2. Store this token in your application state
3. Include the token in the headers of all subsequent state-changing requests

Example using Axios:

```javascript
// Configure Axios to automatically handle CSRF tokens
axios.interceptors.response.use(response => {
  // Extract CSRF token from response headers if present
  const csrfToken = response.headers['x-csrf-token'];
  if (csrfToken) {
    // Store the token for future requests
    localStorage.setItem('csrfToken', csrfToken);
  }
  return response;
});

axios.interceptors.request.use(config => {
  // Add CSRF token to request headers for state-changing methods
  if (['post', 'put', 'delete'].includes(config.method.toLowerCase())) {
    const csrfToken = localStorage.getItem('csrfToken');
    if (csrfToken) {
      config.headers['X-CSRF-Token'] = csrfToken;
    }
  }
  return config;
});
```

### React Example

```jsx
import React, { useState, useEffect } from 'react';
import axios from 'axios';

// Configure axios as shown above

function ProfileForm() {
  const [profile, setProfile] = useState({});
  
  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      // The CSRF token will be automatically included by the axios interceptor
      await axios.put('/api/user/profile', profile);
      alert('Profile updated successfully');
    } catch (error) {
      console.error('Error updating profile:', error);
    }
  };
  
  return (
    <form onSubmit={handleSubmit}>
      {/* Form fields */}
      <button type="submit">Update Profile</button>
    </form>
  );
}
```

### Vue.js Example

```javascript
// In your main.js or a plugin file
import axios from 'axios';

// Configure axios interceptors as shown above

// Then in your component
export default {
  methods: {
    async updateProfile() {
      try {
        // The CSRF token will be automatically included by the axios interceptor
        await axios.put('/api/user/profile', this.profile);
        this.$notify({
          type: 'success',
          message: 'Profile updated successfully'
        });
      } catch (error) {
        console.error('Error updating profile:', error);
      }
    }
  }
}
```

### Angular Example

```typescript
// In an HTTP interceptor
import { Injectable } from '@angular/core';
import { HttpInterceptor, HttpRequest, HttpHandler, HttpEvent, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators';

@Injectable()
export class CSRFInterceptor implements HttpInterceptor {
  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    // Add CSRF token to state-changing requests
    if (req.method !== 'GET' && req.method !== 'HEAD' && req.method !== 'OPTIONS') {
      const csrfToken = localStorage.getItem('csrfToken');
      if (csrfToken) {
        req = req.clone({
          setHeaders: {
            'X-CSRF-Token': csrfToken
          }
        });
      }
    }
    
    return next.handle(req).pipe(
      tap(event => {
        if (event instanceof HttpResponse) {
          // Extract and store CSRF token from response
          const csrfToken = event.headers.get('x-csrf-token');
          if (csrfToken) {
            localStorage.setItem('csrfToken', csrfToken);
          }
        }
      })
    );
  }
}
```

## Mobile App Integration

For mobile applications, the same principles apply:

1. Extract the CSRF token from response headers
2. Store it securely in the app's state
3. Include it in headers for state-changing requests

### iOS (Swift) Example

```swift
import Foundation

class APIClient {
    static let shared = APIClient()
    private var csrfToken: String?
    
    func request<T: Decodable>(endpoint: String, method: String, body: Data? = nil, completion: @escaping (Result<T, Error>) -> Void) {
        guard let url = URL(string: "https://api.invoiceb2b.com\(endpoint)") else {
            completion(.failure(NSError(domain: "Invalid URL", code: 0, userInfo: nil)))
            return
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = method
        
        // Add CSRF token for state-changing methods
        if method == "POST" || method == "PUT" || method == "DELETE", let token = csrfToken {
            request.addValue(token, forHTTPHeaderField: "X-CSRF-Token")
        }
        
        // Add other headers as needed
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")
        
        if let body = body {
            request.httpBody = body
        }
        
        let task = URLSession.shared.dataTask(with: request) { [weak self] data, response, error in
            // Extract CSRF token from response headers
            if let httpResponse = response as? HTTPURLResponse,
               let token = httpResponse.allHeaderFields["X-CSRF-Token"] as? String {
                self?.csrfToken = token
            }
            
            // Handle response...
        }
        
        task.resume()
    }
}
```

### Android (Kotlin) Example

```kotlin
import okhttp3.OkHttpClient
import okhttp3.Request
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory

class CSRFTokenInterceptor : Interceptor {
    private var csrfToken: String? = null

    override fun intercept(chain: Interceptor.Chain): Response {
        val original = chain.request()
        
        // Add CSRF token for state-changing methods
        val request = if (original.method in listOf("POST", "PUT", "DELETE") && csrfToken != null) {
            original.newBuilder()
                .header("X-CSRF-Token", csrfToken!!)
                .build()
        } else {
            original
        }
        
        val response = chain.proceed(request)
        
        // Extract CSRF token from response headers
        response.header("X-CSRF-Token")?.let {
            csrfToken = it
        }
        
        return response
    }
}

// Setup Retrofit with the interceptor
val client = OkHttpClient.Builder()
    .addInterceptor(CSRFTokenInterceptor())
    .build()

val retrofit = Retrofit.Builder()
    .baseUrl("https://api.invoiceb2b.com/")
    .client(client)
    .addConverterFactory(GsonConverterFactory.create())
    .build()
```

## Troubleshooting

### Common Issues

1. **CSRF Token Missing**: If you receive a 403 Forbidden error with a message about a missing CSRF token, ensure you're extracting the token from the response headers and including it in your request headers.

2. **CSRF Token Mismatch**: If you receive a 403 Forbidden error with a message about a token mismatch, it could be due to:
   - Using an expired token (tokens expire after 24 hours)
   - Using a token from a different session
   - Cookie issues (ensure cookies are being sent with your requests)

3. **Cookie Issues**: The CSRF protection relies on cookies. Ensure that:
   - Your frontend is configured to send cookies with requests (withCredentials: true for AJAX)
   - You're not in a cross-origin situation without proper CORS configuration

### Testing CSRF Protection

You can verify that CSRF protection is working by:

1. Making a GET request to any endpoint to obtain a CSRF token
2. Attempting a POST/PUT/DELETE request without the token (should fail with 403)
3. Attempting the same request with the token (should succeed)

## Security Considerations

- Never store the CSRF token in cookies that are accessible via JavaScript
- For SPAs, storing the token in memory is preferable to localStorage when possible
- The token is tied to the user's session, so it will change when the user logs out and back in
- The implementation uses cryptographic signatures to prevent token tampering

## API Reference

### Headers

- `X-CSRF-Token` (response): Contains the CSRF token that should be used for subsequent requests
- `X-CSRF-Token` (request): Should contain the CSRF token for state-changing operations

### Status Codes

- `403 Forbidden`: Returned when CSRF validation fails, with a message explaining the reason

## Conclusion

By following this guide, your frontend application will be properly integrated with the CSRF protection mechanism in the InvoiceB2B API, helping to prevent CSRF attacks and enhance the security of your application.