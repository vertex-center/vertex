import { Text, Title } from "../../../components/Text/Text";
import styles from "../Prometheus/Prometheus.module.sass";
import { Vertical } from "../../../components/Layouts/Layouts";
import InstanceInstaller from "../../../components/InstanceInstaller/InstanceInstaller";
import { api } from "../../../backend/api/backend";

export default function Grafana() {
    return (
        <Vertical gap={30}>
            <Vertical gap={20}>
                <Title className={styles.title}>Grafana</Title>
                <Text className={styles.content}>
                    Grafana allows you to visualize metrics gathered by a
                    Collector.
                </Text>
                <InstanceInstaller
                    name="Grafana"
                    tag="vertex-grafana-visualizer"
                    install={api.vxMonitoring.visualizer("grafana").install}
                />
            </Vertical>
        </Vertical>
    );
}
