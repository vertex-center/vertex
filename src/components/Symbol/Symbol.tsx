import classNames from "classnames";

import styles from "./Symbol.module.sass";
import { HTMLProps } from "react";

export type SymbolProps = HTMLProps<HTMLSpanElement> & {
    name: string;
    rotating?: boolean;
};

export default function Symbol(props: Readonly<SymbolProps>) {
    const { name, rotating, className, ...others } = props;
    return (
        <span
            className={classNames({
                "material-symbols-rounded": true,
                [styles.rotating]: rotating,
                [className]: true,
            })}
            {...others}
        >
            {name}
        </span>
    );
}
