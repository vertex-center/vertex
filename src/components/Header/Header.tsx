import {
    Link,
    LinkProps as RouterLinkProps,
    useLocation,
} from "react-router-dom";
import { useApps } from "../../hooks/useApps";
import {
    Header,
    LinkProps,
    MaterialIcon,
    PageContext,
} from "@vertex-center/components";
import { useContext } from "react";

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

    const { title, navigation } = useContext(PageContext);

    let leading = undefined;
    const nav = navigation?.split("/") ?? [];
    let linkBack: LinkProps<RouterLinkProps> = undefined;
    if (nav.length > 0) {
        linkBack = {
            as: Link,
            to: `/app/${app?.id}/${nav.slice(0, -1).join("/")}`,
        };
        leading = <MaterialIcon icon={"arrow_back"} />;
    }

    const linkLogo: LinkProps<RouterLinkProps> = {
        as: Link,
        to,
    };

    return (
        <Header
            onClick={onClick}
            appName={app?.name}
            linkLogo={linkLogo}
            linkBack={linkBack}
            leading={leading}
        />
    );
}
