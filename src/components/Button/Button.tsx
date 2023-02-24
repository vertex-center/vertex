import { PropsWithChildren } from "react";

import styles from "./Button.module.sass";
import Symbol from "../Symbol/Symbol";
import classNames from "classnames";

type Props = PropsWithChildren<{
    rightSymbol: string;
    type?: "normal" | "large";
}>;

export default function Button(props: Props) {
    const { children, rightSymbol, type } = props;

    return (
        <button
            className={classNames({
                [styles.button]: true,
                [styles.buttonLarge]: type === "large",
            })}
            type="button"
        >
            <div>{children}</div>
            <Symbol name={rightSymbol} />
        </button>
    );
}
