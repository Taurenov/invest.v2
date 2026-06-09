/** Логотип Svalrix — угловая «S» внутри градиентного квадрата */
export default function BrandLogo({ size = 32 }: { size?: number }) {
  return (
    <span className="logo-dot" style={{ width: size, height: size }} aria-hidden>
      <svg viewBox="0 0 32 32" width={size} height={size} fill="none">
        <defs>
          <linearGradient id="svalrixS" x1="8" y1="6" x2="24" y2="26">
            <stop offset="0%" stopColor="#ffffff" stopOpacity="0.95" />
            <stop offset="100%" stopColor="#e0e7ff" stopOpacity="0.85" />
          </linearGradient>
        </defs>
        {/* Svalrix: острый S + молния */}
        <path
          d="M21 9H11l5 6.5L11 23h10l-4.5-6L21 9z"
          fill="url(#svalrixS)"
          opacity="0.95"
        />
        <path
          d="M22 8l2 2-3 3 2 2-3 3 2 2"
          stroke="#00d4aa"
          strokeWidth="1.2"
          strokeLinecap="round"
          strokeLinejoin="round"
          opacity="0.9"
        />
      </svg>
    </span>
  );
}
