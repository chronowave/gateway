# Instrumenting web apps with opentelemetry-web using WebWorker
This repo demonstrates how to instrument a web app with opentelemetry using WebWorker

## Architecture
### description
Instrument web app with opentelemery-web to capture [ReadableSpan](https://github.com/open-telemetry/opentelemetry-js/blob/master/packages/opentelemetry-tracing/src/export/ReadableSpan.ts). [Copy](https://developer.mozilla.org/en-US/docs/Web/API/Web_Workers_API/Structured_clone_algorithm) spans to [WebWorker](https://developer.mozilla.org/en-US/docs/Web/API/Web_Workers_API/Using_web_workers) thread. Serialize and forward spans to collector via WebSocket in WebWorker thread.

#### high level diagram
```text
                 ┌────────────────────────────────────┐
                 │   ┌─────────┐        Kubernetes    │
                 │ ┌─┤HOT R.O.D├────┐   namespace     │
                 │ │ └─────────┘    │                 │
┌───────┐        . │                │                 │
│browser│◀──────( )┤  ┌───────┐     │     ┌──────────┐│
└───────┘        ' │  │Comlink│     │     │  Jaeger  ││
                 │ └──┤gateway│     ├────▶│all-in-one││
                 │    └──┬────┘     │     └──────────┘│
                 │       │          │                 │
                 │ ┌─────▼──────┐   │                 │
                 │ │OpenTelemtry│   │                 │
                 │ │Collector   ├───┘                 │
                 │ └────────────┘                     │
                 └────────────────────────────────────┘
```
## running example on Kubernetes
   * checkout and deploy
```shell script
git clone https://github.com/chronowave/gateway.git
cd gateway
kubectl create namespace example
helm install hotrod -n example k8s
```
   * edit /etc/hosts for ingress routing
```text
127.0.0.1        hotrod.local
```
   * simulate demo traffic via http://hotrod.local/
   * investigate traced via Jaeger UI http://localhost:16686
```shell script
kubectl port-forward svc/jaeger 16686:16686 -n example
```






