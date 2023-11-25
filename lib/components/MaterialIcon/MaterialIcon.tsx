import { HTMLAttributes } from "react";
import classNames from "classnames";
import "./MaterialIcon.sass";

type Props = HTMLAttributes<HTMLSpanElement> & {
    icon: string;
};

export function MaterialIcon(props: Readonly<Props>) {
    const { className, icon, ...others } = props;
    return (
        <span
            className={classNames(
                "material-icon",
                "material-symbols-rounded",
                className,
            )}
            {...others}
        >
            {icon}
        </span>
    );
}
