import { PropsWithChildren } from "react";
import styles from "./VersionTag.module.sass";

export default function VersionTag({ children }: PropsWithChildren) {
    return <span className={styles.tag}>{children}</span>;
}
