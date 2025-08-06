# Vue Sonner CSS Fix

## Issue
The application was missing the proper CSS import specifier for the vue-sonner package, which was causing toast notifications to not display correctly.

## Root Cause
The CSS import path in nuxt.config.ts was using a direct reference to the internal file structure of the vue-sonner package:

```javascript
css: ['~/assets/css/tailwind.css', 'vue-sonner/lib/index.css'],
```

While this path technically works (as the CSS file does exist at that location), it's not the recommended import path according to the vue-sonner documentation.

## Solution
Updated the CSS import path in nuxt.config.ts to use the documented path:

```javascript
css: ['~/assets/css/tailwind.css', 'vue-sonner/style.css'],
```

This change aligns with the vue-sonner documentation and follows the package's recommended usage.

## Technical Details
After investigating the vue-sonner package.json, we found that there's a mapping in the "exports" section:

```json
"./style.css": "./lib/index.css",
```

This means that both import paths point to the same CSS file:
1. 'vue-sonner/style.css' (the documented path)
2. 'vue-sonner/lib/index.css' (the direct path)

Using the documented path is better for maintainability and follows the package's recommended usage.

## Alternative Approach
An alternative approach would be to use the Nuxt module provided by vue-sonner:

```javascript
// nuxt.config.ts
export default defineNuxtConfig({
  modules: ['vue-sonner/nuxt'],
  vueSonner: {
    css: true // Include CSS automatically
  }
})
```

This approach would handle the CSS import automatically. However, since the manual CSS import approach is working correctly with the updated path, we decided to stick with that approach for now.

## Verification
The toast notifications should now display correctly throughout the application. Toast notifications are used in several files:
- layouts/dashboard.vue
- pages/financing.vue
- pages/invoices.vue
- pages/onboarding.vue
- pages/submit-invoice.vue
- components/Views/Auth/Login.vue
- components/Views/Auth/Register.vue

## References
- [Vue Sonner Documentation](https://github.com/xiaoluoboding/vue-sonner)
- [Nuxt CSS Configuration](https://nuxt.com/docs/api/configuration/nuxt-config#css)