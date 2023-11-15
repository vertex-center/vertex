import SyntaxHighlighter from "react-syntax-highlighter";
import { atomOneDark as style } from "react-syntax-highlighter/dist/esm/styles/hljs";
import "./Code.sass";
import { HTMLProps } from "react";
import cx from "classnames";

export type CodeProps = HTMLProps<HTMLDivElement> & {
    language: string;
};

export function Code(props: Readonly<CodeProps>) {
    const { language, children, className, ...others } = props;

    if (typeof children !== "string") {
        console.error(
            "Code component must receive a string as children, received: ",
            typeof children,
        );
        return null;
    }

    return (
        <div className={cx("code", className)} {...others}>
            <SyntaxHighlighter language={language} style={style}>
                {children}
            </SyntaxHighlighter>
        </div>
    );
}
