type IconName =
  | "home"
  | "transactions"
  | "analytics"
  | "markets"
  | "portfolio"
  | "goals"
  | "calculator"
  | "forecast"
  | "extras"
  | "settings";

function G(id: string) {
  return (
    <defs>
      <linearGradient id={id} x1="0" y1="0" x2="24" y2="24">
        <stop offset="0%" stopColor="#5b6ef5" />
        <stop offset="100%" stopColor="#00d4aa" />
      </linearGradient>
    </defs>
  );
}

const paths = (g: string): Record<IconName, JSX.Element> => ({
  home: (
    <path
      d="M4 10.5L12 4l8 6.5V20a1 1 0 01-1 1h-5v-6H10v6H5a1 1 0 01-1-1v-9.5z"
      stroke="currentColor"
      strokeWidth="1.6"
      fill="none"
      strokeLinejoin="round"
    />
  ),
  transactions: (
    <>
      <path d="M7 8h10M7 12h7" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" />
      <rect x="4" y="5" width="16" height="14" rx="3" stroke="currentColor" strokeWidth="1.6" fill="none" />
    </>
  ),
  analytics: (
    <>
      <path d="M6 18V12M10 18V8M14 18V14M18 18V6" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
      <path d="M5 20h14" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" opacity="0.5" />
    </>
  ),
  markets: (
  <>
      <path d="M4 16l4-5 4 3 5-7 3 4" stroke={`url(#${g})`} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" fill="none" />
      <circle cx="18" cy="7" r="2" fill={`url(#${g})`} />
    </>
  ),
  portfolio: (
    <path
      d="M6 8V6a2 2 0 012-2h8a2 2 0 012 2v2M5 8h14v10a2 2 0 01-2 2H7a2 2 0 01-2-2V8z"
      stroke="currentColor"
      strokeWidth="1.6"
      fill="none"
      strokeLinejoin="round"
    />
  ),
  goals: (
    <>
      <circle cx="12" cy="12" r="8" stroke="currentColor" strokeWidth="1.6" fill="none" />
      <circle cx="12" cy="12" r="4" stroke={`url(#${g})`} strokeWidth="1.6" fill="none" />
      <circle cx="12" cy="12" r="1.2" fill={`url(#${g})`} />
    </>
  ),
  calculator: (
    <>
      <rect x="5" y="3" width="14" height="18" rx="3" stroke="currentColor" strokeWidth="1.6" fill="none" />
      <rect x="8" y="6" width="8" height="3" rx="1" fill={`url(#${g})`} opacity="0.8" />
      <circle cx="8.5" cy="13" r="1" fill="currentColor" />
      <circle cx="12" cy="13" r="1" fill="currentColor" />
      <circle cx="15.5" cy="13" r="1" fill="currentColor" />
      <circle cx="8.5" cy="17" r="1" fill="currentColor" />
      <circle cx="12" cy="17" r="1" fill="currentColor" />
      <circle cx="15.5" cy="17" r="1" fill={`url(#${g})`} />
    </>
  ),
  forecast: (
    <>
      <path d="M5 14c2-4 4-6 7-6s5 2 7 6" stroke={`url(#${g})`} strokeWidth="2" fill="none" strokeLinecap="round" />
      <path d="M12 4v2M16.2 5.8l-1.4 1.4M19 10h-2" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" opacity="0.7" />
    </>
  ),
  extras: (
    <>
      <path d="M12 3l1.8 4.5L18 9l-4.2 1.5L12 15l-1.8-4.5L6 9l4.2-1.5L12 3z" fill={`url(#${g})`} opacity="0.9" />
      <path d="M5 18h14" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" opacity="0.4" />
    </>
  ),
  settings: (
    <>
      <circle cx="12" cy="12" r="3" stroke="currentColor" strokeWidth="1.6" fill="none" />
      <path
        d="M12 3v2M12 19v2M3 12h2M19 12h2M5.6 5.6l1.4 1.4M17 17l1.4 1.4M5.6 18.4l1.4-1.4M17 7l1.4-1.4"
        stroke="currentColor"
        strokeWidth="1.4"
        strokeLinecap="round"
        opacity="0.75"
      />
    </>
  ),
});

export default function NavIcon({ name, active }: { name: IconName; active?: boolean }) {
  const gid = `navGrad-${name}`;
  return (
    <svg
      className={`nav-svg ${active ? "active" : ""}`}
      width="22"
      height="22"
      viewBox="0 0 24 24"
      aria-hidden
    >
      {G(gid)}
      {paths(gid)[name]}
    </svg>
  );
}
