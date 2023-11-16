import { HTMLProps } from "react";
import cx from "classnames";

export type ListDescriptionProps = HTMLProps<HTMLDivElement>;

export function ListDescription(props: Readonly<ListDescriptionProps>) {
    const { className, ...others } = props;
    return <div className={cx("list-description", className)} {...others} />;
}
