import { Title } from "../../../components/Text/Text";
import styles from "../Prometheus/Prometheus.module.sass";
import { Vertical } from "../../../components/Layouts/Layouts";
import ContainerInstaller from "../../../components/ContainerInstaller/ContainerInstaller";
import { api } from "../../../backend/api/backend";
import { Paragraph } from "@vertex-center/components";

export default function Grafana() {
    return (
        <Vertical gap={30}>
            <Vertical gap={20}>
                <Title className={styles.title}>Grafana</Title>
                <Paragraph className={styles.content}>
                    Grafana allows you to visualize metrics gathered by a
                    Collector.
                </Paragraph>
                <ContainerInstaller
                    name="Grafana"
                    tag="Vertex Monitoring - Grafana Visualizer"
                    install={api.vxMonitoring.visualizer("grafana").install}
                />
            </Vertical>
        </Vertical>
    );
}
