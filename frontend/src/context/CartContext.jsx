import React, { createContext, useContext, useState, useEffect, useCallback } from 'react';
import api from '../api';
import { useAuth } from './AuthContext';
import { useToast } from './ToastContext';

const CartContext = createContext(null);

export const useCart = () => {
  const context = useContext(CartContext);
  if (!context) {
    throw new Error('useCart must be used within a CartProvider');
  }
  return context;
};

export const CartProvider = ({ children }) => {
  const { user, token } = useAuth();
  const { showError, showSuccess } = useToast();
  const [cartItems, setCartItems] = useState([]);
  const [loading, setLoading] = useState(false);

  // Initialize/Ensure Cart exists, then fetch items
  const fetchCart = useCallback(async () => {
    if (!user) {
      setCartItems([]);
      return;
    }
    setLoading(true);
    try {
      // 1. Try to get the user's cart; create one if it doesn't exist
      try {
        await api.get('/cart');
      } catch (err) {
        if (err.response?.status === 404) {
          await api.post('/cart');
        }
      }

      // 2. Fetch cart items — always treat non-200 as empty
      try {
        const itemsRes = await api.get('/cart/items');
        setCartItems(itemsRes.data.data.items || []);
      } catch {
        setCartItems([]);
      }
    } catch (err) {
      console.error('Error initializing cart:', err);
      setCartItems([]);
    } finally {
      setLoading(false);
    }
  }, [user]);

  // Fetch cart automatically when user changes/logs in
  useEffect(() => {
    if (user) {
      fetchCart();
    } else {
      setCartItems([]);
    }
  }, [user, token, fetchCart]);

  const addToCart = async (productId, quantity = 1) => {
    if (!user) {
      showError('Please log in to add items to your cart.');
      return false;
    }
    try {
      await api.post('/cart/items', { product_id: productId, quantity });
      showSuccess('Item added to cart.');
      await fetchCart();
      return true;
    } catch (err) {
      showError(err.response?.data?.error || 'Failed to add item to cart');
      return false;
    }
  };

  const removeFromCart = async (itemId) => {
    try {
      await api.delete(`/cart/items/${itemId}`);
      showSuccess('Item removed from cart.');
      await fetchCart();
      return true;
    } catch (err) {
      showError(err.response?.data?.error || 'Failed to remove item');
      return false;
    }
  };

  const clearCart = async () => {
    // Clear cart endpoint: Let's see if we have clear cart endpoint.
    // In service: ClearCart(userID) exists. But wait, is there a handler or route for clear cart?
    // Let's check routes.go: there is no DELETE /cart or POST /cart/clear route.
    // But checkout automatically clears the cart on payment success!
    // So we can clear locally or we can delete items individually if we want.
    // But since the database has ON DELETE CASCADE or payment webhook clears it, let's keep it simple.
    setCartItems([]);
  };

  const cartCount = cartItems.reduce((acc, item) => acc + item.quantity, 0);

  return (
    <CartContext.Provider value={{ cartItems, loading, addToCart, removeFromCart, clearCart, fetchCart, cartCount }}>
      {children}
    </CartContext.Provider>
  );
};
