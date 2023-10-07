import { HTMLProps } from "react";
import classNames from "classnames";

import styles from "./List.module.sass";

export type ListInfoProps = HTMLProps<HTMLDivElement>;

export default function ListInfo(props: Readonly<ListInfoProps>) {
    const { className, ...others } = props;
    return <div className={classNames(styles.info, className)} {...others} />;
}
