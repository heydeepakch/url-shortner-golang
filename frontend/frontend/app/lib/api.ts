import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/api', // Your Go backend URL
});

// Add a request interceptor to attach the token if it exists
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token'); // Or use cookies
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default api;