import { HTMLProps } from "react";
import cx from "classnames";

type TabItemProps = HTMLProps<HTMLDivElement> & {
    label?: string;
    active?: boolean;
};

export function TabItem(props: Readonly<TabItemProps>) {
    const { label, active, className, ...others } = props;
    return (
        <div
            className={cx("tab-item", { "tab-item-active": active }, className)}
            {...others}
        >
            {label}
        </div>
    );
}
