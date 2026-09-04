# Carbon-Aware Kubernetes Scheduler

A Kubernetes-native tool that checks real-time electricity grid carbon intensity 
and decides whether it's a good time to run non-urgent workloads (batch jobs, 
CI runs, ML training, reports) — helping reduce the carbon footprint of 
computing by shifting work to cleaner energy windows.

## Why this project

Cloud/compute workloads consume electricity whose carbon intensity varies 
significantly throughout the day depending on the energy mix (solar, wind, 
gas, etc.). Non-urgent, delay-tolerant workloads can be shifted to run when 
the grid is cleanest — reducing emissions with zero cost and minimal 
engineering complexity. This project explores that idea using real, live 
grid data for Italy.

## Status: Work in progress

### ✅ Completed so far
- Fetches real-time carbon intensity data for a given grid zone (Italy) via 
  the [ElectricityMaps API](https://www.electricitymaps.com/)
- Threshold-based decision logic (LOW vs HIGH carbon intensity)
- Containerized with Docker
- Deployed as a Kubernetes **CronJob**, running hourly inside a local 
  Minikube cluster
- API key managed securely via a Kubernetes **Secret** (not hardcoded)
- Verified end-to-end: CronJob → Secret → Container → Live API call → 
  Real decision, confirmed via manual test runs and pod logs

### 🔜 Next steps
- Move from "prints a decision" to "actually triggers/holds back a real 
  Kubernetes Job" based on the carbon intensity decision
- Replace the fixed threshold (currently 250 gCO2eq/kWh, based on observed 
  daily min/max for Italy) with a dynamic, percentile-based threshold using 
  historical data, so it adapts to daily variation (e.g. cloudy days) 
  instead of a hardcoded number
- Refactor into a proper Kubernetes Operator/Controller pattern (custom 
  resource: `CarbonAwareSchedule`) so any workload can opt in declaratively
- Add a simple dashboard showing carbon/cost savings over time

## Architecture (current)