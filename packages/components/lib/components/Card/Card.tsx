import cx from "classnames";
import { HTMLProps } from "react";
import "./Card.sass";

export type CardProps = HTMLProps<HTMLDivElement>;

export function Card(props: CardProps) {
    const { className, ...others } = props;
    return <div className={cx("card", className)} {...others} />;
}
