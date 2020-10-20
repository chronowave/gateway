import {ExportResult} from '@opentelemetry/core';
import {ReadableSpan, SpanExporter} from '@opentelemetry/tracing';
import {connect, send, setServiceName} from '../worker/index';

export class ComlinkExporter implements SpanExporter {
  constructor(location: Location, svcName: string | null) {
    if (svcName) {
      setServiceName(svcName)
    }
    
    const url = (location.protocol === "https:" ? "wss://" : "ws://") + location.host
    connect(url).then()
  }
  
  export(
    spans: ReadableSpan[],
    resultCallback: (result: ExportResult) => void
  ): void {
    send(spans).then(() => {
      if (resultCallback) {
        resultCallback(ExportResult.SUCCESS);
      }
    })
  }
  
  shutdown(): Promise<void> {
    return send([])
  }
}