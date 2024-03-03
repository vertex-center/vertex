import ContainerInstaller from "../../../components/ContainerInstaller/ContainerInstaller";
import { api } from "../../../backend/api/backend";
import { Paragraph, Title } from "@vertex-center/components";
import Content from "../../../components/Content/Content";

export default function CloudflareTunnels() {
    return (
        <Content>
            <Title variant="h2">Cloudflare Tunnel</Title>
            <Paragraph>
                Cloudflare Tunnel allows you to expose your services to the
                internet, without having to open ports or manage firewalls.
            </Paragraph>
            <ContainerInstaller
                name="Cloudflare Tunnel"
                tag="Vertex Tunnels - Cloudflare"
                install={api.vxTunnels.provider("cloudflared").install}
            />
        </Content>
    );
}
