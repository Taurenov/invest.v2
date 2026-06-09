type Props = {
  value: number;
  onChange: (v: number) => void;
  min?: number;
  max?: number;
  step?: number;
  suffix?: string;
  label?: string;
};

export default function NumberField({ value, onChange, min = 0, max, step = 1, suffix, label }: Props) {
  const clamp = (n: number) => {
    let v = n;
    if (min != null) v = Math.max(min, v);
    if (max != null) v = Math.min(max, v);
    return v;
  };

  return (
    <label className="number-field-wrap">
      {label && <span className="field-label">{label}</span>}
      <div className="number-field">
        <button type="button" className="nf-btn" onClick={() => onChange(clamp(value - step))} aria-label="-">
          −
        </button>
        <input
          className="nf-input"
          type="text"
          inputMode="decimal"
          value={Number.isFinite(value) ? value : ""}
          onChange={(e) => {
            const raw = e.target.value.replace(",", ".");
            const n = parseFloat(raw);
            if (!Number.isNaN(n)) onChange(clamp(n));
            else if (raw === "") onChange(min);
          }}
        />
        <button type="button" className="nf-btn" onClick={() => onChange(clamp(value + step))} aria-label="+">
          +
        </button>
        {suffix && <span className="nf-suffix">{suffix}</span>}
      </div>
    </label>
  );
}
