import { HTMLProps } from "react";
import cx from "classnames";

export type ListActionsProps = HTMLProps<HTMLDivElement>;

export function ListActions(props: Readonly<ListActionsProps>) {
    const { className, ...others } = props;
    return <div className={cx("list-actions", className)} {...others} />;
}
