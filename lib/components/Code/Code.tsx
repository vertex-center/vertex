import SyntaxHighlighter from "react-syntax-highlighter";
import { atomOneDark as style } from "react-syntax-highlighter/dist/esm/styles/hljs";
import "./Code.sass";
import { HTMLProps } from "react";
import cx from "classnames";

export type CodeProps = HTMLProps<HTMLDivElement> & {
    code: string;
    language: string;
};

export function Code(props: Readonly<CodeProps>) {
    const { language, code, className, ...others } = props;

    return (
        <div className={cx("code", className)} {...others}>
            <SyntaxHighlighter language={language} style={style}>
                {code}
            </SyntaxHighlighter>
        </div>
    );
}
