import React, { createContext, useContext, useState, useEffect } from 'react';
import api from '../api';
import { useToast } from './ToastContext';

const AuthContext = createContext(null);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

// Helper function to decode JWT claims
const decodeToken = (token) => {
  try {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split('')
        .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    );
    return JSON.parse(jsonPayload);
  } catch (error) {
    return null;
  }
};

export const AuthProvider = ({ children }) => {
  const [token, setToken] = useState(localStorage.getItem('token'));
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const { showError, showSuccess } = useToast();

  useEffect(() => {
    if (token) {
      const claims = decodeToken(token);
      if (claims && claims.exp * 1000 > Date.now()) {
        setUser({
          id: claims.user_id,
          role: claims.role,
        });
      } else {
        // Expired
        logout();
      }
    } else {
      setUser(null);
    }
    setLoading(false);
  }, [token]);

  const login = async (email, password) => {
    try {
      const response = await api.post('/auth/login', { email, password });
      const jwtToken = response.data.data.token; // check structure: { data: { token: '...' } }
      localStorage.setItem('token', jwtToken);
      setToken(jwtToken);
      showSuccess('Successfully logged in!');
      return true;
    } catch (err) {
      const errorMsg = err.response?.data?.error || 'Failed to login';
      showError(errorMsg);
      return false;
    }
  };

  const register = async (fullName, email, password) => {
    try {
      await api.post('/auth/register', { full_name: fullName, email, password });
      showSuccess('Account created successfully! Please log in.');
      return true;
    } catch (err) {
      const errorMsg = err.response?.data?.error || 'Failed to register';
      showError(errorMsg);
      return false;
    }
  };

  const logout = () => {
    localStorage.removeItem('token');
    setToken(null);
    setUser(null);
    showSuccess('Successfully logged out.');
  };

  const isAdmin = user?.role === 'admin';

  return (
    <AuthContext.Provider value={{ user, token, loading, login, register, logout, isAdmin }}>
      {children}
    </AuthContext.Provider>
  );
};
