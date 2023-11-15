import { Title, useTitle } from "@vertex-center/components";
import "./Documentation.sass";
import { useMemo } from "react";

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
            h1: (props: any) => <Title {...props} variant="h2" />,
            h2: (props: any) => <Title {...props} variant="h3" />,
            h3: (props: any) => <Title {...props} variant="h4" />,
            h4: (props: any) => <Title {...props} variant="h5" />,
            h5: (props: any) => <Title {...props} variant="h6" />,
        }),
        []
    );

    return (
        <div className="documentation">
            <Content components={components} />
        </div>
    );
}
