import ContainerInstaller from "../../../components/ContainerInstaller/ContainerInstaller";
import { api } from "../../../backend/api/backend";
import { Paragraph, Title } from "@vertex-center/components";
import Content from "../../../components/Content/Content";

export default function Grafana() {
    return (
        <Content>
            <Title variant="h2">Grafana</Title>
            <Paragraph>
                Grafana allows you to visualize metrics gathered by a Collector.
            </Paragraph>
            <ContainerInstaller
                name="Grafana"
                tag="Vertex Monitoring - Grafana Visualizer"
                install={api.vxMonitoring.visualizer("grafana").install}
            />
        </Content>
    );
}
