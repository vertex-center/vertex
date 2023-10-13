import { Link, useLocation } from "react-router-dom";
import styles from "./Header.module.sass";
import Logo from "../Logo/Logo";
import { useApps } from "../../hooks/useApps";

type Props = {
    title?: string;
    onClick?: () => void;
};

export default function Header(props: Readonly<Props>) {
    const { title, onClick } = props;
    const { apps } = useApps();

    const location = useLocation();

    let name = title;
    let to = "/";

    if (!title) {
        to = "/app/vx-containers";
        if (location.pathname.startsWith("/app/")) {
            const app = apps?.find((app) =>
                location.pathname.includes(`/app/${app.id}`)
            );
            if (app) {
                name = app.name;
                to = `/app/${app.id}`;
            }
        }
    }

    return (
        <header className={styles.header} onClick={onClick}>
            <Link to={to} className={styles.logo}>
                <Logo name={name} />
            </Link>
        </header>
    );
}
