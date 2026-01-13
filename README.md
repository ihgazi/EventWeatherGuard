# EventWeatherGuard

EventWeatherGuard is a Go-based backend service that helps event organizers assess weather risks for outdoor events. It fetches real-time weather data from the Open-Meteo API, classifies weather conditions according to customizable rules, and provides actionable recommendations via a simple HTTP API.

---

## Setup Instructions

### Prerequisites

- Go 1.18 or later installed ([Download Go](https://golang.org/dl/))

### Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/ihgazi/EventWeatherGuard.git
   cd EventWeatherGuard
   ```

2. **Install dependencies:**
   ```sh
   go mod tidy
   ```

3. **Build and run the server:**
   ```sh
   go run main.go
   ```
   By default, the server listens on `localhost:8080`.

---

## API Usage Example

### Endpoint: `/event-forecast`

**Method:** `POST`

**Request Body:**
```json
{
  "name": "Football Match",
  "location": {
    "latitude": -28.06,
    "longitude": 156.48
  },
  "start_time": "2026-01-13T01:00:00",
  "end_time": "2026-01-13T03:00:00"
}
```

**Response Example:**
```json
{
  "classification": "Risky",
  "severity": 84,
  "summary": "Moderate rainfall and winds are expected during the event.",
  "reasons": [
    "Moderate risk: 4.5 mm rain, 34.3 km/h wind, 100% rain probability at 01:00",
    "Moderate risk: 4.5 mm rain, 36.1 km/h wind, 100% rain probability at 02:00"
  ],
  "forecast_window": [
    {
      "time": "2026-01-13T01:00:00Z",
      "rain_prob": 100,
      "precip_mm": 4.5,
      "wind_kmh": 34.3,
      "weather": "Rain Showers"
    },
    {
      "time": "2026-01-13T02:00:00Z",
      "rain_prob": 100,
      "precip_mm": 4.5,
      "wind_kmh": 36.1,
      "weather": "Rain Showers"
    }
  ]
}
```

---
## API Documentation

Interactive API documentation is available via Swagger UI.

- Start the server (`go run main.go`)
- Open [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) in your browser to explore and test the API.

---
## Weather Classification & Severity Evaluation

This service uses a **deterministic, rule-based evaluation engine** to classify event safety and compute a **granular severity score (0–100)** based on hourly weather forecasts.

The design follows a **worst-case safety model**:  
every hour within the event window is evaluated independently, and the **highest risk detected in any hour determines the final event classification**.

---

### 1. Risk Classification Logic

Each hourly forecast is evaluated using a rule-based verification logic, that classifies the given weather conditions as a specific risk level.
The classifier attempts to find the most severe rule from the current set that satisfies the given conditions.

### Weather Parameters Considered
- Hourly precipitation (mm)
- Rain probability (%)
- Wind speed (km/h)
- Weather condition (WMO-derived symbols)

### Risk Levels

| Risk Level | Trigger Conditions (Any) | Rationale |
|----------|--------------------------|-----------|
| ❌ **Unsafe** | Precipitation ≥ **10.0 mm**<br>Wind ≥ **40 km/h**<br>Weather = **Thunderstorm** | Severe risk to people, equipment, and temporary structures. |
| ⚠️ **Risky** | Precipitation ≥ **2.5 mm**<br>Wind ≥ **30 km/h**<br>Rain Probability ≥ **40%**<br>Weather = **Heavy Rain** | Conditions requiring caution or mitigation planning. |
| ✅ **Safe** | None of the above | Favorable outdoor conditions. |

> If **any single hour** is classified as *Unsafe*, the **entire event** is marked *Unsafe*.  
> Otherwise, the highest remaining level (*Risky* or *Safe*) is used.

---

### 2. Severity Score (0–100)

Beyond categorical labels, the service computes a **continuous severity score** to represent intensity.

Each numeric variable is first normalized using the threshold values used earlier. The normalized scores are combined using configurable weights:

$$
\text{Severity} =
(W_{rain} \times S_{rain}) +
(W_{prob} \times S_{prob}) +
(W_{wind} \times S_{wind})
$$

These weights emphasize each of three metrics as stronger indicators of risk.

---

### 3. WMO Weather Code Handling (Severity Caps)

Weather symbols (derived from WMO codes) are used as indicators of severe weather phenomenons (such as Thunderstorms).
To incorporate these, the system applies **minimum severity caps** based on detected conditions:

| Weather Condition | Minimum Severity |
|------------------|------------------|
| Clear / Cloudy | 0.00 |
| Rain Showers | 0.25 |
| Heavy Rain | 0.50 |
| Thunderstorm | 1.00 |

The final severity score is computed as:

$$
\text{Final Severity} = \max(\text{Weighted Score}, \text{WMO Cap})
$$

This ensures that **dangerous but low-rain scenarios** (e.g., dry thunderstorms with lightning) are still treated as severe.

---

### 4. Event-Level Aggregation

- Severity is calculated **per hour**
- The **maximum severity across all hours** is selected
- Final output severity is scaled to **0–100**
- Human-readable reasons are aggregated from all non-safe hours

---

### 5. Configuration

All classification rules and their corresponding thresholds can be configured at: `service/classification/rules.go`


---

## Key Assumptions and Trade-offs

- **Assumptions:**
  - The Open-Meteo API is available and reliable.
  - Event organizers provide accurate latitude, longitude, and event time.
  - Weather data granularity is sufficient for event planning.

- **Trade-offs:**
  - **Real-time Data:** The service fetches current/forecast data, but cannot guarantee accuracy for rapidly changing conditions.
  - **Rule Simplicity:** Classification rules are intentionally simple for clarity and maintainability; more complex models (e.g., ML-based) are not used.
  - **No Persistent Storage:** The service is stateless and does not store event or user data.
  - **Extensibility:** The modular design allows for easy integration of additional weather providers or more advanced classification logic in the future.

---
