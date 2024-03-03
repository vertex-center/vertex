import { api } from "../../../backend/api/backend";
import ContainerInstaller from "../../../components/ContainerInstaller/ContainerInstaller";
import { Paragraph, Title, Vertical } from "@vertex-center/components";
import Content from "../../../components/Content/Content";
import { APIError } from "../../../components/Error/APIError";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useCollector } from "../hooks/useCollector";
import Metrics from "../Metrics/Metrics";

export default function Prometheus() {
    const { collector, isLoadingCollector, errorCollector } =
        useCollector("prometheus");

    return (
        <Content>
            <Vertical gap={20}>
                <Title variant="h2">Prometheus</Title>
                <ProgressOverlay show={isLoadingCollector} />
                <Paragraph>
                    Prometheus is a Collector that gathers metrics from your
                    Vertex installation.
                </Paragraph>
                <ContainerInstaller
                    name="Prometheus"
                    tag="Vertex Monitoring - Prometheus Collector"
                    install={api.vxMonitoring.collector("prometheus").install}
                />
            </Vertical>

            <Vertical gap={20}>
                {collector?.metrics && <Title variant="h3">Metrics</Title>}
                <APIError error={errorCollector} />
                {collector?.metrics && <Metrics metrics={collector?.metrics} />}
            </Vertical>
        </Content>
    );
}
