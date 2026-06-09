import React, { useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Trash2, ShoppingBag, ArrowRight, ShieldCheck } from 'lucide-react';
import { useCart } from '../context/CartContext';
import api from '../api';

export default function Cart() {
  const { cartItems, removeFromCart, loading } = useCart();
  const [productMap, setProductMap] = useState({});
  const [fetchingProducts, setFetchingProducts] = useState(true);
  const navigate = useNavigate();

  // Fetch products to map product ID -> Name, Category, etc.
  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const res = await api.get('/products');
        const list = res.data.data.products || [];
        const map = {};
        list.forEach((p) => {
          map[p.id] = p;
        });
        setProductMap(map);
      } catch (err) {
        console.error('Failed to map products:', err);
      } finally {
        setFetchingProducts(false);
      }
    };
    fetchProducts();
  }, [cartItems]);

  const subtotal = cartItems.reduce((acc, item) => {
    const price = productMap[item.product_id]?.price || item.price || 0;
    return acc + price * item.quantity;
  }, 0);

  if (loading || fetchingProducts) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <div className="spinner"></div>
      </div>
    );
  }

  if (cartItems.length === 0) {
    return (
      <div
        style={{
          minHeight: '100vh',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          padding: '120px 24px 80px',
        }}
        className="animate-fade"
      >
        <div
          className="card-glass text-center"
          style={{
            padding: '60px 40px',
            maxWidth: '480px',
            width: '100%',
          }}
        >
          <div
            style={{
              width: '80px',
              height: '80px',
              borderRadius: '50%',
              background: 'var(--bg-secondary)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              margin: '0 auto 24px',
            }}
          >
            <ShoppingBag size={36} color="var(--text-muted)" />
          </div>
          <h2 style={{ fontSize: '24px', fontWeight: 800, marginBottom: '12px' }}>Your Cart is Empty</h2>
          <p style={{ color: 'var(--text-secondary)', fontSize: '15px', marginBottom: '32px' }}>
            Looks like you haven't added any products to your cart yet.
          </p>
          <Link to="/" className="btn btn-primary btn-full">
            Explore Collection
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div style={{ paddingTop: '120px', paddingBottom: '100px' }} className="animate-fade">
      <div className="container">
        <h1 style={{ fontSize: '32px', fontWeight: 800, marginBottom: '8px' }}>Your Shopping Cart</h1>
        <p style={{ color: 'var(--text-secondary)', marginBottom: '40px' }}>
          Manage items and proceed to Stripe checkout.
        </p>

        {/* Layout Grid */}
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(320px, 1fr))',
            gap: '40px',
            alignItems: 'start',
          }}
        >
          {/* Items list */}
          <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
            {cartItems.map((item) => {
              const product = productMap[item.product_id] || {
                name: 'Loading Item...',
                category: 'General',
                price: item.price,
              };

              return (
                <div
                  key={item.id}
                  className="card flex"
                  style={{
                    padding: '20px',
                    gap: '20px',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                  }}
                >
                  {/* Thumbnail */}
                  <div
                    style={{
                      width: '64px',
                      height: '64px',
                      borderRadius: 'var(--radius-md)',
                      background: 'var(--bg-secondary)',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      border: '1px solid var(--border)',
                    }}
                  >
                    <ShoppingBag size={24} color="var(--accent-from)" />
                  </div>

                  {/* Info */}
                  <div style={{ flex: 1 }}>
                    <span className="badge badge-purple" style={{ fontSize: '9px', padding: '2px 6px', marginBottom: '4px' }}>
                      {product.category}
                    </span>
                    <h4 style={{ fontSize: '16px', fontWeight: 700, color: 'var(--text-primary)' }}>
                      {product.name}
                    </h4>
                    <span style={{ fontSize: '13px', color: 'var(--text-secondary)' }}>
                      Qty: {item.quantity} &times; ${product.price?.toFixed(2) || '0.00'}
                    </span>
                  </div>

                  {/* Price & Delete */}
                  <div style={{ display: 'flex', alignItems: 'center', gap: '20px' }}>
                    <span style={{ fontSize: '16px', fontWeight: 800, color: 'var(--text-primary)' }}>
                      ${((product.price || 0) * item.quantity).toFixed(2)}
                    </span>
                    <button
                      onClick={() => removeFromCart(item.id)}
                      style={{
                        background: 'none',
                        border: 'none',
                        color: 'var(--red)',
                        cursor: 'pointer',
                        display: 'flex',
                        alignItems: 'center',
                        padding: '8px',
                        borderRadius: '50%',
                        transition: 'background-color 0.2s',
                      }}
                      onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = 'rgba(239, 68, 68, 0.1)')}
                      onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = 'transparent')}
                    >
                      <Trash2 size={16} />
                    </button>
                  </div>
                </div>
              );
            })}
          </div>

          {/* Checkout Summary Card */}
          <div className="card-glass" style={{ padding: '32px' }}>
            <h3 style={{ fontSize: '20px', fontWeight: 700, marginBottom: '20px' }}>Order Summary</h3>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
              <div className="flex justify-between" style={{ fontSize: '15px', color: 'var(--text-secondary)' }}>
                <span>Subtotal</span>
                <span style={{ color: 'var(--text-primary)', fontWeight: 500 }}>
                  ${subtotal.toFixed(2)}
                </span>
              </div>
              <div className="flex justify-between" style={{ fontSize: '15px', color: 'var(--text-secondary)' }}>
                <span>Shipping</span>
                <span className="badge badge-green" style={{ fontSize: '10px' }}>
                  Free
                </span>
              </div>
              <div className="flex justify-between" style={{ fontSize: '15px', color: 'var(--text-secondary)' }}>
                <span>Tax</span>
                <span style={{ color: 'var(--text-primary)', fontWeight: 500 }}>$0.00</span>
              </div>

              <div className="divider" style={{ margin: '8px 0' }}></div>

              <div className="flex justify-between" style={{ fontSize: '18px', fontWeight: 800 }}>
                <span>Total</span>
                <span className="gradient-text">${subtotal.toFixed(2)}</span>
              </div>

              <button
                onClick={() => navigate('/checkout')}
                className="btn btn-primary btn-full mt-4"
                style={{ height: '48px' }}
              >
                Proceed to Checkout
                <ArrowRight size={16} />
              </button>

              <div
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: '8px',
                  justifyContent: 'center',
                  marginTop: '16px',
                  fontSize: '12px',
                  color: 'var(--text-secondary)',
                }}
              >
                <ShieldCheck size={16} color="var(--green)" />
                <span>Encrypted checkout using Stripe.</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
