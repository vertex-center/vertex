import { Link, useLocation } from "react-router-dom";
import styles from "./Header.module.sass";
import Logo from "../Logo/Logo";
import { useApps } from "../../hooks/useApps";

export default function Header() {
    const { apps } = useApps();

    const location = useLocation();

    let name = "Vertex";
    let to = "/app/vx-instances";
    if (location.pathname.startsWith("/app/")) {
        const app = apps?.find((app) =>
            location.pathname.includes(`/app/${app.id}`)
        );
        if (app) {
            name = app.name;
            to = `/app/${app.id}`;
        }
    }

    return (
        <header className={styles.header}>
            <Link to={to} className={styles.logo}>
                <Logo name={name} />
            </Link>
        </header>
    );
}
