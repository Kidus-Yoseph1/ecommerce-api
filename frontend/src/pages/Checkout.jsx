import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { loadStripe } from '@stripe/stripe-js';
import { Elements, CardElement, useStripe, useElements } from '@stripe/react-stripe-js';
import { CreditCard, ShoppingBag, ShieldCheck, ArrowLeft } from 'lucide-react';
import api from '../api';
import { useCart } from '../context/CartContext';
import { useToast } from '../context/ToastContext';

// Use env key or standard test key for placeholder
const stripePromise = loadStripe(import.meta.env.VITE_STRIPE_PUBLIC_KEY || 'pk_test_51I1zV7H8C7pXyEuxW12H4w5J7K7lR8m9P0q1s2t3u4v5w6x7y8z');

const CARD_ELEMENT_OPTIONS = {
  style: {
    base: {
      color: '#f1f5f9',
      fontFamily: '"Inter", sans-serif',
      fontSmoothing: 'antialiased',
      fontSize: '16px',
      '::placeholder': {
        color: '#475569',
      },
    },
    invalid: {
      color: '#ef4444',
      iconColor: '#ef4444',
    },
  },
};

function CheckoutForm({ clientSecret, orderId }) {
  const stripe = useStripe();
  const elements = useElements();
  const { showError, showSuccess } = useToast();
  const { clearCart } = useCart();
  const navigate = useNavigate();
  const [processing, setProcessing] = useState(false);

  const handleSubmit = async (event) => {
    event.preventDefault();

    if (!stripe || !elements) return;

    setProcessing(true);

    try {
      const cardElement = elements.getElement(CardElement);
      const result = await stripe.confirmCardPayment(clientSecret, {
        payment_method: {
          card: cardElement,
        },
      });

      if (result.error) {
        showError(result.error.message);
        setProcessing(false);
      } else {
        if (result.paymentIntent.status === 'succeeded') {
          showSuccess('Payment successful! Your order has been placed.');
          clearCart();
          navigate(`/order-confirmation/${orderId}`);
        }
      }
    } catch (err) {
      showError('Payment processing failed.');
      setProcessing(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
      <div className="input-group">
        <label>Card Information</label>
        <div
          style={{
            background: 'var(--bg-secondary)',
            border: '1px solid var(--border)',
            borderRadius: 'var(--radius-md)',
            padding: '16px',
          }}
        >
          <CardElement options={CARD_ELEMENT_OPTIONS} />
        </div>
      </div>

      <button
        type="submit"
        disabled={!stripe || processing}
        className="btn btn-primary btn-full"
        style={{ height: '52px' }}
      >
        {processing ? (
          <span className="spinner" style={{ width: '20px', height: '20px', borderWidth: '2px' }}></span>
        ) : (
          <>
            <CreditCard size={18} />
            Pay Now
          </>
        )}
      </button>
    </form>
  );
}

export default function Checkout() {
  const { cartItems } = useCart();
  const { showError } = useToast();
  const navigate = useNavigate();
  const [clientSecret, setClientSecret] = useState(null);
  const [orderId, setOrderId] = useState(null);
  const [loading, setLoading] = useState(true);
  const [productMap, setProductMap] = useState({});

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
        console.error(err);
      }
    };
    fetchProducts();
  }, []);

  useEffect(() => {
    const initiateCheckout = async () => {
      if (cartItems.length === 0) {
        navigate('/cart');
        return;
      }
      try {
        const res = await api.post('/checkout');
        setClientSecret(res.data.data.client_secret);
        setOrderId(res.data.data.order_id);
      } catch (err) {
        showError(err.response?.data?.error || 'Failed to initialize checkout');
        navigate('/cart');
      } finally {
        setLoading(false);
      }
    };
    initiateCheckout();
  }, [cartItems, navigate, showError]);

  const subtotal = cartItems.reduce((acc, item) => {
    const price = productMap[item.product_id]?.price || item.price || 0;
    return acc + price * item.quantity;
  }, 0);

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <div className="spinner"></div>
      </div>
    );
  }

  return (
    <div style={{ paddingTop: '120px', paddingBottom: '100px' }} className="animate-fade">
      <div className="container">
        <button
          onClick={() => navigate('/cart')}
          className="btn btn-ghost btn-sm"
          style={{ marginBottom: '40px' }}
        >
          <ArrowLeft size={16} />
          Return to Cart
        </button>

        <h1 style={{ fontSize: '32px', fontWeight: 800, marginBottom: '8px' }}>Stripe Secure Checkout</h1>
        <p style={{ color: 'var(--text-secondary)', marginBottom: '40px' }}>
          Complete payment using our secure interface.
        </p>

        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(320px, 1fr))',
            gap: '40px',
            alignItems: 'start',
          }}
        >
          {/* Form */}
          <div className="card-glass" style={{ padding: '36px' }}>
            <h3 style={{ fontSize: '20px', fontWeight: 700, marginBottom: '24px' }}>Payment details</h3>
            {clientSecret && (
              <Elements stripe={stripePromise} options={{ clientSecret }}>
                <CheckoutForm clientSecret={clientSecret} orderId={orderId} />
              </Elements>
            )}
          </div>

          {/* Cart review */}
          <div className="card" style={{ padding: '32px' }}>
            <h3 style={{ fontSize: '18px', fontWeight: 700, marginBottom: '20px' }}>Review Order</h3>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
              {cartItems.map((item) => {
                const product = productMap[item.product_id] || { name: 'Item', price: item.price };
                return (
                  <div key={item.id} className="flex justify-between items-center" style={{ fontSize: '14px' }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                      <div
                        style={{
                          width: '40px',
                          height: '40px',
                          borderRadius: 'var(--radius-sm)',
                          background: 'var(--bg-secondary)',
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                          border: '1px solid var(--border)',
                        }}
                      >
                        <ShoppingBag size={18} color="var(--accent-from)" />
                      </div>
                      <div>
                        <div style={{ fontWeight: 600, color: 'var(--text-primary)' }}>{product.name}</div>
                        <div style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>Qty: {item.quantity}</div>
                      </div>
                    </div>
                    <span style={{ fontWeight: 700 }}>
                      ${((product.price || 0) * item.quantity).toFixed(2)}
                    </span>
                  </div>
                );
              })}

              <div className="divider"></div>

              <div className="flex justify-between" style={{ fontSize: '16px', fontWeight: 700 }}>
                <span>Subtotal</span>
                <span>${subtotal.toFixed(2)}</span>
              </div>
              <div className="flex justify-between" style={{ fontSize: '16px', fontWeight: 700 }}>
                <span>Shipping</span>
                <span className="badge badge-green" style={{ fontSize: '9px' }}>Free</span>
              </div>

              <div className="divider" style={{ margin: '8px 0' }}></div>

              <div className="flex justify-between" style={{ fontSize: '18px', fontWeight: 800 }}>
                <span>Total Due</span>
                <span className="gradient-text">${subtotal.toFixed(2)}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
