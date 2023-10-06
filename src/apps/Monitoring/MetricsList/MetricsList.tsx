import { Vertical } from "../../../components/Layouts/Layouts";
import styles from "./MetricsList.module.sass";
import Metrics from "../Metrics/Metrics";
import { useFetch } from "../../../hooks/useFetch";
import { api } from "../../../backend/backend";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { APIError } from "../../../components/Error/Error";
import { Title } from "../../../components/Text/Text";
import { Metric } from "../../../models/metrics";

export default function MetricsList() {
    const {
        data: metrics,
        loading: loadingMetrics,
        error: errorMetrics,
    } = useFetch<Metric[]>(api.metrics.get);

    return (
        <Vertical gap={20}>
            <ProgressOverlay show={loadingMetrics} />
            <Title className={styles.title}>Metrics</Title>
            <APIError error={errorMetrics} />
            <Metrics metrics={metrics} />
        </Vertical>
    );
}
