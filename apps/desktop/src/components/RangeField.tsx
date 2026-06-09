type Props = {
  value: number;
  onChange: (v: number) => void;
  min: number;
  max: number;
  step?: number;
  label: string;
  unit?: string;
};

export default function RangeField({ value, onChange, min, max, step = 1, label, unit }: Props) {
  const pct = ((value - min) / (max - min)) * 100;
  return (
    <div className="range-field">
      <div className="range-header">
        <span className="field-label">{label}</span>
        <span className="range-value">
          {value}
          {unit && <small> {unit}</small>}
        </span>
      </div>
      <div className="range-track-wrap">
        <div className="range-fill" style={{ width: `${pct}%` }} />
        <input
          className="range-slider"
          type="range"
          min={min}
          max={max}
          step={step}
          value={value}
          onChange={(e) => onChange(Number(e.target.value))}
        />
      </div>
    </div>
  );
}
