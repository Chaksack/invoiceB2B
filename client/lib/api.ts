import axios from 'axios';

// Get the API base URL from environment variables or use a default
const API_BASE_URL = process.env.NUXT_PUBLIC_API_BASE_URL || 'http://localhost:8080/api/v1';

/**
 * Create an axios instance with default configuration
 * @param token - Optional authentication token
 * @returns Configured axios instance
 */
export const createApiClient = (token?: string) => {
  return axios.create({
    baseURL: API_BASE_URL,
    headers: {
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      'Content-Type': 'application/json',
    },
  });
};

/**
 * Financial Institutions API
 */
export const financialInstitutionsApi = {
  /**
   * Get a paginated list of financial institutions with optional filters
   */
  getAll: (apiClient, params = {}) => {
    return apiClient.get('/admin/financial-institutions', { params });
  },

  /**
   * Get a specific financial institution by ID
   */
  getById: (apiClient, id) => {
    return apiClient.get(`/admin/financial-institutions/${id}`);
  },

  /**
   * Create a new financial institution
   */
  create: (apiClient, data) => {
    return apiClient.post('/admin/financial-institutions', data);
  },

  /**
   * Update an existing financial institution
   */
  update: (apiClient, id, data) => {
    return apiClient.put(`/admin/financial-institutions/${id}`, data);
  },

  /**
   * Delete a financial institution
   */
  delete: (apiClient, id) => {
    return apiClient.delete(`/admin/financial-institutions/${id}`);
  },

  /**
   * Get products for a financial institution
   */
  getProducts: (apiClient, fiId) => {
    return apiClient.get(`/admin/financial-institutions/${fiId}/products`);
  },

  /**
   * Get terms for a financial institution
   */
  getTerms: (apiClient, fiId) => {
    return apiClient.get(`/admin/financial-institutions/${fiId}/terms`);
  },
};

/**
 * User API
 */
export const userApi = {
  /**
   * Get the current user's profile
   */
  getProfile: (apiClient) => {
    return apiClient.get('/user/profile');
  },

  /**
   * Update the current user's profile
   */
  updateProfile: (apiClient, data) => {
    return apiClient.put('/user/profile', data);
  },
};

/**
 * Dashboard API
 */
export const dashboardApi = {
  /**
   * Get dashboard analytics
   */
  getAnalytics: (apiClient) => {
    return apiClient.get('/admin/dashboard/analytics');
  },
};