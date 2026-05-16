use axum::{
    routing::{get, post},
    Json, Router,
};
use fin_math::{cagr_percent, forecast_linear_mean, roi_percent, savings_future_value};
use serde::{Deserialize, Serialize};
use std::net::SocketAddr;
use tower_http::cors::{Any, CorsLayer};

#[derive(Deserialize)]
struct RoiBody {
    initial: f64,
    current: f64,
}

#[derive(Serialize)]
struct RoiJson {
    roi_percent: f64,
}

#[derive(Deserialize)]
struct PredictBody {
    prices: Vec<f64>,
    horizon_days: u32,
    symbol: String,
}

#[derive(Serialize)]
struct PredictJson {
    predicted_value: f64,
    change_percent: f64,
    confidence: f64,
    model_version: String,
}

async fn http_roi(Json(body): Json<RoiBody>) -> Result<Json<RoiJson>, (axum::http::StatusCode, String)> {
    let roi = roi_percent(body.initial, body.current)
        .ok_or((axum::http::StatusCode::BAD_REQUEST, "invalid values".into()))?;
    Ok(Json(RoiJson { roi_percent: roi }))
}

async fn http_predict(Json(body): Json<PredictBody>) -> Result<Json<PredictJson>, (axum::http::StatusCode, String)> {
    if body.prices.len() < 2 {
        return Err((axum::http::StatusCode::BAD_REQUEST, "need at least 2 prices".into()));
    }
    let horizon = body.horizon_days.max(1) as usize;
    let f = forecast_linear_mean(&body.prices, horizon)
        .ok_or((axum::http::StatusCode::BAD_REQUEST, "forecast failed".into()))?;
    Ok(Json(PredictJson {
        predicted_value: f.predicted_value,
        change_percent: f.change_percent,
        confidence: f.confidence,
        model_version: "linear-mean-v1".into(),
    }))
}

#[derive(Deserialize)]
struct CagrBody {
    initial: f64,
    #[serde(rename = "final")]
    final_value: f64,
    years: f64,
}

#[derive(Serialize)]
struct CagrJson {
    cagr_percent: f64,
}

#[derive(Deserialize)]
struct SavingsBody {
    monthly: f64,
    annual_rate_pct: f64,
    months: u32,
    initial_balance: f64,
}

#[derive(Serialize)]
struct SavingsJson {
    future_value: f64,
}

async fn http_cagr(Json(body): Json<CagrBody>) -> Result<Json<CagrJson>, (axum::http::StatusCode, String)> {
    let v = cagr_percent(body.initial, body.final_value, body.years)
        .ok_or((axum::http::StatusCode::BAD_REQUEST, "invalid values".into()))?;
    Ok(Json(CagrJson { cagr_percent: v }))
}

async fn http_savings(Json(body): Json<SavingsBody>) -> Result<Json<SavingsJson>, (axum::http::StatusCode, String)> {
    if body.months == 0 {
        return Err((axum::http::StatusCode::BAD_REQUEST, "months must be > 0".into()));
    }
    let v = savings_future_value(body.monthly, body.annual_rate_pct, body.months, body.initial_balance);
    Ok(Json(SavingsJson { future_value: v }))
}

async fn health() -> &'static str {
    "ok"
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let http_addr: SocketAddr = std::env::var("ENGINE_HTTP_ADDR")
        .unwrap_or_else(|_| "127.0.0.1:50052".into())
        .parse()?;

    let cors = CorsLayer::new()
        .allow_origin(Any)
        .allow_methods(Any)
        .allow_headers(Any);

    let app = Router::new()
        .route("/health", get(health))
        .route("/v1/roi", post(http_roi))
        .route("/v1/predict", post(http_predict))
        .route("/v1/cagr", post(http_cagr))
        .route("/v1/savings", post(http_savings))
        .layer(cors);

    let listener = tokio::net::TcpListener::bind(http_addr).await?;
    println!("fin-engine HTTP listening on {http_addr}");
    axum::serve(listener, app).await?;
    Ok(())
}
