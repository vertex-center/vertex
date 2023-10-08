import { Link, useLocation } from "react-router-dom";
import styles from "./Header.module.sass";
import Logo from "../Logo/Logo";
import { apps } from "../../models/app";

export default function Header() {
    const location = useLocation();

    let name = "Vertex";
    let to = "/app/vx-instances";
    if (location.pathname.startsWith("/app/")) {
        const app = apps.find((app) => location.pathname.includes(app.to));
        if (app) {
            name = app.name;
            to = app.to;
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
