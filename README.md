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

ElectricityMaps API (live carbon intensity, zone=IT)
↓
Go program (fetch + threshold logic)
↓
Docker container
↓
Kubernetes CronJob (runs hourly, reads API key from Secret)
↓
Prints decision: run now (LOW) or wait (HIGH)


## Tech stack
- Go
- Docker
- Kubernetes (developed/tested on Minikube)
- ElectricityMaps API

## Setup

1. Get a free API key from [ElectricityMaps](https://www.electricitymaps.com/free-tier)
2. Build the image inside your Minikube Docker environment:
    eval $(minikube docker-env)
    docker build -t carbon-checker:v1 .
3. Create the Kubernetes Secret:
    kubectl create secret generic electricitymaps-secret
    --from-literal=api-key=YOUR_KEY
4. Apply the CronJob:
    kubectl apply -f k8s/cronjob.yaml
5. Trigger a manual test run:
    kubectl create job --from=cronjob/carbon-checker manual-test
    kubectl logs -l job-name=manual-test


## Lessons learned (real debugging notes)
- Docker's `buildx` builder can fail to boot its buildkit container in some 
  local setups — worked around by falling back to the classic builder 
  (`DOCKER_BUILDKIT=0`)
- Multiple local Docker daemons/contexts can cause an image built 
  successfully to still be invisible to Minikube — resolved by explicitly 
  loading the built image into Minikube (`minikube image load`)
- `imagePullPolicy: IfNotPresent` can still attempt a remote registry pull 
  in some cases; `Never` is more reliable for purely local images

## License
MIT