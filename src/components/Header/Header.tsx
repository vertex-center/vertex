import {
    Link,
    LinkProps as RouterLinkProps,
    useLocation,
} from "react-router-dom";
import { useApps } from "../../hooks/useApps";
import {
    DropdownItem,
    Header,
    HeaderItem,
    LinkProps,
    ProfilePicture,
} from "@vertex-center/components";

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

    const accountItems = (
        <DropdownItem icon="logout" red>
            Logout
        </DropdownItem>
    );

    const account = (
        <HeaderItem items={accountItems}>
            <ProfilePicture size={36} />
        </HeaderItem>
    );

    return (
        <Header
            onClick={onClick}
            appName={app?.name}
            linkLogo={linkLogo}
            trailing={account}
        />
    );
}
