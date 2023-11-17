import { api } from "../../../backend/api/backend";
import ContainerInstaller from "../../../components/ContainerInstaller/ContainerInstaller";
import { Paragraph, Title } from "@vertex-center/components";
import Content from "../../../components/Content/Content";

export default function Prometheus() {
    return (
        <Content>
            <Title variant="h2">Prometheus</Title>
            <Paragraph>
                Prometheus is a Collector that gathers metrics from your Vertex
                installation.
            </Paragraph>
            <ContainerInstaller
                name="Prometheus"
                tag="Vertex Monitoring - Prometheus Collector"
                install={api.vxMonitoring.collector("prometheus").install}
            />
        </Content>
    );
}
