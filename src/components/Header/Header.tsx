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
    Title,
} from "@vertex-center/components";
import {
    createContext,
    PropsWithChildren,
    useContext,
    useMemo,
    useState,
} from "react";

type Props = {
    title?: string;
    onClick?: () => void;
};

export const HeaderContext = createContext<{
    title: string;
    setTitle: any;
    navigation: string;
    setNavigation: any;
}>({
    title: undefined,
    setTitle: undefined,
    navigation: undefined,
    setNavigation: undefined,
});

export function HeaderProvider(props: PropsWithChildren) {
    const { children } = props;

    const [title, setTitle] = useState<string>(undefined);
    const [navigation, setNavigation] = useState<string>(undefined);

    const value = useMemo(
        () => ({ title, setTitle, navigation, setNavigation }),
        [title, setTitle, navigation, setNavigation]
    );

    return (
        <HeaderContext.Provider value={value}>
            {children}
        </HeaderContext.Provider>
    );
}

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

    const { title, navigation } = useContext(HeaderContext);

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
        >
            {title && <Title variant="h1">{title}</Title>}
        </Header>
    );
}
