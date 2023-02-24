import { PropsWithChildren } from "react";

import styles from "./Button.module.sass";
import Symbol from "../Symbol/Symbol";
import classNames from "classnames";

type Props = PropsWithChildren<{
    rightSymbol: string;
    type?: "normal" | "large";
    downloading?: boolean;
    onClick: () => void;
}>;

export default function Button(props: Props) {
    const { children, rightSymbol, type, downloading, onClick } = props;

    const content = (
        <div className={styles.content}>
            <div>{children}</div>
            <Symbol name={rightSymbol} />
        </div>
    );

    const downloadingContent = (
        <div className={classNames(styles.content, styles.contentDownloading)}>
            <div>Downloading</div>
            <Symbol name="sentiment_satisfied" rotating />
        </div>
    );

    return (
        <button
            className={classNames({
                [styles.button]: true,
                [styles.buttonLarge]: type === "large",
                [styles.buttonDownloading]: downloading,
            })}
            type="button"
            onClick={onClick}
        >
            {content}
            {downloadingContent}
        </button>
    );
}
