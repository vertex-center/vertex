import { Vertical } from "../../../components/Layouts/Layouts";
import { Text, Title } from "../../../components/Text/Text";
import styles from "./CloudflareTunnels.module.sass";
import ContainerInstaller from "../../../components/ContainerInstaller/ContainerInstaller";
import { api } from "../../../backend/api/backend";

type Props = {};

export default function CloudflareTunnels(props: Readonly<Props>) {
    return (
        <Vertical gap={30}>
            <Vertical gap={20}>
                <Title className={styles.title}>Cloudflare Tunnel</Title>
                <Text className={styles.content}>
                    Cloudflare Tunnel allows you to expose your services to the
                    internet, without having to open ports or manage firewalls.
                </Text>
                <ContainerInstaller
                    name="Cloudflare Tunnel"
                    tag="Vertex Tunnels - Cloudflare"
                    install={api.vxTunnels.provider("cloudflared").install}
                />
            </Vertical>
        </Vertical>
    );
}
