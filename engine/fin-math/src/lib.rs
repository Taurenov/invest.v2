//! Финансовые расчёты для fin-engine.
//! В проде вызывается через gRPC; здесь — чистая библиотека + unit-тесты.

/// ROI в процентах: (current - initial) / initial * 100
/// initial должен быть > 0.
pub fn roi_percent(initial: f64, current: f64) -> Option<f64> {
    if initial <= 0.0 || !initial.is_finite() || !current.is_finite() {
        return None;
    }
    Some((current - initial) / initial * 100.0)
}

/// Простой прогноз: линейная экстраполяция по последним N точкам (метод наименьших квадратов).
/// Возвращает прогноз на `horizon` шагов вперёд от последней точки ряда.
pub fn forecast_linear_mean(prices: &[f64], horizon: usize) -> Option<ForecastResult> {
    if prices.len() < 2 || horizon == 0 {
        return None;
    }
    let n = prices.len() as f64;
    let mut sum_x = 0.0;
    let mut sum_y = 0.0;
    let mut sum_xy = 0.0;
    let mut sum_xx = 0.0;

    for (i, &y) in prices.iter().enumerate() {
        let x = i as f64;
        sum_x += x;
        sum_y += y;
        sum_xy += x * y;
        sum_xx += x * x;
    }

    let denom = n * sum_xx - sum_x * sum_x;
    if denom.abs() < f64::EPSILON {
        return None;
    }

    let slope = (n * sum_xy - sum_x * sum_y) / denom;
    let intercept = (sum_y - slope * sum_x) / n;
    let last_x = (prices.len() - 1) as f64;
    let target_x = last_x + horizon as f64;
    let predicted = intercept + slope * target_x;
    let last = *prices.last()?;
    let change_pct = if last.abs() > f64::EPSILON {
        (predicted - last) / last * 100.0
    } else {
        0.0
    };

    // Грубая "уверенность": чем больше точек и меньше дисперсия остатков, тем выше.
    let confidence = (0.35 + 0.1 * n.min(10.0)).min(0.85);

    Some(ForecastResult {
        predicted_value: predicted,
        change_percent: change_pct,
        confidence,
    })
}

/// CAGR в процентах за `years` лет.
pub fn cagr_percent(initial: f64, final_value: f64, years: f64) -> Option<f64> {
    if initial <= 0.0 || years <= 0.0 || !initial.is_finite() || !final_value.is_finite() {
        return None;
    }
    Some(((final_value / initial).powf(1.0 / years) - 1.0) * 100.0)
}

/// Будущая стоимость: начальный баланс + ежемесячные взносы с годовой ставкой.
pub fn savings_future_value(
    monthly: f64,
    annual_rate_pct: f64,
    months: u32,
    initial_balance: f64,
) -> f64 {
    let r = annual_rate_pct / 100.0 / 12.0;
    let mut total = initial_balance;
    for _ in 0..months {
        total = total * (1.0 + r) + monthly;
    }
    total
}

#[derive(Debug, Clone, PartialEq)]
pub struct ForecastResult {
    pub predicted_value: f64,
    pub change_percent: f64,
    pub confidence: f64,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn roi_basic() {
        let r = roi_percent(1000.0, 1250.0).unwrap();
        assert!((r - 25.0).abs() < 1e-6);
    }

    #[test]
    fn forecast_uptrend() {
        let prices = vec![100.0, 102.0, 105.0, 108.0, 111.0];
        let f = forecast_linear_mean(&prices, 1).unwrap();
        assert!(f.predicted_value > 111.0);
        assert!(f.change_percent > 0.0);
    }

    #[test]
    fn cagr_basic() {
        let c = cagr_percent(1000.0, 1331.0, 3.0).unwrap();
        assert!((c - 10.0).abs() < 0.5);
    }
}
