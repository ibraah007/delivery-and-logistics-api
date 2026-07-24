import React from 'react';

export default function ManLionLogo({ className = '', style = {} }) {
  return (
    <svg viewBox="0 0 500 350" xmlns="http://www.w3.org/2000/svg" className={className} style={style}>
      <defs>
        <linearGradient id="chromeGradient" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor="#e6f0fa" />
          <stop offset="35%" stopColor="#8a9ba8" />
          <stop offset="50%" stopColor="#ffffff" />
          <stop offset="65%" stopColor="#4a5a68" />
          <stop offset="100%" stopColor="#1e2936" />
        </linearGradient>
        <filter id="badgeGlow" x="-20%" y="-20%" width="140%" height="140%">
          <feGaussianBlur stdDeviation="8" result="blur" />
          <feComponentTransfer in="blur" result="glow1">
            <feFuncA type="linear" slope="0.6" />
          </feComponentTransfer>
          <feMerge>
            <feMergeNode in="glow1" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>
      </defs>
      <g filter="url(#badgeGlow)">
        <path d="M 50 40 L 450 40 C 450 160 380 270 250 310 C 120 270 50 160 50 40 Z" fill="none" stroke="url(#chromeGradient)" strokeWidth="12" strokeLinejoin="round"/>
        <path d="M 68 56 L 432 56 C 432 165 368 258 250 294 C 132 258 68 165 68 56 Z" fill="rgba(15, 23, 42, 0.85)" stroke="url(#chromeGradient)" strokeWidth="4"/>
        <g fill="url(#chromeGradient)">
          <path d="M 235 90 C 245 82 265 82 275 92 C 285 102 280 115 292 120 C 302 122 312 112 320 125 C 326 135 315 145 322 155 C 328 162 338 165 335 178 C 330 190 315 188 308 198 C 300 208 305 220 295 228 C 285 235 272 222 260 230 C 250 238 248 250 235 248 C 225 245 228 230 218 222 C 210 215 198 220 190 210 C 182 200 195 190 190 180 C 185 172 172 170 175 158 C 178 148 190 148 195 138 C 200 128 190 118 200 108 C 210 98 222 102 235 90 Z" />
          <path d="M 285 145 L 340 130 L 325 150 L 350 155 L 310 178 Z" />
          <path d="M 185 185 C 160 170 150 140 162 115 C 168 102 178 95 175 88 C 170 82 155 95 148 112 C 135 145 150 190 180 208 Z" />
        </g>
        <rect x="120" y="250" width="260" height="35" rx="6" fill="url(#chromeGradient)" />
        <text x="250" y="276" fill="#0f172a" fontSize="24" fontWeight="900" fontFamily="Arial, sans-serif" textAnchor="middle" letterSpacing="6">MAN</text>
      </g>
    </svg>
  );
}
