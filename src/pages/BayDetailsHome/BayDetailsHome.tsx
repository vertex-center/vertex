import { Fragment } from "react";
import { Title } from "../../components/Text/Text";
import Symbol from "../../components/Symbol/Symbol";

import styles from "./BayDetailsHome.module.sass";
import { bayNavItems } from "../BayDetails/BayDetails";
import { Link, useParams } from "react-router-dom";

export default function BayDetailsHome() {
    const { uuid } = useParams();

    return (
        <Fragment>
            <Title>Home</Title>
            <nav className={styles.nav}>
                {bayNavItems.map((item) => (
                    <Link
                        to={`/bay/${uuid}${item.to}`}
                        className={styles.navItem}
                    >
                        <Symbol
                            className={styles.navItemSymbol}
                            name={item.symbol}
                        />
                        <div>{item.label}</div>
                    </Link>
                ))}
            </nav>
        </Fragment>
    );
}
