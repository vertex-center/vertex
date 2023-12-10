export type MetricType = "metric_type_on_off" | "metric_type_number";

export type Metric = {
    name: string;
    description: string;
    type: MetricType;
};

export type Collector = {
    metrics: Metric[];
};
