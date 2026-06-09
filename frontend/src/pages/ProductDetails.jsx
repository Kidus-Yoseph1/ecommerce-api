import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { ArrowLeft, ShoppingCart, Plus, Minus, ShieldCheck, Truck, RefreshCw } from 'lucide-react';
import api from '../api';
import { useCart } from '../context/CartContext';

export default function ProductDetails() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { addToCart } = useCart();
  const [product, setProduct] = useState(null);
  const [loading, setLoading] = useState(true);
  const [quantity, setQuantity] = useState(1);

  useEffect(() => {
    const fetchProduct = async () => {
      try {
        const res = await api.get(`/products/${id}`);
        setProduct(res.data.data.product);
      } catch (err) {
        console.error(err);
      } finally {
        setLoading(false);
      }
    };
    fetchProduct();
  }, [id]);

  const handleIncrement = () => {
    if (quantity < (product?.stock_quantity || 1)) {
      setQuantity((prev) => prev + 1);
    }
  };

  const handleDecrement = () => {
    if (quantity > 1) {
      setQuantity((prev) => prev - 1);
    }
  };

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <div className="spinner"></div>
      </div>
    );
  }

  if (!product) {
    return (
      <div style={{ textAlign: 'center', padding: '120px 24px' }}>
        <h2>Product not found.</h2>
        <button onClick={() => navigate('/')} className="btn btn-primary mt-4">
          Go Back Home
        </button>
      </div>
    );
  }

  return (
    <div style={{ paddingTop: '120px', paddingBottom: '80px' }} className="animate-fade">
      <div className="container">
        {/* Back Button */}
        <button
          onClick={() => navigate(-1)}
          className="btn btn-ghost btn-sm"
          style={{ marginBottom: '40px' }}
        >
          <ArrowLeft size={16} />
          Back to Listings
        </button>

        {/* Layout */}
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(320px, 1fr))',
            gap: '60px',
            alignItems: 'start',
          }}
        >
          {/* Visual Container */}
          <div
            className="card-glass"
            style={{
              height: '440px',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              background: 'var(--bg-secondary)',
              border: '1px solid var(--border)',
              borderRadius: 'var(--radius-lg)',
              position: 'relative',
            }}
          >
            <span
              className="badge badge-purple"
              style={{ position: 'absolute', top: '24px', left: '24px' }}
            >
              {product.category}
            </span>

            {/* Simulated item icon render */}
            <div
              style={{
                width: '120px',
                height: '120px',
                borderRadius: '50%',
                background: 'linear-gradient(135deg, var(--accent-glow), transparent)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                boxShadow: 'var(--shadow-glow)',
              }}
            >
              <ShoppingCart size={54} color="var(--accent-from)" />
            </div>
          </div>

          {/* Details Container */}
          <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
            <div>
              <h1 style={{ fontSize: '38px', fontWeight: 800, marginBottom: '12px' }}>{product.name}</h1>
              <span style={{ fontSize: '28px', fontWeight: 800, color: 'var(--text-primary)' }}>
                ${product.price.toFixed(2)}
              </span>
            </div>

            <div className="divider" style={{ margin: 0 }}></div>

            <div>
              <h3 style={{ fontSize: '16px', fontWeight: 600, marginBottom: '10px', color: 'var(--text-secondary)' }}>
                Description
              </h3>
              <p style={{ color: 'var(--text-secondary)', fontSize: '15px', lineHeight: 1.7 }}>
                {product.description || 'No description provided for this premium item.'}
              </p>
            </div>

            {/* Stock status indicator */}
            <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
              <span style={{ fontSize: '14px', fontWeight: 500, color: 'var(--text-secondary)' }}>Status:</span>
              {product.stock_quantity === 0 ? (
                <span className="badge badge-red">Out of Stock</span>
              ) : product.stock_quantity < 5 ? (
                <span className="badge badge-yellow">Low Stock — Only {product.stock_quantity} remaining</span>
              ) : (
                <span className="badge badge-green">In Stock — {product.stock_quantity} available</span>
              )}
            </div>

            {product.stock_quantity > 0 && (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                <label style={{ fontSize: '13px', fontWeight: 600, color: 'var(--text-secondary)', textTransform: 'uppercase' }}>
                  Quantity
                </label>
                <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
                  <div
                    style={{
                      display: 'flex',
                      alignItems: 'center',
                      background: 'var(--bg-secondary)',
                      border: '1px solid var(--border)',
                      borderRadius: 'var(--radius-md)',
                      height: '46px',
                    }}
                  >
                    <button
                      onClick={handleDecrement}
                      style={{
                        background: 'none',
                        border: 'none',
                        color: 'var(--text-primary)',
                        padding: '0 16px',
                        cursor: 'pointer',
                      }}
                    >
                      <Minus size={16} />
                    </button>
                    <span style={{ fontSize: '16px', fontWeight: 600, width: '20px', textAlign: 'center' }}>
                      {quantity}
                    </span>
                    <button
                      onClick={handleIncrement}
                      style={{
                        background: 'none',
                        border: 'none',
                        color: 'var(--text-primary)',
                        padding: '0 16px',
                        cursor: 'pointer',
                      }}
                    >
                      <Plus size={16} />
                    </button>
                  </div>

                  <button
                    onClick={() => addToCart(product.id, quantity)}
                    className="btn btn-primary"
                    style={{ height: '46px', flex: 1 }}
                  >
                    <ShoppingCart size={18} />
                    Add to Cart
                  </button>
                </div>
              </div>
            )}

            <div className="divider" style={{ margin: 0 }}></div>

            {/* Dynamic trust signals */}
            <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '10px', fontSize: '13px', color: 'var(--text-secondary)' }}>
                <ShieldCheck size={16} color="var(--accent-from)" />
                <span>Secure payment with Stripe checkout encryption.</span>
              </div>
              <div style={{ display: 'flex', alignItems: 'center', gap: '10px', fontSize: '13px', color: 'var(--text-secondary)' }}>
                <Truck size={16} color="var(--accent-from)" />
                <span>Free tracked delivery for premium items.</span>
              </div>
              <div style={{ display: 'flex', alignItems: 'center', gap: '10px', fontSize: '13px', color: 'var(--text-secondary)' }}>
                <RefreshCw size={16} color="var(--accent-from)" />
                <span>30-day return policy for peace of mind.</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
