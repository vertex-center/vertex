import { Link, useLocation } from "react-router-dom";
import styles from "./Header.module.sass";
import Logo from "../Logo/Logo";
import { apps } from "../../models/app";

export default function Header() {
    const location = useLocation();

    let name = "Vertex";
    if (location.pathname.startsWith("/app/")) {
        const app = apps.find((app) => location.pathname.includes(app.to));
        if (app) {
            name = app.name;
        }
    }

    return (
        <header className={styles.header}>
            <Link to="/app/vx-instances" className={styles.logo}>
                <Logo name={name} />
            </Link>
        </header>
    );
}
