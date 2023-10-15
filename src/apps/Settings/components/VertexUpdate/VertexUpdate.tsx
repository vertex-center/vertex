import Button from "../../../../components/Button/Button";
import Logo from "../../../../components/Logo/Logo";
import Card from "../../../../components/Card/Card";
import { Horizontal, Vertical } from "../../../../components/Layouts/Layouts";
import Spacer from "../../../../components/Spacer/Spacer";
import URL from "../../../../components/URL/URL";
import { Caption } from "../../../../components/Text/Text";

type Props = {
    version?: string;
    description?: string;
    install: () => void;
    isInstalling?: boolean;
};

export default function VertexUpdate(props: Readonly<Props>) {
    const { version, description, install, isInstalling } = props;

    return (
        <Card>
            <Vertical gap={20}>
                <Horizontal alignItems="flex-start">
                    <Logo name={`Vertex ${version}`} />
                    <Spacer />
                    <div>
                        <Button
                            rightIcon="download"
                            onClick={install}
                            disabled={isInstalling}
                        >
                            Update
                        </Button>
                    </div>
                </Horizontal>
                <Caption>{description}</Caption>
                <div>
                    <URL
                        href="https://vertex.quentinguidee.dev/docs/changelog"
                        target="_blank"
                    >
                        Changelog
                    </URL>
                </div>
            </Vertical>
        </Card>
    );
}
