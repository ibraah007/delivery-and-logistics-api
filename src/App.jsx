import React from 'react';
import ManLionLogo from './components/ManLionLogo';

export default function App() {
  return (
    <div style={{
      position: 'relative',
      minHeight: '100vh',
      backgroundColor: '#060b13',
      color: '#fff',
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      overflow: 'hidden'
    }}>
      {/* Background Shining Logo Watermark */}
      <div style={{
        position: 'absolute',
        top: '50%',
        left: '50%',
        transform: 'translate(-50%, -50%)',
        width: '500px',
        maxWidth: '80vw',
        opacity: 0.2,
        pointerEvents: 'none'
      }}>
        <ManLionLogo style={{ width: '100%', height: 'auto' }} />
      </div>

      {/* Main UI Content Layer */}
      <div style={{ position: 'relative', zIndex: 1, textAlign: 'center' }}>
        <h1 style={{ fontSize: '2rem', marginBottom: '1rem', color: '#ffcc00' }}>
          Mobikey MAN Fleet Command
        </h1>
        <p style={{ color: '#94a3b8' }}>Background watermark logo is active!</p>
      </div>
    </div>
  );
}
