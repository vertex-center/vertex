import SyntaxHighlighter from "react-syntax-highlighter";
import { atomOneDark } from "react-syntax-highlighter/dist/esm/styles/hljs";

import styles from "./Code.module.sass";

type Props = {
    code: string;
    language: string;
};

export default function Code(props: Props) {
    const { language, code } = props;

    return (
        <div className={styles.code}>
            <SyntaxHighlighter language={language} style={atomOneDark}>
                {code}
            </SyntaxHighlighter>
        </div>
    );
}
