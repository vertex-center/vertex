import { Metric } from "../../models/metrics";

import { createServer } from "../server";

// @ts-ignore
const server = createServer(window.api_urls.monitoring);

const getMetrics = async () => {
    const { data } = await server.get<Metric[]>(`/metrics`);
    return data;
};

const installCollector = (collector: string) => {
    return server.post(`/collector/${collector}/install`);
};

const installVisualizer = (visualizer: string) => {
    return server.post(`/visualizer/${visualizer}/install`);
};

const collectorRoutes = (collector: string) => ({
    install: () => installCollector(collector),
});

const visualizerRoutes = (visualizer: string) => ({
    install: () => installVisualizer(visualizer),
});

export const vxMonitoringRoutes = {
    metrics: getMetrics,
    collector: collectorRoutes,
    visualizer: visualizerRoutes,
};
