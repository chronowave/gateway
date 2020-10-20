import {CollectorTraceExporter, toCollectorExportTraceServiceRequest} from '@opentelemetry/exporter-collector';
import {ReadableSpan} from '@opentelemetry/tracing';
import {opentelemetryProto} from "@opentelemetry/exporter-collector/build/src/types";

let serviceName: string = '{{ .ServiceName }}';
let urlOverride: string = '{{ .WsUrl }}';

let wsc: WebSocket | null
let buffered: opentelemetryProto.collector.trace.v1.ExportTraceServiceRequest[] = []
let exporter = new CollectorTraceExporter({serviceName: serviceName});

export function setServiceName(svc: string) {
  exporter = new CollectorTraceExporter({serviceName: svc});
}

export async function send(spans: ReadableSpan[]) {
  let req = toCollectorExportTraceServiceRequest(spans, exporter);
  if (wsc) {
    wsc.send(JSON.stringify(req));
  } else {
    if (buffered.length < 200) {
      buffered.push(req)
    }
  }
}

let timeout: number = 0

export async function connect(url: string) {
  const connURL = urlOverride && urlOverride.length > 0 ? urlOverride : (url + "/ws")
  let ws = new WebSocket(connURL, ["access_token", "{{ .AccessToken }}"]);
  let tid: NodeJS.Timeout
  ws.onopen = () => {
    clearTimeout(tid)
    timeout = 0
    buffered.forEach((item, index) => {
      ws.send(JSON.stringify(item))
      buffered.splice(index, 1);
    });
    wsc = ws
  };
  
  ws.onmessage = (e) => {
    // ignore
  };
  
  ws.onclose = (e) => {
    tid = setTimeout(() => connect(url), timeout);
  };
  
  ws.onerror = (e) => {
    clearTimeout(tid)
    timeout += 1000
    wsc = null
    ws.close();
  };
}
