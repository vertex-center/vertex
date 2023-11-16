import { Title } from "../../../components/Text/Text";
import styles from "./Prometheus.module.sass";
import { Vertical } from "../../../components/Layouts/Layouts";
import { api } from "../../../backend/api/backend";
import ContainerInstaller from "../../../components/ContainerInstaller/ContainerInstaller";
import { Paragraph } from "@vertex-center/components";

export default function Prometheus() {
    return (
        <Vertical gap={30}>
            <Vertical gap={20}>
                <Title className={styles.title}>Prometheus</Title>
                <Paragraph className={styles.content}>
                    Prometheus is a Collector that gathers metrics from your
                    Vertex installation.
                </Paragraph>
                <ContainerInstaller
                    name="Prometheus"
                    tag="Vertex Monitoring - Prometheus Collector"
                    install={api.vxMonitoring.collector("prometheus").install}
                />
            </Vertical>
        </Vertical>
    );
}
