import Logo from "../../../../components/Logo/Logo";
import Card from "../../../../components/Card/Card";
import { Horizontal, Vertical } from "../../../../components/Layouts/Layouts";
import URL from "../../../../components/URL/URL";
import { Caption } from "../../../../components/Text/Text";

type Props = {
    version?: string;
    description?: string;
};

export default function VertexUpdate(props: Readonly<Props>) {
    const { version, description } = props;

    return (
        <Card>
            <Vertical gap={20}>
                <Horizontal alignItems="flex-start">
                    <Logo name={`Vertex ${version}`} />
                </Horizontal>
                <Caption>{description}</Caption>
                <div>
                    <URL
                        href="https://docs.vertex.arra.red/changelog"
                        target="_blank"
                    >
                        Changelog
                    </URL>
                </div>
            </Vertical>
        </Card>
    );
}
