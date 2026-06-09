import React, { useState, useEffect } from 'react';
import { Plus, Edit2, Trash2, X, AlertCircle } from 'lucide-react';
import api from '../api';
import { useToast } from '../context/ToastContext';

export default function AdminDashboard() {
  const { showSuccess, showError } = useToast();
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingProduct, setEditingProduct] = useState(null);

  // Form states
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [category, setCategory] = useState('');
  const [price, setPrice] = useState('');
  const [stockQuantity, setStockQuantity] = useState('');

  const fetchProducts = async () => {
    setLoading(true);
    try {
      const res = await api.get('/products');
      setProducts(res.data.data.products || []);
    } catch (err) {
      console.error(err);
      showError('Failed to fetch product catalog.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProducts();
  }, []);

  const openAddModal = () => {
    setEditingProduct(null);
    setName('');
    setDescription('');
    setCategory('Electronics');
    setPrice('');
    setStockQuantity('');
    setIsModalOpen(true);
  };

  const openEditModal = (product) => {
    setEditingProduct(product);
    setName(product.name);
    setDescription(product.description || '');
    setCategory(product.category);
    setPrice(product.price.toString());
    setStockQuantity(product.stock_quantity.toString());
    setIsModalOpen(true);
  };

  const closeModal = () => {
    setIsModalOpen(false);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!name || !category || !price || !stockQuantity) {
      showError('All fields are required.');
      return;
    }

    const payload = {
      name,
      description,
      category,
      price: parseFloat(price),
      stock_quantity: parseInt(stockQuantity, 10),
    };

    try {
      if (editingProduct) {
        // Update product
        await api.put(`/products/${editingProduct.id}`, payload);
        showSuccess('Product updated successfully.');
      } else {
        // Create product
        await api.post('/products', payload);
        showSuccess('Product created successfully.');
      }
      closeModal();
      await fetchProducts();
    } catch (err) {
      showError(err.response?.data?.error || 'Operation failed.');
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this product? This is a soft delete action.')) {
      return;
    }
    try {
      await api.delete(`/products/${id}`);
      showSuccess('Product deleted successfully.');
      await fetchProducts();
    } catch (err) {
      showError(err.response?.data?.error || 'Delete failed.');
    }
  };

  return (
    <div style={{ paddingTop: '120px', paddingBottom: '100px' }} className="animate-fade">
      <div className="container">
        {/* Header toolbar */}
        <div
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            marginBottom: '40px',
          }}
        >
          <div>
            <h1 style={{ fontSize: '32px', fontWeight: 800, marginBottom: '8px' }}>Store Admin Portal</h1>
            <p style={{ color: 'var(--text-secondary)' }}>Manage store product inventory, pricing, and stock.</p>
          </div>
          <button onClick={openAddModal} className="btn btn-primary">
            <Plus size={18} />
            Create Product
          </button>
        </div>

        {/* Dashboard table */}
        {loading ? (
          <div style={{ display: 'flex', justifyContent: 'center', padding: '60px' }}>
            <div className="spinner"></div>
          </div>
        ) : products.length === 0 ? (
          <div style={{ textAlign: 'center', padding: '60px 24px', background: 'var(--bg-secondary)', borderRadius: 'var(--radius-lg)' }}>
            <p style={{ color: 'var(--text-secondary)' }}>Inventory is empty. Get started by adding a product.</p>
          </div>
        ) : (
          <div className="card" style={{ overflowX: 'auto', border: '1px solid var(--border)' }}>
            <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left', fontSize: '14px' }}>
              <thead>
                <tr style={{ background: 'var(--bg-secondary)', borderBottom: '1px solid var(--border)', color: 'var(--text-secondary)' }}>
                  <th style={{ padding: '16px 24px', fontWeight: 600 }}>Product Name</th>
                  <th style={{ padding: '16px 24px', fontWeight: 600 }}>Category</th>
                  <th style={{ padding: '16px 24px', fontWeight: 600 }}>Price</th>
                  <th style={{ padding: '16px 24px', fontWeight: 600 }}>Stock Level</th>
                  <th style={{ padding: '16px 24px', fontWeight: 600, textAlign: 'right' }}>Actions</th>
                </tr>
              </thead>
              <tbody>
                {products.map((p) => (
                  <tr key={p.id} style={{ borderBottom: '1px solid var(--border)', transition: 'background 0.2s' }}>
                    <td style={{ padding: '16px 24px', fontWeight: 600, color: 'var(--text-primary)' }}>{p.name}</td>
                    <td style={{ padding: '16px 24px' }}>
                      <span className="badge badge-purple" style={{ fontSize: '10px' }}>{p.category}</span>
                    </td>
                    <td style={{ padding: '16px 24px', fontWeight: 700 }}>${p.price.toFixed(2)}</td>
                    <td style={{ padding: '16px 24px' }}>
                      {p.stock_quantity === 0 ? (
                        <span className="badge badge-red" style={{ display: 'inline-flex', alignItems: 'center', gap: '4px' }}>
                          <AlertCircle size={10} />
                          Out of stock
                        </span>
                      ) : p.stock_quantity < 5 ? (
                        <span className="badge badge-yellow">{p.stock_quantity} left</span>
                      ) : (
                        <span className="badge badge-green">{p.stock_quantity} units</span>
                      )}
                    </td>
                    <td style={{ padding: '16px 24px', textAlign: 'right' }}>
                      <div style={{ display: 'flex', gap: '8px', justifyContent: 'flex-end' }}>
                        <button
                          onClick={() => openEditModal(p)}
                          className="btn btn-secondary btn-sm"
                          style={{ padding: '8px', borderRadius: '50%' }}
                          title="Edit"
                        >
                          <Edit2 size={14} />
                        </button>
                        <button
                          onClick={() => handleDelete(p.id)}
                          className="btn btn-danger btn-sm"
                          style={{ padding: '8px', borderRadius: '50%' }}
                          title="Delete"
                        >
                          <Trash2 size={14} />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {/* Modal Overlay Form */}
        {isModalOpen && (
          <div
            style={{
              position: 'fixed',
              top: 0,
              left: 0,
              right: 0,
              bottom: 0,
              zIndex: 2000,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              background: 'rgba(0, 0, 0, 0.7)',
              backdropFilter: 'blur(8px)',
            }}
          >
            <div
              className="card-glass"
              style={{
                width: '100%',
                maxWidth: '500px',
                padding: '36px',
                animation: 'fadeInUp 0.3s ease',
              }}
            >
              <div className="flex justify-between items-center" style={{ marginBottom: '24px' }}>
                <h3 style={{ fontSize: '22px', fontWeight: 800 }}>
                  {editingProduct ? 'Update Product' : 'Add New Product'}
                </h3>
                <button
                  onClick={closeModal}
                  style={{ background: 'none', border: 'none', color: 'var(--text-secondary)', cursor: 'pointer' }}
                >
                  <X size={20} />
                </button>
              </div>

              <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                <div className="input-group">
                  <label>Product Name</label>
                  <input
                    type="text"
                    placeholder="E.g., Wireless Headset"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    required
                    className="input"
                  />
                </div>

                <div className="input-group">
                  <label>Description</label>
                  <textarea
                    placeholder="Short details about the product..."
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                    className="input"
                    rows="3"
                    style={{ resize: 'none' }}
                  />
                </div>

                <div className="input-group">
                  <label>Category</label>
                  <select value={category} onChange={(e) => setCategory(e.target.value)} className="input">
                    <option value="Electronics">Electronics</option>
                    <option value="Clothing">Clothing</option>
                    <option value="Books">Books</option>
                    <option value="Home">Home</option>
                    <option value="Beauty">Beauty</option>
                    <option value="Sports">Sports</option>
                  </select>
                </div>

                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px' }}>
                  <div className="input-group">
                    <label>Price ($)</label>
                    <input
                      type="number"
                      step="0.01"
                      placeholder="99.99"
                      value={price}
                      onChange={(e) => setPrice(e.target.value)}
                      required
                      className="input"
                    />
                  </div>
                  <div className="input-group">
                    <label>Stock Quantity</label>
                    <input
                      type="number"
                      placeholder="50"
                      value={stockQuantity}
                      onChange={(e) => setStockQuantity(e.target.value)}
                      required
                      className="input"
                    />
                  </div>
                </div>

                <div className="flex gap-4" style={{ marginTop: '16px' }}>
                  <button type="button" onClick={closeModal} className="btn btn-secondary btn-full">
                    Cancel
                  </button>
                  <button type="submit" className="btn btn-primary btn-full">
                    Save Changes
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
