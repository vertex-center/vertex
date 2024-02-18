import cx from "classnames";
import { HTMLProps } from "react";

export type ListItemProps = HTMLProps<HTMLDivElement>;

export function ListItem(props: Readonly<ListItemProps>) {
    const { className, onClick, ...others } = props;

    return (
        <div
            className={cx(
                "list-item",
                {
                    "list-item-clickable": !!onClick,
                },
                className,
            )}
            onClick={onClick}
            {...others}
        />
    );
}
