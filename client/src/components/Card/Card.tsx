import styles from "./Card.module.sass";
import { HTMLProps } from "react";
import classNames from "classnames";

type Props = HTMLProps<HTMLDivElement>;

export default function Card({ className, ...others }: Props) {
    return <div className={classNames(styles.card, className)} {...others} />;
}
