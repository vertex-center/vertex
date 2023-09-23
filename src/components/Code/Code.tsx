import SyntaxHighlighter from "react-syntax-highlighter";
import { atomOneDark as style } from "react-syntax-highlighter/dist/esm/styles/hljs";

import styles from "./Code.module.sass";
import { HTMLProps } from "react";
import classNames from "classnames";

type Props = HTMLProps<HTMLDivElement> & {
    code: string;
    language: string;
};

export default function Code(props: Props) {
    const { language, code, className, ...others } = props;

    return (
        <div className={classNames(styles.code, className)} {...others}>
            <SyntaxHighlighter language={language} style={style}>
                {code}
            </SyntaxHighlighter>
        </div>
    );
}
