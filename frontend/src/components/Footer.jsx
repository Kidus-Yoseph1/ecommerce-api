import React from 'react';
import { ShoppingBag } from 'lucide-react';

export default function Footer() {
  return (
    <footer
      style={{
        background: 'var(--bg-secondary)',
        borderTop: '1px solid var(--border)',
        padding: '40px 0',
        marginTop: 'auto',
      }}
    >
      <div
        className="container"
        style={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          gap: '16px',
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px', fontSize: '18px', fontWeight: 700 }}>
          <ShoppingBag size={20} style={{ color: 'var(--accent-from)' }} />
          <span>Gebeya</span>
        </div>
        <p style={{ color: 'var(--text-secondary)', fontSize: '14px', textAlign: 'center' }}>
          &copy; {new Date().getFullYear()} Gebeya Inc. All rights reserved.
        </p>
      </div>
    </footer>
  );
}
