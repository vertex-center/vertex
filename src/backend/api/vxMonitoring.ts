import { Metric } from "../../models/metrics";
import { server } from "./backend";

const BASE_URL = "/app/vx-monitoring";

const getMetrics = async () => {
    const { data } = await server.get<Metric[]>(`${BASE_URL}/metrics`);
    return data;
};

const installCollector = (collector: string) => {
    return server.post(`${BASE_URL}/collector/${collector}/install`);
};

const installVisualizer = (visualizer: string) => {
    return server.post(`${BASE_URL}/visualizer/${visualizer}/install`);
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
