# Sonner Toast Component

This is a wrapper around the [vue-sonner](https://github.com/wobsoriano/vue-sonner) toast notification library, styled to match our application's design system.

## Usage

### 1. Import the Toaster Component

The Toaster component is already included in the application layout (app.vue), so you don't need to add it to your components.

### 2. Import the toast function

```javascript
import { toast } from 'vue-sonner';
```

### 3. Use the toast function to display notifications

```javascript
// Success notification
toast.success('Operation completed successfully!');

// Error notification
toast.error('An error occurred. Please try again.');

// Info notification
toast.info('Your session will expire in 5 minutes.');

// Warning notification
toast.warning('This action cannot be undone.');

// Custom notification
toast('Custom notification', {
  description: 'This is a custom notification with a description',
  action: {
    label: 'Undo',
    onClick: () => console.log('Undo clicked')
  }
});
```

### 4. Customize toast appearance

You can customize the appearance of toasts by passing options to the toast function:

```javascript
toast.success('Success message', {
  duration: 5000, // Duration in milliseconds
  position: 'bottom-center', // Position of the toast
  icon: '👍', // Custom icon
  className: 'my-custom-class', // Custom CSS class
  description: 'This is a more detailed description of the success message'
});
```

## Available Options

For a complete list of available options, refer to the [vue-sonner documentation](https://github.com/wobsoriano/vue-sonner).

## Examples

### Form Validation Error

```javascript
const validateForm = () => {
  if (!formData.email) {
    toast.error('Please enter your email address');
    return false;
  }
  return true;
};
```

### Successful Form Submission

```javascript
const submitForm = () => {
  // Submit form data to API
  api.submitForm(formData)
    .then(() => {
      toast.success('Form submitted successfully!');
      router.push('/dashboard');
    })
    .catch(error => {
      toast.error(`Error: ${error.message}`);
    });
};
```

### Loading State

```javascript
const fetchData = async () => {
  const toastId = toast.loading('Fetching data...');
  
  try {
    const data = await api.fetchData();
    toast.success('Data loaded successfully!', { id: toastId });
    return data;
  } catch (error) {
    toast.error(`Error: ${error.message}`, { id: toastId });
    return null;
  }
};
```

## Styling

The toast component uses the application's design tokens for consistent styling. The default styling can be found in the Sonner.vue component.