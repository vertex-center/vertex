import styles from "./DockDrawer.module.sass";
import classNames from "classnames";
import { BigTitle, Caption, Title } from "../Text/Text";
import { Link } from "react-router-dom";
import { Vertical } from "../Layouts/Layouts";
import { apps } from "./Dock";
import Icon from "../Icon/Icon";
import LogoIcon from "../Logo/LogoIcon";

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
            <Icon name={icon} className={styles.appIcon} />
            <Vertical gap={8}>
                <Title className={styles.appName}>
                    <LogoIcon />
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
    const { show, onClose } = props;
    return (
        <div
            className={classNames({
                [styles.drawer]: true,
                [styles.drawerOpen]: show,
            })}
        >
            <BigTitle className={styles.title}>Apps</BigTitle>
            <div className={styles.apps}>
                {apps.map((app) => (
                    <DrawerApp
                        key={app.to}
                        to={app.to}
                        icon={app.icon}
                        name={app.name}
                        description={app.description}
                        onClick={onClose}
                    />
                ))}
            </div>
        </div>
    );
}
