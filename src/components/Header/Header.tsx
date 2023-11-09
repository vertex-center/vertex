import {
    Link,
    LinkProps as RouterLinkProps,
    useLocation,
} from "react-router-dom";
import { useApps } from "../../hooks/useApps";
import { Header, LinkProps } from "@vertex-center/components";

type Props = {
    title?: string;
    onClick?: () => void;
};

export default function (props: Readonly<Props>) {
    const { onClick } = props;
    const { apps } = useApps();

    const location = useLocation();

    let to = "/app/vx-containers";
    let app = undefined;
    if (location.pathname.startsWith("/app/")) {
        app = apps?.find((app) => location.pathname.includes(`/app/${app.id}`));
    }

    if (app) {
        to = `/app/${app.id}`;
    }

    const linkLogo: LinkProps<RouterLinkProps> = {
        as: Link,
        to,
    };

    return <Header onClick={onClick} appName={app?.name} linkLogo={linkLogo} />;
}
