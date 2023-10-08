import { Link } from "react-router-dom";
import styles from "./Header.module.sass";
import Logo from "../Logo/Logo";

export default function Header() {
    return (
        <header className={styles.header}>
            <Link to="/instances" className={styles.logo}>
                <Logo />
            </Link>
        </header>
    );
}
