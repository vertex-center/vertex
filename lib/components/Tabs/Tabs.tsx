import React, { Children, HTMLProps, useState } from "react";
import cx from "classnames";
import "./Tabs.sass";

export type TabsProps = HTMLProps<HTMLDivElement>;

export function Tabs(props: Readonly<TabsProps>) {
    const { className, children, ...others } = props;
    const [active, setActive] = useState(0);
    // @ts-ignore
    const child = children[active] ?? <></>;
    return (
        <div className={cx("tabs", className)} {...others}>
            <div className="tabs-row">
                {Children.map(children, (child: any, i) => {
                    return React.cloneElement(child, {
                        onClick: () => setActive(i),
                    });
                })}
            </div>
            <div key={child.props.label}>{child.props.children}</div>
        </div>
    );
}
