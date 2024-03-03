import {
    Code,
    Paragraph,
    Table,
    Title,
    useTitle,
    InlineCode,
    Box,
} from "@vertex-center/components";
import "./Documentation.sass";
import { useMemo } from "react";
import { useLocation } from "react-router-dom";
import { Mermaid } from "mdx-mermaid/lib/Mermaid";

type Props = {
    content: any;
    title?: string;
};

const docs = import.meta.glob("/{docs,api}/*/{doc,api}.json", {
    eager: true,
});

export default function Documentation(props: Props) {
    const { content, title } = props;

    const location = useLocation();

    const app = location.pathname.split("/")?.[1];
    const doc: any = docs?.[`/docs/${app}/doc.json`];

    useTitle(title ?? doc?.title ?? "-");

    if (content === undefined) return null;

    const Content = content;

    const components = useMemo(
        () => ({
            p: (props: any) => <Paragraph {...props} />,
            h1: (props: any) => <Title {...props} variant="h2" />,
            h2: (props: any) => <Title {...props} variant="h3" />,
            h3: (props: any) => <Title {...props} variant="h4" />,
            h4: (props: any) => <Title {...props} variant="h5" />,
            h5: (props: any) => <Title {...props} variant="h6" />,
            hr: (props: any) => (
                <hr {...props} style={{ margin: "30px 0", opacity: 0.3 }} />
            ),
            code: (props: any) => {
                if (props.className === undefined)
                    return <InlineCode {...props} />;
                const language = /language-(\w+)/.exec(props.className || "");
                console.log(language);
                if (language?.some((l) => l === "mermaid"))
                    return (
                        <Mermaid
                            chart={props.children}
                            config={{
                                theme: {
                                    light: "dark",
                                    dark: "dark",
                                },
                            }}
                        />
                    );
                return (
                    <Code
                        style={{ marginBottom: 15 }}
                        {...props}
                        language={language?.[1]}
                    />
                );
            },
            table: (props: any) => <Table {...props} />,
            info: (props: any) => <Box type="info" {...props} />,
            tip: (props: any) => <Box type="tip" {...props} />,
            warning: (props: any) => <Box type="warning" {...props} />,
        }),
        []
    );

    return (
        <div className="documentation">
            <Content components={components} />
        </div>
    );
}
