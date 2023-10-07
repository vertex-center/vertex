import styles from "./List.module.sass";
import classNames from "classnames";
import { HTMLProps } from "react";

export type ListItemProps = HTMLProps<HTMLDivElement>;

export default function ListItem(props: Readonly<ListItemProps>) {
    const { className, ...others } = props;
    return <div className={classNames(styles.item, className)} {...others} />;
}
