import { Vertical } from "../../../components/Layouts/Layouts";
import styles from "./MetricsList.module.sass";
import Metrics from "../Metrics/Metrics";
import { api } from "../../../backend/backend";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { APIError } from "../../../components/Error/APIError";
import { Title } from "../../../components/Text/Text";
import { useQuery } from "@tanstack/react-query";

export default function MetricsList() {
    const {
        data: metrics,
        isLoading,
        error,
    } = useQuery({
        queryKey: ["metrics"],
        queryFn: api.vxMonitoring.metrics,
    });

    return (
        <Vertical gap={20}>
            <ProgressOverlay show={isLoading} />
            <Title className={styles.title}>Metrics</Title>
            <APIError error={error} />
            <Metrics metrics={metrics} />
        </Vertical>
    );
}
