import styles from "./List.module.sass";
import classNames from "classnames";
import { HTMLProps } from "react";

export type ListItemProps = HTMLProps<HTMLDivElement>;

export default function ListItem(props: Readonly<ListItemProps>) {
    const { className, onClick, ...others } = props;

    return (
        <div
            className={classNames({
                [styles.item]: true,
                [styles.itemClickable]: !!onClick,
                [className]: true,
            })}
            onClick={onClick}
            {...others}
        />
    );
}
