import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Eye, ShoppingCart, ShoppingBag, Search, Filter } from 'lucide-react';
import api from '../api';
import { useCart } from '../context/CartContext';

const CATEGORIES = ['All', 'Electronics', 'Clothing', 'Books', 'Home', 'Beauty', 'Sports'];

export default function Home() {
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [selectedCategory, setSelectedCategory] = useState('All');
  const [searchQuery, setSearchQuery] = useState('');
  const { addToCart } = useCart();
  const navigate = useNavigate();

  useEffect(() => {
    const fetchProducts = async () => {
      setLoading(true);
      try {
        const categoryParam = selectedCategory === 'All' ? '' : selectedCategory;
        const res = await api.get(`/products?category=${categoryParam}`);
        setProducts(res.data.data.products || []);
      } catch (err) {
        console.error('Failed to fetch products:', err);
      } finally {
        setLoading(false);
      }
    };
    fetchProducts();
  }, [selectedCategory]);

  const filteredProducts = products.filter((p) =>
    (p.name || '').toLowerCase().includes((searchQuery || '').toLowerCase()) ||
    (p.description || '').toLowerCase().includes((searchQuery || '').toLowerCase())
  );

  return (
    <div style={{ minHeight: '100vh', paddingBottom: '80px' }} className="animate-fade">
      {/* Hero Section */}
      <section
        style={{
          position: 'relative',
          padding: '160px 0 100px',
          background: 'radial-gradient(circle at 50% -20%, rgba(245, 158, 11, 0.1) 0%, rgba(248, 250, 252, 0) 60%)',
          textAlign: 'center',
          borderBottom: '1px solid var(--border)',
        }}
      >
        <div className="container">
          <span
            className="badge badge-purple"
            style={{ marginBottom: '16px', textTransform: 'uppercase', letterSpacing: '0.1em' }}
          >
            Spring Collection 2026
          </span>
          <h1
            style={{
              fontSize: 'clamp(40px, 6vw, 72px)',
              fontWeight: 850,
              lineHeight: 1.1,
              letterSpacing: '-0.03em',
              marginBottom: '20px',
              maxWidth: '800px',
              margin: '0 auto 20px',
            }}
          >
            Elevate Your Style with <span className="gradient-text">Gebeya</span>
          </h1>
          <p
            style={{
              fontSize: 'clamp(16px, 2vw, 19px)',
              color: 'var(--text-secondary)',
              maxWidth: '600px',
              margin: '0 auto 32px',
            }}
          >
            Discover carefully curated premium quality essentials designed to enhance your daily lifestyle.
          </p>

          {/* Search bar inside hero */}
          <div
            style={{
              position: 'relative',
              maxWidth: '500px',
              margin: '0 auto',
            }}
          >
            <input
              type="text"
              placeholder="Search premium goods..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="input"
              style={{
                paddingLeft: '48px',
                height: '52px',
                borderRadius: '99px',
                fontSize: '16px',
                boxShadow: 'var(--shadow-card)',
              }}
            />
            <Search
              size={18}
              style={{
                position: 'absolute',
                left: '20px',
                top: '50%',
                transform: 'translateY(-50%)',
                color: 'var(--text-muted)',
              }}
            />
          </div>
        </div>
      </section>

      {/* Main Listing Section */}
      <section className="container" style={{ marginTop: '60px' }}>
        {/* Filter Toolbar */}
        <div
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            flexWrap: 'wrap',
            gap: '20px',
            marginBottom: '40px',
          }}
        >
          <div
            style={{
              display: 'flex',
              gap: '8px',
              overflowX: 'auto',
              paddingBottom: '8px',
              maxWidth: '100%',
            }}
          >
            {CATEGORIES.map((cat) => (
              <button
                key={cat}
                onClick={() => setSelectedCategory(cat)}
                className={`btn btn-sm ${selectedCategory === cat ? 'btn-primary' : 'btn-secondary'}`}
                style={{ borderRadius: '99px' }}
              >
                {cat}
              </button>
            ))}
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px', color: 'var(--text-secondary)' }}>
            <Filter size={16} />
            <span style={{ fontSize: '14px', fontWeight: 500 }}>
              Showing {filteredProducts.length} products
            </span>
          </div>
        </div>

        {/* Loading skeleton or listing */}
        {loading ? (
          <div className="products-grid">
            {[...Array(6)].map((_, idx) => (
              <div key={idx} className="card" style={{ height: '380px', display: 'flex', flexDirection: 'column', gap: '16px', padding: '16px' }}>
                <div className="skeleton" style={{ flex: 1 }}></div>
                <div className="skeleton" style={{ height: '24px', width: '60%' }}></div>
                <div className="skeleton" style={{ height: '16px', width: '90%' }}></div>
                <div className="skeleton" style={{ height: '36px', width: '40%' }}></div>
              </div>
            ))}
          </div>
        ) : filteredProducts.length === 0 ? (
          <div style={{ textAlign: 'center', padding: '80px 24px', background: 'var(--bg-secondary)', borderRadius: 'var(--radius-lg)' }}>
            <p style={{ color: 'var(--text-secondary)', fontSize: '16px' }}>No products found matching your criteria.</p>
          </div>
        ) : (
          <div className="products-grid">
            {filteredProducts.map((product) => (
              <div key={product.id} className="card flex flex-col" style={{ overflow: 'hidden' }}>
                {/* Thumbnail container */}
                <div
                  style={{
                    height: '240px',
                    background: 'var(--bg-secondary)',
                    position: 'relative',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    borderBottom: '1px solid var(--border)',
                  }}
                >
                  <span
                    className="badge badge-purple"
                    style={{ position: 'absolute', top: '16px', left: '16px', fontSize: '10px' }}
                  >
                    {product.category}
                  </span>

                  {product.stock_quantity === 0 ? (
                    <span
                      className="badge badge-red"
                      style={{ position: 'absolute', top: '16px', right: '16px', fontSize: '10px' }}
                    >
                      Out of Stock
                    </span>
                  ) : product.stock_quantity < 5 ? (
                    <span
                      className="badge badge-yellow"
                      style={{ position: 'absolute', top: '16px', right: '16px', fontSize: '10px' }}
                    >
                      Only {product.stock_quantity} Left
                    </span>
                  ) : (
                    <span
                      className="badge badge-green"
                      style={{ position: 'absolute', top: '16px', right: '16px', fontSize: '10px' }}
                    >
                      In Stock
                    </span>
                  )}

                  {/* Placeholder SVG style rendering */}
                  <div
                    style={{
                      width: '80px',
                      height: '80px',
                      borderRadius: '50%',
                      background: 'linear-gradient(135deg, var(--accent-glow), transparent)',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                    }}
                  >
                    <ShoppingBag size={36} color="var(--accent-from)" />
                  </div>
                </div>

                {/* Details */}
                <div className="p-6 flex flex-col" style={{ flex: 1 }}>
                  <h3
                    style={{
                      fontSize: '18px',
                      fontWeight: 700,
                      marginBottom: '8px',
                      color: 'var(--text-primary)',
                    }}
                  >
                    {product.name}
                  </h3>
                  <p
                    style={{
                      fontSize: '13px',
                      color: 'var(--text-secondary)',
                      display: '-webkit-box',
                      WebkitLineClamp: 2,
                      WebkitBoxOrient: 'vertical',
                      overflow: 'hidden',
                      marginBottom: '20px',
                      height: '40px',
                    }}
                  >
                    {product.description}
                  </p>

                  <div className="flex items-center justify-between" style={{ marginTop: 'auto' }}>
                    <span style={{ fontSize: '20px', fontWeight: 800, color: 'var(--text-primary)' }}>
                      ${product.price.toFixed(2)}
                    </span>
                    <div style={{ display: 'flex', gap: '8px' }}>
                      <button
                        onClick={() => navigate(`/products/${product.id}`)}
                        className="btn btn-secondary btn-sm"
                        style={{ padding: '8px', borderRadius: '50%' }}
                        title="View Details"
                      >
                        <Eye size={16} />
                      </button>
                      <button
                        onClick={() => addToCart(product.id, 1)}
                        disabled={product.stock_quantity === 0}
                        className="btn btn-primary btn-sm"
                        style={{ padding: '8px', borderRadius: '50%' }}
                        title="Add to Cart"
                      >
                        <ShoppingCart size={16} />
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
