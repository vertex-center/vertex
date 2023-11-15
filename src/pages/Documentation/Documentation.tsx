import {
    Code,
    Paragraph,
    Table,
    Title,
    useTitle,
} from "@vertex-center/components";
import "./Documentation.sass";
import { useMemo } from "react";
import { InlineCode } from "../../../../vertex-components/lib";

type Props = {
    content: any;
};

export default function Documentation(props: Props) {
    const { content } = props;

    useTitle("Documentation");

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
            code: (props: any) => {
                if (props.className === undefined)
                    return <InlineCode {...props} />;
                const language = /language-(\w+)/.exec(props.className || "");
                return (
                    <Code
                        style={{ marginBottom: 15 }}
                        {...props}
                        language={language?.[1]}
                    />
                );
            },
            table: (props: any) => <Table {...props} />,
        }),
        []
    );

    return (
        <div className="documentation">
            <Content components={components} />
        </div>
    );
}
