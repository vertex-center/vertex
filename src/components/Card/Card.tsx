import styles from "./Card.module.sass";
import { PropsWithChildren } from "react";

export default function Card({ children }: PropsWithChildren) {
    return <div className={styles.card}>{children}</div>;
}
