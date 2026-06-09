import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { CheckCircle, ShoppingBag, ArrowRight, ShieldCheck, Mail } from 'lucide-react';
import api from '../api';

export default function OrderConfirmation() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [order, setOrder] = useState(null);
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [productMap, setProductMap] = useState({});

  useEffect(() => {
    const fetchOrderDetails = async () => {
      try {
        // Fetch order basic details
        const orderRes = await api.get(`/orders/${id}`);
        setOrder(orderRes.data.data.order);

        // Fetch order items
        const itemsRes = await api.get(`/orders/${id}/items`);
        setItems(itemsRes.data.data.items || []);

        // Fetch products mapping
        const productsRes = await api.get('/products');
        const list = productsRes.data.data.products || [];
        const map = {};
        list.forEach((p) => {
          map[p.id] = p;
        });
        setProductMap(map);
      } catch (err) {
        console.error('Failed to load order confirmation:', err);
      } finally {
        setLoading(false);
      }
    };
    fetchOrderDetails();
  }, [id]);

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <div className="spinner"></div>
      </div>
    );
  }

  return (
    <div style={{ paddingTop: '120px', paddingBottom: '100px' }} className="animate-fade">
      <div className="container" style={{ maxWidth: '640px' }}>
        <div className="card-glass text-center" style={{ padding: '40px', marginBottom: '32px' }}>
          <div
            style={{
              width: '72px',
              height: '72px',
              borderRadius: '50%',
              background: 'rgba(16, 185, 129, 0.15)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              margin: '0 auto 24px',
              border: '1px solid rgba(16, 185, 129, 0.3)',
            }}
          >
            <CheckCircle size={36} color="var(--green)" />
          </div>
          <h1 style={{ fontSize: '32px', fontWeight: 800, marginBottom: '12px' }}>Order Completed</h1>
          <p style={{ color: 'var(--text-secondary)', fontSize: '15px' }}>
            Thank you for your purchase! Your payment has been verified.
          </p>
          <div
            className="badge badge-purple"
            style={{ marginTop: '16px', fontSize: '11px', padding: '6px 12px' }}
          >
            ORDER ID: {id}
          </div>
        </div>

        <div className="card" style={{ padding: '32px', display: 'flex', flexDirection: 'column', gap: '20px' }}>
          <h3 style={{ fontSize: '18px', fontWeight: 700 }}>Summary Receipt</h3>

          <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
            {items.map((item) => {
              const product = productMap[item.product_id] || { name: 'Item Snapshot', price: item.unit_price };
              return (
                <div key={item.id} className="flex justify-between items-center" style={{ fontSize: '14px' }}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
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
                      <div style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>
                        Qty: {item.quantity} &times; ${item.unit_price.toFixed(2)}
                      </div>
                    </div>
                  </div>
                  <span style={{ fontWeight: 700, color: 'var(--text-primary)' }}>
                    ${(item.unit_price * item.quantity).toFixed(2)}
                  </span>
                </div>
              );
            })}

            <div className="divider" style={{ margin: '8px 0' }}></div>

            <div className="flex justify-between" style={{ fontSize: '15px', color: 'var(--text-secondary)' }}>
              <span>Total Paid</span>
              <span className="gradient-text" style={{ fontSize: '18px', fontWeight: 800 }}>
                ${order?.total_amount?.toFixed(2) || '0.00'}
              </span>
            </div>
            <div className="flex justify-between" style={{ fontSize: '15px', color: 'var(--text-secondary)' }}>
              <span>Status</span>
              <span className={`badge badge-${order?.status === 'paid' ? 'green' : 'yellow'}`}>
                {order?.status}
              </span>
            </div>
          </div>
        </div>

        <div style={{ display: 'flex', gap: '16px', marginTop: '32px' }}>
          <Link to="/" className="btn btn-primary btn-full" style={{ height: '48px' }}>
            Continue Shopping
            <ArrowRight size={16} />
          </Link>
        </div>
      </div>
    </div>
  );
}
