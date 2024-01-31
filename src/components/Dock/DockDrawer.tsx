import styles from "./DockDrawer.module.sass";
import classNames from "classnames";
import { Caption } from "../Text/Text";
import { Link } from "react-router-dom";
import { Vertical } from "../Layouts/Layouts";
import { useApps } from "../../hooks/useApps";
import { Logo, MaterialIcon, Title } from "@vertex-center/components";

type AppProps = {
    to: string;
    icon: string;
    name: string;
    description: string;
    onClick?: () => void;
};

function DrawerApp(props: AppProps) {
    const { to, icon, name, description, onClick } = props;

    return (
        <Link to={to} className={styles.app} onClick={onClick}>
            <MaterialIcon icon={icon} className={styles.appIcon} />
            <Vertical>
                <Title variant="h4" className={styles.appName}>
                    {name}
                </Title>
                <Caption>{description}</Caption>
            </Vertical>
        </Link>
    );
}

type Props = {
    show: boolean;
    onClose: () => void;
};

export default function DockDrawer(props: Props) {
    const { apps } = useApps();
    const { show, onClose } = props;

    return (
        <div
            className={classNames({
                [styles.drawer]: true,
                [styles.drawerOpen]: show,
            })}
        >
            <div className={styles.header} onClick={onClose}>
                <Logo size={24} />
                <div>
                    <Title variant="h3">Vertex</Title>
                </div>
            </div>
            <div className={styles.apps}>
                {[...(apps ?? [])]
                    ?.filter((app) => app.category !== "devtools")
                    ?.filter((app) => !app.hidden)
                    ?.sort((a, b) => (a.name > b.name ? 1 : -1))
                    ?.map((app) => (
                        <DrawerApp
                            key={app.id}
                            to={`/${app.id}`}
                            {...app}
                            onClick={onClose}
                        />
                    ))}
            </div>
        </div>
    );
}
