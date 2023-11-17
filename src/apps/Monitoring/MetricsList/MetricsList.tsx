import Metrics from "../Metrics/Metrics";
import { api } from "../../../backend/api/backend";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { APIError } from "../../../components/Error/APIError";
import { useQuery } from "@tanstack/react-query";
import Content from "../../../components/Content/Content";
import { Title } from "@vertex-center/components";

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
        <Content>
            <Title variant="h2">Metrics</Title>
            <APIError error={error} />
            <ProgressOverlay show={isLoading} />
            <Metrics metrics={metrics} />
        </Content>
    );
}
