# Vue Sonner Implementation

## Overview

This document outlines the implementation of toast notifications in the InvoiceB2B application using the `vue-sonner` library.

## Changes Made

1. **Added CSS Import**: Added the vue-sonner CSS file to the Nuxt configuration to ensure proper styling of toast notifications.
   ```javascript
   // nuxt.config.ts
   export default defineNuxtConfig({
     css: ['~/assets/css/tailwind.css', 'vue-sonner/lib/index.css'],
     // ...
   });
   ```

2. **Enhanced Toaster Configuration**: Updated the Toaster component in app.vue with additional props to improve functionality and appearance.
   ```vue
   <!-- app.vue -->
   <Toaster 
     richColors 
     position="top-right" 
     :expand="false"
     closeButton
     :duration="4000"
     theme="light"
   />
   ```

## How to Use Toast Notifications

### Basic Usage

To display a toast notification, import the `toast` function from `vue-sonner` and call it with your message:

```javascript
import { toast } from 'vue-sonner';

// Display a default toast
toast('This is a toast message');
```

### Toast Types

Vue Sonner provides different types of toast notifications:

```javascript
// Success toast (green)
toast.success('Operation completed successfully!');

// Error toast (red)
toast.error('An error occurred!');

// Info toast (blue)
toast.info('Here is some information.');

// Warning toast (yellow)
toast.warning('Be careful!');
```

### Toast Options

You can customize individual toast notifications with options:

```javascript
toast.success('Custom toast', {
  description: 'This is a more detailed description',
  duration: 5000, // 5 seconds
  icon: '🚀',
  action: {
    label: 'Undo',
    onClick: () => console.log('Undo clicked')
  }
});
```

### Promise Toasts

You can show toast notifications for async operations:

```javascript
const promise = fetch('/api/data');

toast.promise(promise, {
  loading: 'Loading data...',
  success: 'Data loaded successfully!',
  error: 'Failed to load data'
});
```

## Toaster Component Props

The `<Toaster>` component accepts the following props:

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `position` | String | `'bottom-right'` | Position of the toast notifications |
| `expand` | Boolean | `false` | Whether toasts expand to fill the width |
| `richColors` | Boolean | `false` | Use rich colors for toast types |
| `closeButton` | Boolean | `false` | Show close button on toasts |
| `duration` | Number | `4000` | Default duration in milliseconds |
| `theme` | String | `'light'` | Theme ('light' or 'dark') |

## Best Practices

1. **Keep Messages Concise**: Toast messages should be short and to the point.
2. **Use Appropriate Types**: Use the correct toast type (success, error, etc.) to convey the right meaning.
3. **Set Reasonable Durations**: Important messages should stay longer, while less important ones can be shorter.
4. **Add Actions When Helpful**: For toasts that might require user action, add action buttons.
5. **Don't Overuse**: Toast notifications should be used sparingly to avoid overwhelming the user.

## Troubleshooting

If toast notifications are not appearing or not styled correctly:

1. Ensure the CSS file is properly imported in nuxt.config.ts
2. Check that the Toaster component is included in app.vue
3. Verify that toast functions are being called correctly in components
4. Check for any z-index conflicts that might be hiding the toast notifications