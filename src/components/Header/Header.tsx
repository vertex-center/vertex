import {
    Link,
    LinkProps as RouterLinkProps,
    useLocation,
    useNavigate,
} from "react-router-dom";
import { useApps } from "../../hooks/useApps";
import {
    DropdownItem,
    Header,
    HeaderItem,
    LinkProps,
    ProfilePicture,
} from "@vertex-center/components";
import useAuth from "../../apps/Auth/hooks/useAuth";

type Props = {
    title?: string;
    onClick?: () => void;
};

export default function (props: Readonly<Props>) {
    const { onClick } = props;
    const { apps } = useApps();
    const { isLoggedIn } = useAuth();

    const navigate = useNavigate();
    const location = useLocation();

    let to = "/app/containers";
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

    let accountItems;
    if (isLoggedIn) {
        accountItems = (
            <DropdownItem icon="logout" red onClick={() => navigate("/logout")}>
                Logout
            </DropdownItem>
        );
    } else {
        accountItems = (
            <DropdownItem icon="login" onClick={() => navigate("/login")}>
                Login
            </DropdownItem>
        );
    }

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
