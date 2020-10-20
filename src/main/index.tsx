import {propagation} from '@opentelemetry/api';
import {XMLHttpRequestPlugin} from '@opentelemetry/plugin-xml-http-request';
import {FetchPlugin} from '@opentelemetry/plugin-fetch';
import {SimpleSpanProcessor} from '@opentelemetry/tracing';
import {WebTracerProvider} from '@opentelemetry/web';
import {JaegerHttpTracePropagator} from '@opentelemetry/propagator-jaeger';
import {ComlinkExporter} from './exporter';

propagation.setGlobalPropagator(new JaegerHttpTracePropagator());

const webTracerWithZone = new WebTracerProvider({
  plugins: [
    new XMLHttpRequestPlugin(),
    new FetchPlugin(),
  ],
});

webTracerWithZone.addSpanProcessor(new SimpleSpanProcessor(new ComlinkExporter(window.location, "demo-ui")));
webTracerWithZone.register();