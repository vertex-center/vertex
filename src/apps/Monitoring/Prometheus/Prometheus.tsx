import { Text, Title } from "../../../components/Text/Text";
import styles from "./Prometheus.module.sass";
import { Vertical } from "../../../components/Layouts/Layouts";
import { api } from "../../../backend/api/backend";
import ContainerInstaller from "../../../components/ContainerInstaller/ContainerInstaller";

export default function Prometheus() {
    return (
        <Vertical gap={30}>
            <Vertical gap={20}>
                <Title className={styles.title}>Prometheus</Title>
                <Text className={styles.content}>
                    Prometheus is a Collector that gathers metrics from your
                    Vertex installation.
                </Text>
                <ContainerInstaller
                    name="Prometheus"
                    tag="vertex-prometheus-collector"
                    install={api.vxMonitoring.collector("prometheus").install}
                />
            </Vertical>
        </Vertical>
    );
}
