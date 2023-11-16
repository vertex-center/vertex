import { Vertical } from "../../../components/Layouts/Layouts";
import { Title } from "../../../components/Text/Text";
import styles from "./CloudflareTunnels.module.sass";
import ContainerInstaller from "../../../components/ContainerInstaller/ContainerInstaller";
import { api } from "../../../backend/api/backend";
import { Paragraph } from "@vertex-center/components";

export default function CloudflareTunnels() {
    return (
        <Vertical gap={30}>
            <Vertical gap={20}>
                <Title className={styles.title}>Cloudflare Tunnel</Title>
                <Paragraph className={styles.content}>
                    Cloudflare Tunnel allows you to expose your services to the
                    internet, without having to open ports or manage firewalls.
                </Paragraph>
                <ContainerInstaller
                    name="Cloudflare Tunnel"
                    tag="Vertex Tunnels - Cloudflare"
                    install={api.vxTunnels.provider("cloudflared").install}
                />
            </Vertical>
        </Vertical>
    );
}
