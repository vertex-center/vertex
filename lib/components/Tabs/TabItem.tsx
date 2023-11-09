import { HTMLProps } from "react";
import cx from "classnames";

type TabItemProps = HTMLProps<HTMLDivElement> & {
    label?: string;
};

export function TabItem(props: Readonly<TabItemProps>) {
    const { label, className, ...others } = props;
    return (
        <div className={cx("tab-item", className)} {...others}>
            {label}
        </div>
    );
}
