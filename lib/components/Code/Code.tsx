import SyntaxHighlighter from "react-syntax-highlighter";
import { atomOneDark as style } from "react-syntax-highlighter/dist/esm/styles/hljs";
import "./Code.sass";
import { HTMLProps } from "react";
import cx from "classnames";

type Props = HTMLProps<HTMLDivElement> & {
    code: string;
    language: string;
};

export default function Code(props: Readonly<Props>) {
    const { language, code, className, ...others } = props;

    return (
        <div className={cx("code", className)} {...others}>
            <SyntaxHighlighter language={language} style={style}>
                {code}
            </SyntaxHighlighter>
        </div>
    );
}
